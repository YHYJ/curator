/*
File: clone.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 15:16:00

Description: 子命令 'clone' 的实现
*/

package cli

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/gookit/color"
	"github.com/pelletier/go-toml"
	"github.com/yhyj/curator/general"
)

// RollingCloneRepos 遍历克隆远端仓库到本地
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
//   - source: 远端仓库源，支持 'github' 和 'gitea'，默认为 'github'
func RollingCloneRepos(configTree *toml.Tree, source string) {
	// 获取配置项
	config, err := general.LoadConfigToStruct(configTree)
	if err != nil {
		color.Error.Println(err)
		return
	}

	// 获取公钥
	publicKeys, err := general.GetPublicKeysByGit(config.SSH.RsaFile)
	if err != nil {
		color.Error.Println(err)
		return
	}

	// 确定仓库源
	githubLink := config.Git.GithubUrl + ":" + config.Git.GithubUsername
	giteaLink := config.Git.GiteaUrl + ":" + config.Git.GiteaUsername
	repoSource := func() map[string]string {
		switch source {
		case "github":
			return map[string]string{
				"repoSourceUrl":      config.Git.GithubUrl,
				"repoSourceUsername": config.Git.GithubUsername,
				"originalLink":       githubLink,
				"newLink":            giteaLink,
			}
		case "gitea":
			return map[string]string{
				"repoSourceUrl":      config.Git.GiteaUrl,
				"repoSourceUsername": config.Git.GiteaUsername,
				"originalLink":       giteaLink,
				"newLink":            githubLink,
			}
		default:
			return map[string]string{
				"repoSourceUrl":      config.Git.GithubUrl,
				"repoSourceUsername": config.Git.GithubUsername,
				"originalLink":       githubLink,
				"newLink":            giteaLink,
			}
		}
	}()

	// 克隆
	color.Info.Tips("%s %s", general.FgWhiteText("Clone repository from"), general.FgGreenText(source))
	color.Info.Tips("%s: %s", general.FgWhiteText("Repository root"), general.PrimaryText(config.Storage.Path))
	// 让用户选择需要 Clone 的存储库
	selectedRepos, err := general.MultipleSelectionFilter(config.Git.Repos)
	if err != nil {
		color.Error.Println(err)
	}
	// 对所选的存储库进行排序
	sort.Strings(selectedRepos)
	// 遍历所选存储库名
	for _, repoName := range selectedRepos {
		repoPath := filepath.Join(config.Storage.Path, repoName)
		// 开始克隆
		color.Printf("%s %s %s: ", general.RunFlag, general.LightText("Cloning"), general.FgCyanText(repoName))
		// 克隆前检测是否存在同名本地仓库或非空文件夹
		if general.FileExist(repoPath) {
			isRepo, _ := general.IsLocalRepo(repoPath)
			if isRepo { // 是本地仓库
				color.Printf("%s %s\n", general.FgBlueText(general.LatestFlag), general.SecondaryText("Local repository already exists"))
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
		repo, err := general.CloneRepoViaSSH(repoPath, repoSource["repoSourceUrl"], repoSource["repoSourceUsername"], repoName, publicKeys)
		if err != nil { // Clone 失败
			color.Error.Println(err)
		} else { // Clone 成功
			length := len(general.RunFlag) + len("Cloning") // 仓库信息缩进长度
			color.Printf("%s %s\n", general.SuccessFlag, general.CommentText("Receive object completed"))
			var errList []string // 使用一个 Slice 存储所有错误信息以美化输出
			// 执行脚本
			for _, scriptName := range config.Script.NameList {
				if err := general.RunScript(repoPath, scriptName); err != nil {
					errList = append(errList, "Run script "+scriptName+": "+err.Error())
				}
			}
			// 处理主仓库的配置文件 .git/config
			configFile := filepath.Join(repoPath, ".git", "config")
			if err = general.ModifyGitConfig(configFile, repoSource["originalLink"], repoSource["newLink"]); err != nil {
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
			color.Printf("%s%s %s [%s]\n", strings.Repeat(" ", length), general.JoinerFinish, general.BranchFlag, general.FgCyanText(strings.Join(localBranchStr, " ")))
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
				color.Printf("%s%s %s %s\n", strings.Repeat(" ", length), joiner, general.SubmoduleFlag, general.FgMagentaText(submodule.Config().Name))
				// 处理子模块的配置文件 .git/modules/<submodule>/config
				configFile := filepath.Join(repoPath, ".git", "modules", submodule.Config().Name, "config")
				if err = general.ModifyGitConfig(configFile, repoSource["originalLink"], repoSource["newLink"]); err != nil {
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
