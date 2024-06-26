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
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s -> Unable to load config: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		return
	}

	// 获取公钥
	publicKeys, err := general.GetPublicKeysByGit(config.SSH.RsaFile)
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s -> Unable to get public key: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
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

	// 为已 Clone 的存储库计数
	totalNum := len(config.Git.Repos) // 总存储库数
	clonedRepo := make([]string, 0)   // 已 Clone 存储库
	for _, repoName := range config.Git.Repos {
		repoPath := filepath.Join(config.Storage.Path, repoName) // 本地存储库路径
		if general.FileExist(repoPath) {
			isRepo, _, _ := general.IsLocalRepo(repoPath)
			if isRepo {
				clonedRepo = append(clonedRepo, repoName)
			}
		}
	}

	// 开始 Clone 提示
	negatives := strings.Builder{}
	negatives.WriteString(color.Sprintf("%s Clone repository from %s, %d/%d cloned\n", general.InfoText("INFO:"), general.FgGreenText(source), len(clonedRepo), totalNum))
	negatives.WriteString(color.Sprintf("%s Repository root: %s\n", general.InfoText("INFO:"), general.PrimaryText(config.Storage.Path)))

	// 让用户选择需要 Clone 的存储库
	selectedRepos, err := general.MultipleSelectionFilter(config.Git.Repos, clonedRepo, negatives.String())
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s -> Unable to start selector: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		return
	}

	// 留屏信息
	if len(selectedRepos) > 0 {
		negatives.WriteString(color.Sprintf("%s Selected: %s\n", general.InfoText("INFO:"), general.FgCyanText(strings.Join(selectedRepos, ", "))))
		negatives.WriteString(color.Sprintf("%s", strings.Repeat(general.Separator1st, general.SeparatorBaseLength)))
		color.Println(negatives.String())
	}

	// 遍历所选存储库名
	for _, repoName := range selectedRepos {
		repoPath := filepath.Join(config.Storage.Path, repoName)
		// 开始克隆提示
		actionPrint := color.Sprintf("%s Cloning %s: ", general.RunFlag, general.FgCyanText(repoName))
		general.WaitSpinner.Prefix = actionPrint
		general.WaitSpinner.Start()
		// 克隆前检测是否存在同名本地仓库或非空文件夹
		if general.FileExist(repoPath) {
			isRepo, _, _ := general.IsLocalRepo(repoPath)
			if isRepo { // 是本地仓库
				general.WaitSpinner.Stop()
				color.Printf("%s%s %s\n", actionPrint, general.FgBlueText(general.LatestFlag), general.SecondaryText("Local repository already exists"))
				// 添加一个延时，使输出更加顺畅
				general.Delay(0.1)
				continue
			} else { // 不是本地仓库
				if general.FolderEmpty(repoPath) { // 是空文件夹，删除后继续克隆
					if err := general.DeleteFile(repoPath); err != nil {
						general.WaitSpinner.Stop()
						color.Printf("%s", actionPrint)
						fileName, lineNo := general.GetCallerInfo()
						color.Printf("%s %s -> Unable to delete file: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
						continue
					}
				} else { // 文件夹非空，处理下一个
					general.WaitSpinner.Stop()
					color.Printf("%s%s %s\n", actionPrint, general.WarningFlag, general.WarnText("Folder is not a local repository and not empty"))
					// 添加一个延时，使输出更加顺畅
					general.Delay(0.1)
					continue
				}
			}
		}

		// 开始克隆
		repo, err := general.CloneRepoViaSSH(repoPath, repoSource["repoSourceUrl"], repoSource["repoSourceUsername"], repoName, publicKeys)

		// 克隆结束
		if err != nil { // Clone 失败
			general.WaitSpinner.Stop()
			color.Printf("%s", actionPrint)
			fileName, lineNo := general.GetCallerInfo()
			color.Printf("%s %s -> Unable to clone repository: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		} else { // Clone 成功
			length := len(general.RunFlag) + len("Cloning") // 仓库信息缩进长度
			general.WaitSpinner.Stop()
			color.Printf("%s%s %s ", actionPrint, general.SuccessFlag, general.FgGreenText("Receive object completed"))
			var errList []string // 使用一个 Slice 存储所有错误信息以美化输出
			// 执行脚本
			for _, scriptName := range config.Script.RunQueue {
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
			remoteBranchs, err := general.GetRepoBranchInfo(worktree, false, "", "remote")
			if err != nil {
				errList = append(errList, "Get local repository branch (remote): "+err.Error())
			}
			// 根据远程分支 refs/remotes/origin/<remoteBranchName> 创建本地分支 refs/heads/<localBranchName>
			otherErrList := general.CreateLocalBranch(repo, remoteBranchs)
			errList = append(errList, otherErrList...)
			// 获取主仓库的本地分支信息
			var localBranchStr []string
			localBranchs, err := general.GetRepoBranchInfo(worktree, false, "", "local")
			if err != nil {
				errList = append(errList, "Get local repository branch (local): "+err.Error())
			}
			for _, localBranch := range localBranchs {
				localBranchStr = append(localBranchStr, localBranch.Name())
			}
			color.Printf("%s\n", general.SecondaryText("[", strings.Join(localBranchStr, " "), "]"))

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
				actionPrint := color.Sprintf("%s%s %s %s ", strings.Repeat(" ", length), joiner, general.SubmoduleFlag, general.FgMagentaText(submodule.Config().Name))
				general.WaitSpinner.Prefix = actionPrint
				general.WaitSpinner.Start()
				// 获取子模块的 worktree
				isRepo, submoduleRepo, _ := general.IsLocalRepo(submodule.Config().Path)
				if isRepo {
					// 获取子模块的远程分支信息
					submoduleRemoteBranchs, err := general.GetRepoBranchInfo(worktree, true, submodule.Config().Name, "remote")
					if err != nil {
						errList = append(errList, "Get local repository branch (remote): "+err.Error())
					}
					// 根据远程分支 modules/<submoduleName>/refs/remotes/origin/<remoteBranchName> 创建本地分支 modules/<submoduleName>/refs/heads/<localBranchName>
					clbErrList := general.CreateLocalBranch(submoduleRepo, submoduleRemoteBranchs)
					errList = append(errList, clbErrList...)
					// 获取子模块的 worktree
					submoduleWorktree, err := submoduleRepo.Worktree()
					if err != nil {
						errList = append(errList, "Get local repository worktree: "+err.Error())
					}
					// 获取子模块默认分支名
					submoduleDefaultBranchName, gdbnErrList := general.GetDefaultBranchName(submoduleRepo, publicKeys)
					errList = append(errList, gdbnErrList...)
					// 切换到默认分支
					if err := general.CheckoutBranch(submoduleWorktree, submoduleDefaultBranchName); err != nil {
						errList = append(errList, "Checkout to default branch: "+err.Error())
					}
					// 获取子模块的本地分支信息
					var submoduleLocalBranchStr []string
					submoduleLocalBranchs, err := general.GetRepoBranchInfo(worktree, true, submodule.Config().Name, "local")
					if err != nil {
						errList = append(errList, "Get local repository branch (local): "+err.Error())
					}
					for _, submoduleLocalBranch := range submoduleLocalBranchs {
						submoduleLocalBranchStr = append(submoduleLocalBranchStr, submoduleLocalBranch.Name())
					}
					general.WaitSpinner.Stop()
					color.Printf("%s%s\n", actionPrint, general.SecondaryText("[", strings.Join(submoduleLocalBranchStr, " "), "]"))
				} else { // 子模块非本地仓库
					general.WaitSpinner.Stop()
					color.Printf("%s%s %s\n", actionPrint, general.ErrorFlag, general.DangerText("Folder is not a local repository"))
				}
			}
			// 输出克隆完成后其他操作产生的错误信息
			fileName, lineNo := general.GetCallerInfo()
			for _, err := range errList {
				color.Printf("%s %s -> Other error info: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+2, "]"), err)
			}
		}
		// 添加一个延时，使输出更加顺畅
		general.Delay(0.1)
	}
}
