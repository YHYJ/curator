/*
File: run.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 15:16:00

Description:
*/

package cli

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/yhyj/clone-repos/general"
)

// 更新.git/config文件
func updateGitConfig(configFile, githubLink, giteaLink string) error {
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
func runScript(filePath, scriptName string) error {
	// 判断是否存在脚本文件，存在则运行脚本，不存在则忽略
	if general.FileExist(filePath + "/" + scriptName) {
		// 进到指定目录
		if err := os.Chdir(filePath); err != nil {
			return err
		}
		// 运行脚本
		bashArgs := []string{scriptName}
		if err := general.RunCommand("bash", bashArgs); err != nil {
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
		githubLink := githubUrl + ":" + githubUsername
		giteaUsername := conf.Get("git.gitea_username").(string)
		giteaLink := giteaUrl + ":" + giteaUsername
		repos := conf.Get("git.repos").([]interface{})
		scriptNameList := conf.Get("script.name_list").([]interface{})
		// 定义变量
		var interval = 100 * time.Millisecond
		// 获取公钥
		publicKeys, err := general.GetPublicKeysByGit(pemfile.(string))
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
			return
		}

		// 克隆
		fmt.Printf("Clone to: \x1b[32;1m%s\x1b[0m\n\n", storagePath)
		for _, repo := range repos {
			// 提示信息
			fmt.Printf("\x1b[32;1m==>\x1b[0m Cloning \x1b[36;1m%s\x1b[0m: ", repo.(string))
			repoPath := storagePath + "/" + repo.(string)
			// 克隆前检测是否存在同名本地仓库或非空文件夹
			if general.FileExist(repoPath) {
				isRepo, _ := general.IsLocalRepo(repoPath)
				if isRepo { // 是本地仓库
					fmt.Printf("\x1b[32m[✔]\x1b[0m \x1b[34m%s\x1b[0m\n", "Local repo already exists")
					// 添加一个延时，使输出更加顺畅
					time.Sleep(interval)
					continue
				} else { // 不是本地仓库
					if general.FolderEmpty(repoPath) { // 空文件夹则删除
						if err := general.DeleteFile(repoPath); err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
						}
					} else { // 文件夹非空
						fmt.Println("Folder is not a local repo and is not empty")
						// 添加一个延时，使输出更加顺畅
						time.Sleep(interval)
						continue
					}
				}
			}
			// 开始克隆
			repo, err := general.CloneRepoViaSSH(repoPath, githubUrl, githubUsername, repo.(string), publicKeys)
			if err != nil { // Clone失败
				fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
			} else { // Clone成功
				fmt.Printf("\x1b[32m[✔]\x1b[0m ")
				var errList []string // 使用一个Slice存储所有错误信息以美化输出
				// 执行脚本
				for _, scriptName := range scriptNameList {
					if err := runScript(repoPath, scriptName.(string)); err != nil {
						errList = append(errList, "Run Script "+scriptName.(string)+": "+err.Error())
					}
				}
				// 处理主仓库的配置文件.git/config
				configFile := repoPath + "/" + ".git/config"
				if err = updateGitConfig(configFile, githubLink, giteaLink); err != nil {
					errList = append(errList, "Update Git Config (main): "+err.Error())
				}
				// 获取主仓库的worktree
				worktree, err := repo.Worktree()
				if err != nil {
					errList = append(errList, "Get Local Repo Worktree: "+err.Error())
				}
				// 获取主仓库的远程分支信息
				remoteBranchs, err := general.GetRepoBranchInfo(worktree, "remote")
				if err != nil {
					errList = append(errList, "Get Local Repo Branch (remote): "+err.Error())
				}
				// 根据远程分支refs/remotes/origin/<remoteBranchName>创建本地分支refs/heads/<localBranchName>
				otherErrList := general.CreateLocalBranch(repo, remoteBranchs)
				errList = append(errList, otherErrList...)
				// 获取主仓库的本地分支信息
				localBranchs, err := general.GetRepoBranchInfo(worktree, "local")
				if err != nil {
					errList = append(errList, "Get Local Repo Branch (local): "+err.Error())
				}
				var localBranchStr string
				for _, localBranch := range localBranchs {
					localBranchStr = localBranchStr + localBranch.Name() + ", "
				}
				// 获取子模块信息
				submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
				var submoduleStr string
				if err != nil {
					errList = append(errList, "Get Local Repo Submodules: "+err.Error())
				}
				for _, submodule := range submodules {
					submoduleStr = submoduleStr + submodule.Config().Name + ", "
					// 处理子模块的配置文件.git/modules/<submodule>/config
					configFile := fmt.Sprintf("%s/%s/%s/%s", repoPath, ".git/modules", submodule.Config().Name, "config")
					if err := updateGitConfig(configFile, githubLink, giteaLink); err != nil {
						errList = append(errList, "Update Git Config (submodule): "+err.Error())
					}
				}
				// 处理并输出本地分支和子模块信息
				// TODO: 需要添加子模块的分支信息，但获取困难 <13-10-23, YJ> //
				localBranchStr = strings.TrimRight(localBranchStr, ", ")
				submoduleStr = strings.TrimRight(submoduleStr, ", ")
				if len(submoduleStr) == 0 { // 分支常有而子模块不常有
					fmt.Printf("Branch: \x1b[33;1m%s\x1b[0m\n", localBranchStr)
				} else {
					fmt.Printf("Branch: \x1b[33;1m%s\x1b[0m Submodule: \x1b[35m%s\x1b[0m\n", localBranchStr, submoduleStr)
				}
				// 输出克隆完成后其他操作产生的错误信息
				for _, err := range errList {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
				}
			}
			// 添加一个延时，使输出更加顺畅
			time.Sleep(interval)
		}
	}
}
