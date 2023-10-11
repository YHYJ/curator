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
	"github.com/go-git/go-git/v5/plumbing/transport"
	gssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pelletier/go-toml"
	"golang.org/x/crypto/ssh"
)

func getTomlConfig(filename string) (*toml.Tree, error) {
	tree, err := toml.LoadFile(filename)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func getSshKeyAuth(privateSshKeyFile string) transport.AuthMethod {
	var auth transport.AuthMethod
	sshKey, _ := os.ReadFile(privateSshKeyFile)
	signer, _ := ssh.ParsePrivateKey([]byte(sshKey))
	auth = &gssh.PublicKeys{User: "git", Signer: signer}
	return auth
}

// 更新.git/config文件
func updateGitConfig(filePath, githubLink, giteaLink string) (err error) {
	fileName := ".git/config"                                               // 文件名
	file, err := os.OpenFile(filePath+"/"+fileName, os.O_RDWR, os.ModePerm) // 打开文件以读写模式
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // 创建一个扫描器来读取文件内容
	var lines []string                // 存储读取到的行

	regexPattern := `.*url\s*=\s*.*github\.com.*` // 正则表达式用于模糊匹配行的内容
	regex := regexp.MustCompile(regexPattern)     // 创建正则表达式

	matched := false // 是否匹配到，用于限制只匹配一次

	pushUrl1 := "" // 第一行pushurl
	pushUrl2 := "" // 第二行pushurl

	// 逐行读取文件内容
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line) // 将读取到的行存入lines

		if !matched {
			if regex.MatchString(line) { // 检索模糊匹配的行
				pushUrl1 = strings.ReplaceAll(line, "url", "pushurl")
				pushUrl2 = strings.ReplaceAll(pushUrl1, githubLink, giteaLink)
				lines = append(lines, pushUrl1)
				lines = append(lines, pushUrl2)
				matched = true
			}
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

func RollingCLoneRepos(confile string) {
	// 加载配置文件
	conf, err := getTomlConfig(confile)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
	} else {
		// 获取配置项
		private_key_file := conf.Get("ssh.private_key_file")
		path := conf.Get("storage.path").(string)
		githubUrl := conf.Get("git.github_url").(string)
		githubUsername := conf.Get("git.github_username").(string)
		giteaUrl := conf.Get("git.gitea_url").(string)
		giteaUsername := conf.Get("git.gitea_username").(string)
		repos := conf.Get("git.repos").([]interface{})
		scriptNameList := conf.Get("script.name_list").([]interface{})
		auth := getSshKeyAuth(private_key_file.(string))
		// 开始克隆
		fmt.Printf("Clone to: \x1b[32;1m%s\x1b[0m\n\n", path)
		for _, repo := range repos {
			fmt.Printf("\x1b[32;1m==>\x1b[0m Cloning \x1b[36;1m%s\x1b[0m: ", repo.(string))
			storagePath := path + "/" + repo.(string)
			_, err := git.PlainClone(storagePath, false, &git.CloneOptions{
				URL:               "git" + "@" + githubUrl + ":" + githubUsername + "/" + repo.(string) + ".git",
				Auth:              auth,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
				Progress:          io.Discard, // os.Stdout会将Clone的详细过程输出到控制台，io.Discard会直接丢弃
			})
			if err != nil {
				if err == git.ErrRepositoryAlreadyExists {
					fmt.Printf("%s\n", err)
				} else {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
				}
			} else {
				fmt.Printf("\x1b[32m%s\x1b[0m\n", "cloning completed")
				// Clone成功后更新.git/config
				githubLink := githubUrl + ":" + githubUsername
				giteaLink := giteaUrl + ":" + giteaUsername
				err := updateGitConfig(storagePath, githubLink, giteaLink)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
				}
				// Clone成功后执行脚本
				for _, scriptName := range scriptNameList {
					err := runScript(storagePath, scriptName.(string))
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
