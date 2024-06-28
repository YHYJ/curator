/*
File: clone.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 15:16:00

Description: 子命令 'clone' 的实现
*/

package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/pelletier/go-toml"
	"github.com/yhyj/curator/general"
)

// RollingCloneRepos 遍历 Clone 远端存储库到本地
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
//   - source: 远端存储库源，支持 'github' 和 'gitea'，默认为 'github'
func RollingCloneRepos(configTree *toml.Tree, source string) {
	// 获取配置项
	config, err := general.LoadConfigToStruct(configTree)
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s -> Unable to load config: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		return
	}

	// 确定存储库源
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

	// 查找已存在的本地存储库
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

	// 输出基础信息
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
		// 构建本地存储库路径
		repoPath := filepath.Join(config.Storage.Path, repoName)

		// Clone
		clone(config, repoSource, repoPath, repoName, config.Script.RunQueue)

		// 添加一个延时，使输出更加顺畅
		general.Delay(general.DelayTime)
	}
}

// clone Clone 远端存储库到本地
//
// 参数：
//   - config: 配置项
//   - source: 存储库源
//   - path: 本地存储库路径
//   - name: 存储库名
//   - scripts: Clone 完成后需要执行的脚本
func clone(config *general.Config, source map[string]string, path, name string, scripts []string) {
	// 获取公钥
	publicKeys, err := general.GetPublicKeysByGit(config.SSH.RsaFile)
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s -> Unable to get public key: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		return
	}

	// 开始 Clone 提示
	actionPrint := color.Sprintf("%s Cloning %s: ", general.RunFlag, general.FgCyanText(name))
	general.WaitSpinner.Prefix = actionPrint
	general.WaitSpinner.Start()

	// Clone 前检测是否存在同名本地存储库或非空文件夹
	if general.FileExist(path) {
		isRepo, _, _ := general.IsLocalRepo(path)
		if isRepo { // 是本地存储库
			general.WaitSpinner.Stop()
			color.Printf("%s%s %s\n", actionPrint, general.FgBlueText(general.LatestFlag), general.SecondaryText("Local repository already exists"))
			// 添加一个延时，使输出更加顺畅
			general.Delay(general.DelayTime)
			return
		} else { // 不是本地存储库
			if general.FolderEmpty(path) { // 是空文件夹，删除后继续 Clone
				if err := general.DeleteFile(path); err != nil {
					general.WaitSpinner.Stop()
					color.Printf("%s", actionPrint)
					fileName, lineNo := general.GetCallerInfo()
					color.Printf("%s %s -> Unable to delete file: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
					return
				}
			} else { // 文件夹非空，处理下一个
				general.WaitSpinner.Stop()
				color.Printf("%s%s %s\n", actionPrint, general.WarningFlag, general.WarnText("Folder is not a local repository and not empty"))
				// 添加一个延时，使输出更加顺畅
				general.Delay(general.DelayTime)
				return
			}
		}
	}

	// 开始 Clone
	repo, err := general.CloneRepoViaSSH(path, source["repoSourceUrl"], source["repoSourceUsername"], name, publicKeys)

	// Clone 结束
	if err != nil { // Clone 失败
		general.WaitSpinner.Stop()
		color.Printf("%s", actionPrint)
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s -> Unable to clone repository: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
	} else { // Clone 成功
		// 成功信息
		length := len(general.RunFlag) + len("Cloning") // 存储库信息缩进长度
		general.WaitSpinner.Stop()
		color.Printf("%s%s %s ", actionPrint, general.SuccessFlag, general.FgGreenText("Receive object completed"))

		// 使用一个切片存储后续所有错误信息以美化输出
		var errList []string

		// Clone 成功后执行存储库中的 Shell 脚本来优化存储库
		for _, script := range scripts {
			if general.FileExist(filepath.Join(path, script)) {
				// 进到指定目录
				if err := os.Chdir(path); err != nil {
					errList = append(errList, "Run script "+script+": "+err.Error())
				}
				// 运行脚本
				bashArgs := []string{script}
				if err := general.RunCommandToOS("bash", bashArgs); err != nil {
					errList = append(errList, "Run script "+script+": "+err.Error())
				}
			}
		}

		// 更新主存储库的配置文件 .git/config
		configFile := filepath.Join(path, ".git", "config")
		if err = general.ModifyGitConfig(configFile, source["originalLink"], source["newLink"]); err != nil {
			errList = append(errList, "Update repository git config (main): "+err.Error())
		}

		// 获取主存储库的 worktree
		worktree, err := repo.Worktree()
		if err != nil {
			errList = append(errList, "Get local repository worktree: "+err.Error())
		}
		// 获取主存储库的远程分支信息
		remoteBranchs, err := general.GetRepoBranchInfo(worktree, false, "", "remote")
		if err != nil {
			errList = append(errList, "Get local repository branch (remote): "+err.Error())
		}
		// 根据远程分支 refs/remotes/origin/<remoteBranchName> 创建本地分支 refs/heads/<localBranchName>
		otherErrList := general.CreateLocalBranch(repo, remoteBranchs)
		errList = append(errList, otherErrList...)

		// 获取主存储库的本地分支信息
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
			// 输出子模块信息
			joiner := func() string { // 主模块和子模块的输出连接符
				if index == len(submodules)-1 {
					return general.JoinerFinish
				}
				return general.JoinerIng
			}()
			actionPrint := color.Sprintf("%s%s %s %s ", strings.Repeat(" ", length), joiner, general.SubmoduleFlag, general.FgMagentaText(submodule.Config().Name))
			general.WaitSpinner.Prefix = actionPrint
			general.WaitSpinner.Start()

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
			} else { // 子模块非本地存储库
				general.WaitSpinner.Stop()
				color.Printf("%s%s %s\n", actionPrint, general.ErrorFlag, general.DangerText("Folder is not a local repository"))
			}
		}

		// 输出 Clone 完成后其他操作产生的错误信息
		fileName, lineNo := general.GetCallerInfo()
		for _, err := range errList {
			color.Printf("%s %s -> Other error info: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+2, "]"), err)
		}
	}
}
