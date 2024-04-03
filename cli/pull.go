/*
File: pull.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-03-05 14:22:54

Description: 子命令 `pull` 的实现
*/

package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
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
		fmt.Printf(general.ErrorBaseFormat, err)
	} else {
		// 获取配置项
		pemfile := conf.Get("ssh.rsa_file")
		storagePath := conf.Get("storage.path").(string)
		repoNames := conf.Get("git.repos").([]interface{})
		// 获取公钥
		publicKeys, err := general.GetPublicKeysByGit(pemfile.(string))
		if err != nil {
			fmt.Printf(general.ErrorBaseFormat, err)
			return
		}

		// 创建运行状态符号
		yesSymbol := fmt.Sprintf("%s%s%s", "[", general.Yes, "]")
		noSymbol := fmt.Sprintf("%s%s%s", "[", general.No, "]")
		// 拉取
		fmt.Printf(general.TipsPrefixFormat, "Fetch from and merge with", " ", source)
		fmt.Println()
		for _, repoName := range repoNames {
			repoPath := filepath.Join(storagePath, repoName.(string))
			// 开始拉取
			fmt.Printf(general.Tips2PSuffixNoNewLineFormat, general.Run, " Pulling ", repoName.(string), ":", " ")
			// 拉取前检测本地仓库是否存在
			if general.FileExist(repoPath) {
				isRepo, repo := general.IsLocalRepo(repoPath)
				if isRepo {
					worktree, leftCommit, rightCommit, err := general.PullRepo(repo, publicKeys)
					if err != nil {
						if err == git.NoErrAlreadyUpToDate {
							fmt.Printf(general.SliceTraverse2PFormat, yesSymbol, " ", "Already up-to-date")
							// 尝试拉取子模块
							submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
							if err != nil {
								fmt.Printf(general.ErrorBaseFormat, err)
								continue
							}
							if len(submodules) != 0 {
								length := len(general.Run) + len("Pulling") // 子模块缩进长度
								for index, submodule := range submodules {
									fmt.Printf(strings.Repeat(" ", length)) // 子模块信息相对主模块进行一次缩进
									// 创建和主模块的连接符
									joiner := fmt.Sprintf("%s%s", general.JoinerIng, " ")
									if index == len(submodules)-1 {
										joiner = fmt.Sprintf("%s%s", general.JoinerFinish, " ")
									}
									fmt.Printf(general.InfoPrefixSuffixNoNewLineFormat, joiner, "[", submodule.Config().Name, "]", "")
									submoduleRepo, err := submodule.Repository()
									if err != nil {
										fmt.Printf(general.ErrorBaseFormat, err)
									} else {
										_, submoduleLeftCommit, submoduleRightCommit, err := general.PullRepo(submoduleRepo, publicKeys)
										if err != nil {
											if err == git.NoErrAlreadyUpToDate {
												fmt.Printf(general.SliceTraverse2PNoNewLineFormat, yesSymbol, " ", "Already up-to-date")
											} else {
												fmt.Printf(general.ErrorBaseFormat, err)
											}
										} else {
											fmt.Printf(general.SliceTraverse2PSuffixNoNewLineFormat, submoduleLeftCommit.Hash.String()[:6], " --> ", submoduleRightCommit.Hash.String()[:6], "")
										}
									}
									fmt.Println() // 子模块处理完成，换行
								}
							}
						} else {
							fmt.Printf(general.ErrorBaseFormat, err)
						}
					} else {
						fmt.Printf(general.SuccessSuffixNoNewLineFormat, yesSymbol, " ", "")
						fmt.Printf(general.SliceTraverse2PSuffixFormat, leftCommit.Hash.String()[:6], " --> ", rightCommit.Hash.String()[:6], "")
						// 尝试拉取子模块
						submodules, err := general.GetLocalRepoSubmoduleInfo(worktree)
						if err != nil {
							fmt.Printf(general.ErrorBaseFormat, err)
							continue
						}
						if len(submodules) != 0 {
							length := len(general.Run) + len("Pulling") // 子模块缩进长度
							for index, submodule := range submodules {
								fmt.Printf(strings.Repeat(" ", length)) // 子模块信息相对主模块进行一次缩进
								// 创建和主模块的连接符
								joiner := fmt.Sprintf("%s%s", general.JoinerIng, " ")
								if index == len(submodules)-1 {
									joiner = fmt.Sprintf("%s%s", general.JoinerFinish, " ")
								}
								fmt.Printf(general.InfoPrefixSuffixNoNewLineFormat, joiner, "[", submodule.Config().Name, "]", "")
								submoduleRepo, err := submodule.Repository()
								if err != nil {
									fmt.Printf(general.ErrorBaseFormat, err)
								} else {
									_, submoduleLeftCommit, submoduleRightCommit, err := general.PullRepo(submoduleRepo, publicKeys)
									if err != nil {
										if err == git.NoErrAlreadyUpToDate {
											fmt.Printf(general.SliceTraverse2PNoNewLineFormat, yesSymbol, " ", "Already up-to-date")
										} else {
											fmt.Printf(general.ErrorBaseFormat, err)
										}
									} else {
										fmt.Printf(general.SliceTraverse2PSuffixNoNewLineFormat, submoduleLeftCommit.Hash.String()[:6], " --> ", submoduleRightCommit.Hash.String()[:6], "")
									}
								}
								fmt.Println() // 子模块处理完成，换行
							}
						}
					}
				} else { // 非本地仓库
					fmt.Printf(general.ErrorSuffixFormat, noSymbol, " ", "Folder is not a local repository")
				}
			} else {
				fmt.Printf(general.ErrorSuffixFormat, noSymbol, " ", "The local repository does not exist")
			}
			// 添加一个延时，使输出更加顺畅
			general.Delay(0.1)
		}
	}
}
