/*
File: pull.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-03-05 14:22:54

Description: 子命令 'pull' 的实现
*/

package cli

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
	"github.com/yhyj/curator/general"
)

// RollingPullRepos 遍历 Pull 远端存储库的更改到本地
//
// 参数：
//   - config: 解析 toml 配置文件得到的配置项
//   - source: 远端存储库源，支持 'github' 和 'gitea'，默认为 'github'
func RollingPullRepos(config *general.Config, source string) {
	// 为已存在的本地存储库计数
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

	// 显示项排序
	sort.Strings(config.Git.Repos)

	// 输出基础信息
	negatives := strings.Builder{}
	negatives.WriteString(color.Sprintf("%s Pull repository from %s: %d/%d cloned\n", general.InfoText("INFO:"), general.FgGreenText(source), len(clonedRepo), totalNum))
	negatives.WriteString(color.Sprintf("%s Repository root: %s\n", general.InfoText("INFO:"), general.PrimaryText(config.Storage.Path)))

	// 让用户选择需要 Pull 的存储库
	selectedRepos, err := general.MultipleSelectionFilter(config.Git.Repos, clonedRepo, negatives.String())
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		return
	}

	// 选择项排序
	sort.Strings(selectedRepos)

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

		// Pull
		pull(config, repoPath, repoName)

		// 添加一个延时，使输出更加顺畅
		general.Delay(general.DelayTime)
	}
}

// pull Pull 远端存储库的更改到本地
//
// 参数：
//   - config: 配置项目
//   - path: 本地存储库路径
//   - name: 存储库名
func pull(config *general.Config, path, name string) {
	// 获取公钥
	publicKeys, err := general.GetPublicKeysByGit(config.SSH.RsaFile)
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		return
	}

	// 开始 Pull 提示
	actionPrint := color.Sprintf("%s Pulling %s: ", general.RunFlag, general.FgCyanText(name))
	general.WaitSpinner.Prefix = actionPrint
	general.WaitSpinner.Start()

	// Pull 前检测本地存储库是否存在
	if general.FileExist(path) {
		isRepo, repo, headRef := general.IsLocalRepo(path)
		if isRepo { // 本地存储库可以 Pull
			// 开始 Pull
			worktree, leftCommit, rightCommit, err := general.PullRepo(repo, publicKeys)
			// Pull 结束
			if err != nil {
				if err == git.NoErrAlreadyUpToDate {
					// 本地存储库已经是最新
					general.WaitSpinner.Stop()
					color.Printf("%s%s %s %s\n", actionPrint, general.FgBlueText(general.LatestFlag), general.SecondaryText("Already up-to-date"), general.SecondaryText("[", headRef.Name().Short(), "]"))

					// 获取子模块信息
					submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
					if err != nil {
						general.WaitSpinner.Stop()
						color.Printf("%s", actionPrint)
						fileName, lineNo := general.GetCallerInfo()
						color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
						return
					}
					if len(submodules) != 0 {
						length := len(general.RunFlag) + len("Pulling") // 子模块缩进长度
						for index, submodule := range submodules {
							joiner := func() string { // 主模块和子模块的输出连接符
								if index == len(submodules)-1 {
									return general.JoinerFinish
								}
								return general.JoinerIng
							}()
							// 开始 Pull 提示
							subActionPrint := color.Sprintf("%s%s %s %s: ", strings.Repeat(" ", length), joiner, general.SubmoduleFlag, general.FgMagentaText(submodule.Config().Name))
							general.WaitSpinner.Prefix = subActionPrint
							general.WaitSpinner.Start()
							submoduleRepo, err := submodule.Repository()
							if err != nil {
								general.WaitSpinner.Stop()
								color.Printf("%s", subActionPrint)
								fileName, lineNo := general.GetCallerInfo()
								color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
							} else {
								// 开始 Pull
								submoduleRepoHeadRef := general.GetRepoHeadRef(submoduleRepo)
								_, submoduleLeftCommit, submoduleRightCommit, err := general.PullRepo(submoduleRepo, publicKeys)
								// Pull 结束
								if err != nil {
									if err == git.NoErrAlreadyUpToDate {
										general.WaitSpinner.Stop()
										color.Printf("%s%s %s %s", subActionPrint, general.FgBlueText(general.LatestFlag), general.SecondaryText("Already up-to-date"), general.SecondaryText("[", submoduleRepoHeadRef.Name().Short(), "]"))
									} else {
										general.WaitSpinner.Stop()
										color.Printf("%s", subActionPrint)
										fileName, lineNo := general.GetCallerInfo()
										color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
									}
								} else {
									general.WaitSpinner.Stop()
									color.Printf("%s%s %s --> %s %s", subActionPrint, general.SuccessFlag, general.FgBlueText(submoduleLeftCommit.Hash.String()[:6]), general.FgGreenText(submoduleRightCommit.Hash.String()[:6]), general.SecondaryText("[", submoduleRepoHeadRef.Name().Short(), "]"))
								}
							}
							color.Println() // 当前子模块处理完成，处理下一个子模块
						}
					}
				} else {
					general.WaitSpinner.Stop()
					color.Printf("%s", actionPrint)
					fileName, lineNo := general.GetCallerInfo()
					color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
				}
			} else {
				// 成功 Pull
				general.WaitSpinner.Stop()
				color.Printf("%s%s %s --> %s %s\n", actionPrint, general.SuccessFlag, general.FgBlueText(leftCommit.Hash.String()[:6]), general.FgGreenText(rightCommit.Hash.String()[:6]), general.SecondaryText("[", headRef.Name().Short(), "]"))

				// 尝试 Pull 子模块
				submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
				if err != nil {
					general.WaitSpinner.Stop()
					color.Printf("%s", actionPrint)
					fileName, lineNo := general.GetCallerInfo()
					color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
					return
				}
				if len(submodules) != 0 {
					length := len(general.RunFlag) + len("Pulling") // 子模块缩进长度
					for index, submodule := range submodules {
						joiner := func() string { // 主模块和子模块的输出连接符
							if index == len(submodules)-1 {
								return general.JoinerFinish
							}
							return general.JoinerIng
						}()
						// 开始 Pull 提示
						subActionPrint := color.Sprintf("%s%s %s %s: ", strings.Repeat(" ", length), joiner, general.SubmoduleFlag, general.FgMagentaText(submodule.Config().Name))
						general.WaitSpinner.Prefix = subActionPrint
						general.WaitSpinner.Start()
						submoduleRepo, err := submodule.Repository()
						if err != nil {
							general.WaitSpinner.Stop()
							color.Printf("%s", subActionPrint)
							fileName, lineNo := general.GetCallerInfo()
							color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
						} else {
							// 开始 Pull
							submoduleRepoHeadRef := general.GetRepoHeadRef(submoduleRepo)
							_, submoduleLeftCommit, submoduleRightCommit, err := general.PullRepo(submoduleRepo, publicKeys)
							// Pull 结束
							if err != nil {
								if err == git.NoErrAlreadyUpToDate {
									general.WaitSpinner.Stop()
									color.Printf("%s%s %s %s", subActionPrint, general.FgBlueText(general.LatestFlag), general.SecondaryText("Already up-to-date"), general.SecondaryText("[", submoduleRepoHeadRef.Name().Short(), "]"))
								} else {
									general.WaitSpinner.Stop()
									color.Printf("%s", subActionPrint)
									fileName, lineNo := general.GetCallerInfo()
									color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
								}
							} else {
								general.WaitSpinner.Stop()
								color.Printf("%s%s %s --> %s %s", subActionPrint, general.SuccessFlag, general.FgBlueText(submoduleLeftCommit.Hash.String()[:6]), general.FgGreenText(submoduleRightCommit.Hash.String()[:6]), general.SecondaryText("[", submoduleRepoHeadRef.Name().Short(), "]"))
							}
						}
						color.Println() // 当前子模块处理完成，处理下一个子模块
					}
				}
			}
		} else { // 非本地存储库无法 Pull
			general.WaitSpinner.Stop()
			color.Printf("%s%s %s\n", actionPrint, general.ErrorFlag, general.DangerText("Folder is not a local repository"))
		}
	} else {
		general.WaitSpinner.Stop()
		color.Printf("%s%s %s\n", actionPrint, general.ErrorFlag, general.DangerText("The local repository does not exist"))
	}

}
