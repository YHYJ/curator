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
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
	"github.com/pelletier/go-toml"
	"github.com/yhyj/curator/general"
)

// RollingPullRepos 遍历拉取远端仓库的更改到本地
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
//   - source: 远端仓库源，支持 'github' 和 'gitea'，默认为 'github'
func RollingPullRepos(configTree *toml.Tree, source string) {
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

	// 为已 Clone 存储库计数
	totalNum := len(config.Git.Repos) // 总存储库数
	clonedNum := 0                    // 已 Clone 存储库数
	for _, repoName := range config.Git.Repos {
		repoPath := filepath.Join(config.Storage.Path, repoName) // 本地存储库路径
		if general.FileExist(repoPath) {
			isRepo, _, _ := general.IsLocalRepo(repoPath)
			if isRepo {
				clonedNum++
			}
		}
	}

	// 开始 Pull 提示
	negatives := strings.Builder{}
	negatives.WriteString(color.Sprintf("%s Pull repository from %s: %d/%d cloned\n", general.InfoText("INFO:"), general.FgGreenText(source), clonedNum, totalNum))
	negatives.WriteString(color.Sprintf("%s Repository root: %s\n", general.InfoText("INFO:"), general.PrimaryText(config.Storage.Path)))

	// 让用户选择需要 Pull 的存储库
	selectedRepos, err := general.MultipleSelectionFilter(config.Git.Repos, negatives.String())
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s -> Unable to start selector: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		return
	}

	// 遍历所选存储库名
	for _, repoName := range selectedRepos {
		repoPath := filepath.Join(config.Storage.Path, repoName)
		// 开始拉取提示
		actionPrint := color.Sprintf("%s Pulling %s: ", general.RunFlag, general.FgCyanText(repoName))
		general.WaitSpinner.Prefix = actionPrint
		general.WaitSpinner.Start()
		// 拉取前检测本地仓库是否存在
		if general.FileExist(repoPath) {
			isRepo, repo, headRef := general.IsLocalRepo(repoPath)
			if isRepo { // 本地仓库可以拉取
				// 开始拉取
				worktree, leftCommit, rightCommit, err := general.PullRepo(repo, publicKeys)
				// 拉取结束
				if err != nil {
					if err == git.NoErrAlreadyUpToDate {
						// 本地仓库已经是最新
						general.WaitSpinner.Stop()
						color.Printf("%s%s %s %s\n", actionPrint, general.FgBlueText(general.LatestFlag), general.SecondaryText("Already up-to-date"), general.SecondaryText("[", headRef.Name().Short(), "]"))

						// 尝试拉取子模块
						submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
						if err != nil {
							general.WaitSpinner.Stop()
							color.Printf("%s", actionPrint)
							fileName, lineNo := general.GetCallerInfo()
							color.Printf("%s %s -> Unable to get local submodule info: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
							continue
						}
						if len(submodules) != 0 {
							length := len(general.RunFlag) + len("Pulling") // 子模块缩进长度
							for index, submodule := range submodules {
								// 创建和主模块的连接符
								joiner := func() string {
									if index == len(submodules)-1 {
										return general.JoinerFinish
									}
									return general.JoinerIng
								}()
								// 开始拉取提示
								subActionPrint := color.Sprintf("%s%s %s %s: ", strings.Repeat(" ", length), joiner, general.SubmoduleFlag, general.FgMagentaText(submodule.Config().Name))
								general.WaitSpinner.Prefix = subActionPrint
								general.WaitSpinner.Start()
								submoduleRepo, err := submodule.Repository()
								if err != nil {
									general.WaitSpinner.Stop()
									color.Printf("%s", subActionPrint)
									fileName, lineNo := general.GetCallerInfo()
									color.Printf("%s %s -> Unable to get local submodule: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
								} else {
									// 开始拉取
									submoduleRepoHeadRef := general.GetRepoHeadRef(submoduleRepo)
									_, submoduleLeftCommit, submoduleRightCommit, err := general.PullRepo(submoduleRepo, publicKeys)
									// 拉取结束
									if err != nil {
										if err == git.NoErrAlreadyUpToDate {
											general.WaitSpinner.Stop()
											color.Printf("%s%s %s %s", subActionPrint, general.FgBlueText(general.LatestFlag), general.SecondaryText("Already up-to-date"), general.SecondaryText("[", submoduleRepoHeadRef.Name().Short(), "]"))
										} else {
											general.WaitSpinner.Stop()
											color.Printf("%s", subActionPrint)
											fileName, lineNo := general.GetCallerInfo()
											color.Printf("%s %s -> Unable to pull submodule: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
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
						color.Printf("%s %s -> Unable to pull repository: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
					}
				} else {
					// 成功拉取
					general.WaitSpinner.Stop()
					color.Printf("%s%s %s --> %s %s\n", actionPrint, general.SuccessFlag, general.FgBlueText(leftCommit.Hash.String()[:6]), general.FgGreenText(rightCommit.Hash.String()[:6]), general.SecondaryText("[", headRef.Name().Short(), "]"))

					// 尝试拉取子模块
					submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
					if err != nil {
						general.WaitSpinner.Stop()
						color.Printf("%s", actionPrint)
						fileName, lineNo := general.GetCallerInfo()
						color.Printf("%s %s -> Unable to get local submodule info: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
						continue
					}
					if len(submodules) != 0 {
						length := len(general.RunFlag) + len("Pulling") // 子模块缩进长度
						for index, submodule := range submodules {
							// 创建和主模块的连接符
							joiner := func() string {
								if index == len(submodules)-1 {
									return general.JoinerFinish
								}
								return general.JoinerIng
							}()
							// 开始拉取提示
							subActionPrint := color.Sprintf("%s%s %s %s: ", strings.Repeat(" ", length), joiner, general.SubmoduleFlag, general.FgMagentaText(submodule.Config().Name))
							general.WaitSpinner.Prefix = subActionPrint
							general.WaitSpinner.Start()
							submoduleRepo, err := submodule.Repository()
							if err != nil {
								general.WaitSpinner.Stop()
								color.Printf("%s", subActionPrint)
								fileName, lineNo := general.GetCallerInfo()
								color.Printf("%s %s -> Unable to get local repository: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
							} else {
								// 开始拉取
								submoduleRepoHeadRef := general.GetRepoHeadRef(submoduleRepo)
								_, submoduleLeftCommit, submoduleRightCommit, err := general.PullRepo(submoduleRepo, publicKeys)
								// 拉取结束
								if err != nil {
									if err == git.NoErrAlreadyUpToDate {
										general.WaitSpinner.Stop()
										color.Printf("%s%s %s %s", subActionPrint, general.FgBlueText(general.LatestFlag), general.SecondaryText("Already up-to-date"), general.SecondaryText("[", submoduleRepoHeadRef.Name().Short(), "]"))
									} else {
										general.WaitSpinner.Stop()
										color.Printf("%s", subActionPrint)
										fileName, lineNo := general.GetCallerInfo()
										color.Printf("%s %s -> Unable to pull submodule: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
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
			} else { // 非本地仓库无法拉取
				general.WaitSpinner.Stop()
				color.Printf("%s%s %s\n", actionPrint, general.ErrorFlag, general.DangerText("Folder is not a local repository"))
			}
		} else {
			general.WaitSpinner.Stop()
			color.Printf("%s%s %s\n", actionPrint, general.ErrorFlag, general.DangerText("The local repository does not exist"))
		}
		// 添加一个延时，使输出更加顺畅
		general.Delay(0.1)
	}
}
