/*
File: pull.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-03-05 14:22:54

Description: 子命令 `pull` 的实现
*/

package cli

import (
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
	"github.com/yhyj/curator/general"
)

// RollingPullRepos 遍历拉取远端仓库的更改到本地
//
// 参数：
//   - confile: 程序配置文件
//   - source: 远端仓库源，支持 'github' 和 'gitea'，默认为 'github'
func RollingPullRepos(confile, source string) {
	// 加载配置文件
	conf, err := GetTomlConfig(confile)
	if err != nil {
		color.Error.Println(err)
	} else {
		// 获取配置项
		pemfile := conf.Get("ssh.rsa_file")
		storagePath := conf.Get("storage.path").(string)
		repoNames := conf.Get("git.repos").([]interface{})
		// 获取公钥
		publicKeys, err := general.GetPublicKeysByGit(pemfile.(string))
		if err != nil {
			color.Error.Println(err)
			return
		}

		// 创建运行状态符号
		// 拉取
		color.Info.Tips("%s %s\n", general.FgWhite("Fetch from and merge with"), general.FgGreen(source))
		for _, repoName := range repoNames {
			repoPath := filepath.Join(storagePath, repoName.(string))
			// 开始拉取
			color.Printf("%s %s %s: ", general.FgGreen(general.Run), general.LightText("Pulling"), general.FgCyan(repoName.(string)))
			// 拉取前检测本地仓库是否存在
			if general.FileExist(repoPath) {
				isRepo, repo := general.IsLocalRepo(repoPath)
				if isRepo {
					worktree, leftCommit, rightCommit, err := general.PullRepo(repo, publicKeys)
					if err != nil {
						if err == git.NoErrAlreadyUpToDate {
							color.Printf("%s %s\n", general.FgBlue(general.Dot), general.SecondaryText("Already up-to-date"))
							// 尝试拉取子模块
							submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
							if err != nil {
								color.Error.Println(err)
								continue
							}
							if len(submodules) != 0 {
								length := len(general.Run) + len("Pulling") // 子模块缩进长度
								for index, submodule := range submodules {
									// 创建和主模块的连接符
									joiner := func() string {
										if index == len(submodules)-1 {
											return general.JoinerFinish
										}
										return general.JoinerIng
									}()
									color.Printf("%s%s %s %s: ", strings.Repeat(" ", length), joiner, "📦", general.FgMagenta(submodule.Config().Name))
									submoduleRepo, err := submodule.Repository()
									if err != nil {
										color.Error.Println(err)
									} else {
										_, submoduleLeftCommit, submoduleRightCommit, err := general.PullRepo(submoduleRepo, publicKeys)
										if err != nil {
											if err == git.NoErrAlreadyUpToDate {
												color.Printf("%s %s", general.FgBlue(general.Dot), general.SecondaryText("Already up-to-date"))
											} else {
												color.Error.Println(err)
											}
										} else {
											color.Printf("%s %s %s", general.FgBlue(submoduleLeftCommit.Hash.String()[:6]), general.LightText("-->"), general.FgGray(submoduleRightCommit.Hash.String()[:6]))
										}
									}
									color.Println() // 子模块处理完成，换行
								}
							}
						} else {
							color.Error.Println(err)
						}
					} else {
						color.Printf("%s %s %s %s\n", general.Yes, general.FgBlue(leftCommit.Hash.String()[:6]), general.LightText("-->"), general.FgGray(rightCommit.Hash.String()[:6]))
						// 尝试拉取子模块
						submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
						if err != nil {
							color.Error.Println(err)
							continue
						}
						if len(submodules) != 0 {
							length := len(general.Run) + len("Pulling") // 子模块缩进长度
							for index, submodule := range submodules {
								// 创建和主模块的连接符
								joiner := func() string {
									if index == len(submodules)-1 {
										return general.JoinerFinish
									}
									return general.JoinerIng
								}()
								color.Printf("%s%s %s %s: ", strings.Repeat(" ", length), joiner, "📦", general.FgMagenta(submodule.Config().Name))
								submoduleRepo, err := submodule.Repository()
								if err != nil {
									color.Error.Println(err)
								} else {
									_, submoduleLeftCommit, submoduleRightCommit, err := general.PullRepo(submoduleRepo, publicKeys)
									if err != nil {
										if err == git.NoErrAlreadyUpToDate {
											color.Printf("%s %s", general.FgBlue(general.Dot), general.SecondaryText("Already up-to-date"))
										} else {
											color.Error.Println(err)
										}
									} else {
										color.Printf("%s %s %s", general.FgBlue(submoduleLeftCommit.Hash.String()[:6]), general.LightText("-->"), general.FgGray(submoduleRightCommit.Hash.String()[:6]))
									}
								}
								color.Println() // 子模块处理完成，换行
							}
						}
					}
				} else { // 非本地仓库
					color.Printf("%s %s\n", general.FgRed(general.No), general.ErrorText("Folder is not a local repository"))
				}
			} else {
				color.Printf("%s %s\n", general.FgRed(general.No), general.ErrorText("The local repository does not exist"))
			}
			// 添加一个延时，使输出更加顺畅
			general.Delay(0.1)
		}
	}
}
