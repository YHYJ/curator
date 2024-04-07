/*
File: clone.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 15:16:00

Description: 子命令 'clone' 的实现
*/

package cli

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gookit/color"
	"github.com/yhyj/curator/general"
)

// updateGitConfig 更新 .git/config 文件
//
// 参数：
//   - configFile: .git/config 文件路径
//   - originalLink: 需要替换的原始链接
//   - newLink: 替换上去的新链接
//
// 返回：
//   - 错误信息
func updateGitConfig(configFile, originalLink, newLink string) error {
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
	pushUrl1 := "" // 第一行 pushurl
	pushUrl2 := "" // 第二行 pushurl

	// 逐行读取文件内容
	for scanner.Scan() {
		line := scanner.Text()

		// 检索一次模糊匹配的行
		if !matched && regex.MatchString(line) {
			// 第一次匹配：将可能存在的 "ssh://" 删除，并在"/"多于1个时将第1个替换为":"
			// 该次匹配是专对子模块的 .git/config 的处理
			line = strings.Replace(line, "ssh://", "", 1)
			if strings.Count(line, "/") >= 2 {
				line = strings.Replace(line, "/", ":", 1)
			}
			lines = append(lines, line)
			// 第二次匹配：创建2行 "pushurl"
			// 该次匹配是对于 .git/config 的通用处理
			pushUrl1 = strings.ReplaceAll(line, "url", "pushurl")
			pushUrl2 = strings.ReplaceAll(pushUrl1, originalLink, newLink)
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

// runScript 运行 shell 脚本
//
// 参数：
//   - filePath: 脚本所在目录
//   - scriptName: 脚本名
//
// 返回：
//   - 错误信息
func runScript(filePath, scriptName string) error {
	// 判断是否存在脚本文件，存在则运行脚本，不存在则忽略
	if general.FileExist(filepath.Join(filePath, scriptName)) {
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

// RollingCloneRepos 遍历克隆远端仓库到本地
//
// 参数：
//   - confile: 程序配置文件
//   - source: 远端仓库源，支持 'github' 和 'gitea'，默认为 'github'
func RollingCloneRepos(confile, source string) {
	// 加载配置文件
	conf, err := GetTomlConfig(confile)
	if err != nil {
		color.Error.Println(err)
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
		repoNames := conf.Get("git.repos").([]interface{})
		scriptNameList := conf.Get("script.name_list").([]interface{})
		// 获取公钥
		publicKeys, err := general.GetPublicKeysByGit(pemfile.(string))
		if err != nil {
			color.Error.Println(err)
			return
		}

		// 确定仓库源
		repoSource := func() map[string]string {
			switch source {
			case "github":
				return map[string]string{
					"repoSourceUrl":      githubUrl,
					"repoSourceUsername": githubUsername,
					"originalLink":       githubLink,
					"newLink":            giteaLink,
				}
			case "gitea":
				return map[string]string{
					"repoSourceUrl":      giteaUrl,
					"repoSourceUsername": giteaUsername,
					"originalLink":       giteaLink,
					"newLink":            githubLink,
				}
			default:
				return map[string]string{
					"repoSourceUrl":      githubUrl,
					"repoSourceUsername": githubUsername,
					"originalLink":       githubLink,
					"newLink":            giteaLink,
				}
			}
		}()

		// 克隆
		color.Info.Tips("%s %s\n", general.FgWhite("Clone to"), general.PrimaryText(storagePath))
		for _, repoName := range repoNames {
			repoPath := filepath.Join(storagePath, repoName.(string))
			// 开始克隆
			color.Printf("%s %s %s: ", general.RunFlag, general.LightText("Cloning"), general.FgCyan(repoName.(string)))
			// 克隆前检测是否存在同名本地仓库或非空文件夹
			if general.FileExist(repoPath) {
				isRepo, _ := general.IsLocalRepo(repoPath)
				if isRepo { // 是本地仓库
					color.Printf("%s %s\n", general.FgBlue(general.UnmodifiedFlag), general.SecondaryText("Local repository already exists"))
					// 添加一个延时，使输出更加顺畅
					general.Delay(0.1)
					continue
				} else { // 不是本地仓库
					if general.FolderEmpty(repoPath) { // 是空文件夹，删除后继续克隆
						if err := general.DeleteFile(repoPath); err != nil {
							color.Error.Println(err)
						}
					} else { // 文件夹非空，处理下一个
						color.Printf("%s %s\n", general.WarningFlag, general.WarnText("Folder is not a local repository and not empty"))
						// 添加一个延时，使输出更加顺畅
						general.Delay(0.1)
						continue
					}
				}
			}
			repo, err := general.CloneRepoViaSSH(repoPath, repoSource["repoSourceUrl"], repoSource["repoSourceUsername"], repoName.(string), publicKeys)
			if err != nil { // Clone 失败
				color.Error.Println(err)
			} else { // Clone 成功
				length := len(general.RunFlag) + len("Cloning") // 仓库信息缩进长度
				color.Printf("%s %s\n", general.SuccessFlag, general.CommentText("Receive object completed"))
				var errList []string // 使用一个 Slice 存储所有错误信息以美化输出
				// 执行脚本
				for _, scriptName := range scriptNameList {
					if err := runScript(repoPath, scriptName.(string)); err != nil {
						errList = append(errList, "Run script "+scriptName.(string)+": "+err.Error())
					}
				}
				// 处理主仓库的配置文件 .git/config
				configFile := filepath.Join(repoPath, ".git", "config")
				if err = updateGitConfig(configFile, repoSource["originalLink"], repoSource["newLink"]); err != nil {
					errList = append(errList, "Update repository git config (main): "+err.Error())
				}
				// 获取主仓库的 worktree
				worktree, err := repo.Worktree()
				if err != nil {
					errList = append(errList, "Get local repository worktree: "+err.Error())
				}
				// 获取主仓库的远程分支信息
				remoteBranchs, err := general.GetRepoBranchInfo(worktree, "remote")
				if err != nil {
					errList = append(errList, "Get local repository branch (remote): "+err.Error())
				}
				// 根据远程分支 refs/remotes/origin/<remoteBranchName> 创建本地分支 refs/heads/<localBranchName>
				otherErrList := general.CreateLocalBranch(repo, remoteBranchs)
				errList = append(errList, otherErrList...)
				// 获取主仓库的本地分支信息
				var localBranchStr []string
				localBranchs, err := general.GetRepoBranchInfo(worktree, "local")
				if err != nil {
					errList = append(errList, "Get local repository branch (local): "+err.Error())
				}
				for _, localBranch := range localBranchs {
					localBranchStr = append(localBranchStr, localBranch.Name())
				}
				color.Printf("%s%s %s [%s]\n", strings.Repeat(" ", length), general.JoinerFinish, general.BranchFlag, general.FgCyan(strings.Join(localBranchStr, " ")))
				// 获取子模块信息
				submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
				if err != nil {
					errList = append(errList, "Get local repository submodules: "+err.Error())
				}
				for index, submodule := range submodules {
					// 创建和主模块的连接符
					joiner := func() string {
						if index == len(submodules)-1 {
							return general.JoinerFinish
						}
						return general.JoinerIng
					}()
					color.Printf("%s%s %s %s\n", strings.Repeat(" ", length), joiner, general.SubmoduleFlag, general.FgMagenta(submodule.Config().Name))
					// 处理子模块的配置文件 .git/modules/<submodule>/config
					configFile := filepath.Join(repoPath, ".git", "modules", submodule.Config().Name, "config")
					if err = updateGitConfig(configFile, repoSource["originalLink"], repoSource["newLink"]); err != nil {
						errList = append(errList, "Update repository git config (submodule): "+err.Error())
					}
				}
				// 输出克隆完成后其他操作产生的错误信息
				for _, err := range errList {
					color.Error.Println(err)
				}
			}
			// 添加一个延时，使输出更加顺畅
			general.Delay(0.1)
		}
	}
}
