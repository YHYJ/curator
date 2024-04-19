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
	// pemfile := conf.Get("ssh.rsa_file")
	// storagePath := conf.Get("storage.path").(string)
	// repoNames := conf.Get("git.repos").([]interface{})
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

	// 拉取
	color.Info.Tips("%s %s", general.FgWhiteText("Fetch from and merge with"), general.FgGreenText(source))
	color.Info.Tips("%s: %s", general.FgWhiteText("Repository root"), general.PrimaryText(config.Storage.Path))
	// 让用户选择需要 Pull 的存储库
	selectedRepos, err := general.MultipleSelectionFilter(config.Git.Repos)
	if err != nil {
		color.Error.Println(err)
	}
	// 对所选的存储库进行排序
	sort.Strings(selectedRepos)
	// 遍历所选存储库名
	for _, repoName := range selectedRepos {
		repoPath := filepath.Join(config.Storage.Path, repoName)
		// 开始拉取
		color.Printf("%s %s %s: ", general.RunFlag, general.LightText("Pulling"), general.FgCyanText(repoName))
		// 拉取前检测本地仓库是否存在
		if general.FileExist(repoPath) {
			isRepo, repo := general.IsLocalRepo(repoPath)
			if isRepo {
				worktree, leftCommit, rightCommit, err := general.PullRepo(repo, publicKeys)
				if err != nil {
					if err == git.NoErrAlreadyUpToDate {
						color.Printf("%s %s\n", general.FgBlueText(general.LatestFlag), general.SecondaryText("Already up-to-date"))
						// 尝试拉取子模块
						submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
						if err != nil {
							color.Error.Println(err)
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
								color.Printf("%s%s %s %s: ", strings.Repeat(" ", length), joiner, general.SubmoduleFlag, general.FgMagentaText(submodule.Config().Name))
								submoduleRepo, err := submodule.Repository()
								if err != nil {
									color.Error.Println(err)
								} else {
									_, submoduleLeftCommit, submoduleRightCommit, err := general.PullRepo(submoduleRepo, publicKeys)
									if err != nil {
										if err == git.NoErrAlreadyUpToDate {
											color.Printf("%s %s", general.FgBlueText(general.LatestFlag), general.SecondaryText("Already up-to-date"))
										} else {
											color.Error.Println(err)
										}
									} else {
										color.Printf("%s %s %s %s", general.SuccessFlag, general.FgBlueText(submoduleLeftCommit.Hash.String()[:6]), general.LightText("-->"), general.FgGreenText(submoduleRightCommit.Hash.String()[:6]))
									}
								}
								color.Println() // 子模块处理完成，换行
							}
						}
					} else {
						color.Error.Println(err)
					}
				} else {
					color.Printf("%s %s %s %s\n", general.SuccessFlag, general.FgBlueText(leftCommit.Hash.String()[:6]), general.LightText("-->"), general.FgGreenText(rightCommit.Hash.String()[:6]))
					// 尝试拉取子模块
					submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
					if err != nil {
						color.Error.Println(err)
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
							color.Printf("%s%s %s %s: ", strings.Repeat(" ", length), joiner, general.SubmoduleFlag, general.FgMagentaText(submodule.Config().Name))
							submoduleRepo, err := submodule.Repository()
							if err != nil {
								color.Error.Println(err)
							} else {
								_, submoduleLeftCommit, submoduleRightCommit, err := general.PullRepo(submoduleRepo, publicKeys)
								if err != nil {
									if err == git.NoErrAlreadyUpToDate {
										color.Printf("%s %s", general.FgBlueText(general.LatestFlag), general.SecondaryText("Already up-to-date"))
									} else {
										color.Error.Println(err)
									}
								} else {
									color.Printf("%s %s %s %s", general.SuccessFlag, general.FgBlueText(submoduleLeftCommit.Hash.String()[:6]), general.LightText("-->"), general.FgGreenText(submoduleRightCommit.Hash.String()[:6]))
								}
							}
							color.Println() // 子模块处理完成，换行
						}
					}
				}
			} else { // 非本地仓库
				color.Printf("%s %s\n", general.ErrorFlag, general.ErrorText("Folder is not a local repository"))
			}
		} else {
			color.Printf("%s %s\n", general.ErrorFlag, general.ErrorText("The local repository does not exist"))
		}
		// 添加一个延时，使输出更加顺畅
		general.Delay(0.1)
	}
}
