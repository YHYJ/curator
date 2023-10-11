/*
File: git_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-11 14:44:19

Description: git操作
*/

package function

import (
	"github.com/go-git/go-git/v5"
)

// 获取git仓库的所有子模块名称
func GetSubModuleNames(repoPath string) (git.Submodules, error) {
	// 打开本地存储库
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	submodules, err := worktree.Submodules()
	if err != nil {
		return nil, err
	}

	return submodules, nil
}
