/*
File: run.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 15:16:00

Description:
*/

package function

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

// 更新.git/config文件
func updateGitConfig(configFile, githubLink, giteaLink string) (err error) {
	// 以读写模式打开文件
	file, err := os.OpenFile(configFile, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取文件
	scanner := bufio.NewScanner(file) // 创建一个扫描器来读取文件内容
	var lines []string                // 存储读取到的行

	// 正则匹配（主仓库和子模块的匹配规则一样）
	regexPattern := `.*url\s*=\s*.*[:\/].*\.git` // 定义正则匹配规则
	regex := regexp.MustCompile(regexPattern)    // 创建正则表达式
	matched := false                             // 是否匹配到，用于限制只匹配一次

	// 需要新增的行
	pushUrl1 := "" // 第一行pushurl
	pushUrl2 := "" // 第二行pushurl

	// 逐行读取文件内容
	for scanner.Scan() {
		line := scanner.Text()

		// 检索一次模糊匹配的行
		if !matched && regex.MatchString(line) {
			// 第一次匹配：将可能存在的"ssh://"删除，并在"/"多于1个时将第1个替换为":"
			// 该次匹配是专对子模块的.git/config的处理
			line = strings.Replace(line, "ssh://", "", 1)
			if strings.Count(line, "/") >= 2 {
				line = strings.Replace(line, "/", ":", 1)
			}
			lines = append(lines, line)
			// 第二次匹配：创建2行"pushurl"
			// 该次匹配是对于.git/config的通用处理
			pushUrl1 = strings.ReplaceAll(line, "url", "pushurl")
			pushUrl2 = strings.ReplaceAll(pushUrl1, githubLink, giteaLink)
			lines = append(lines, pushUrl1)
			lines = append(lines, pushUrl2)
			matched = true
		} else {
			lines = append(lines, line)
		}
	}

	// 将修改后的内容写回文件
	file.Truncate(0) // 清空文件内容
	file.Seek(0, 0)  // 移动光标到文件开头
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, _ = writer.WriteString(line + "\n")
	}
	writer.Flush()

	return nil
}

// 运行脚本
func runScript(filePath, scriptName string) (err error) {
	// 判断是否存在脚本文件，存在则运行脚本，不存在则忽略
	if FileExist(filePath + "/" + scriptName) {
		// 进到指定目录
		err = os.Chdir(filePath)
		if err != nil {
			return err
		}
		// 运行脚本
		bashArgs := []string{scriptName}
		err = RunCommand("bash", bashArgs)
		if err != nil {
			return err
		}
	}
	return nil
}

func RollingCloneRepos(confile string) {
	// 加载配置文件
	conf, err := GetTomlConfig(confile)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
	} else {
		// 获取配置项
		pemfile := conf.Get("ssh.rsa_file")
		storagePath := conf.Get("storage.path").(string)
		githubUrl := conf.Get("git.github_url").(string)
		githubUsername := conf.Get("git.github_username").(string)
		giteaUrl := conf.Get("git.gitea_url").(string)
		giteaUsername := conf.Get("git.gitea_username").(string)
		repos := conf.Get("git.repos").([]interface{})
		scriptNameList := conf.Get("script.name_list").([]interface{})
		publicKeys, err := GetPublicKeysByGit(pemfile.(string), "") // TODO: 需要处理有password的情况 <11-10-23, YJ> //
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
			return
		}
		// 开始克隆
		fmt.Printf("Clone to: \x1b[32;1m%s\x1b[0m\n\n", storagePath)
		for _, repo := range repos {
			fmt.Printf("\x1b[32;1m==>\x1b[0m Cloning \x1b[36;1m%s\x1b[0m: ", repo.(string))
			repoPath := storagePath + "/" + repo.(string)
			_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
				URL:               "git" + "@" + githubUrl + ":" + githubUsername + "/" + repo.(string) + ".git",
				Auth:              publicKeys,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
				Progress:          io.Discard, // os.Stdout会将Clone的详细过程输出到控制台，io.Discard会直接丢弃
			})
			if err != nil { // Clone失败
				if err == git.ErrRepositoryAlreadyExists { // 错误原因是本地git仓库已存在
					fmt.Printf("%s\n", err)
				} else { // 其他错误
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
				}
			} else {
				fmt.Printf("\x1b[32m%s\x1b[0m\n", "cloning completed")
				// Clone成功后更新.git/config
				githubLink := githubUrl + ":" + githubUsername
				giteaLink := giteaUrl + ":" + giteaUsername
				// 处理主仓库的.git/config
				configFile := repoPath + "/" + ".git/config"
				err := updateGitConfig(configFile, githubLink, giteaLink)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
				}
				// 处理子模块的.git/config
				submodules, err := GetSubModuleNames(repoPath)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
					continue
				}
				for _, submodule := range submodules {
					configFile := fmt.Sprintf("%s/%s/%s/%s", repoPath, ".git/modules", submodule.Config().Name, "config")
					err := updateGitConfig(configFile, githubLink, giteaLink)
					if err != nil {
						fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
						continue
					}
				}
				// Clone成功后执行脚本
				for _, scriptName := range scriptNameList {
					err := runScript(repoPath, scriptName.(string))
					if err != nil {
						fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
					}
				}
			}
			// 添加一个0.01秒的延时，使输出更加顺畅
			time.Sleep(100 * time.Millisecond)
		}
	}
}
