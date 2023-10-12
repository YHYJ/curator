/*
File: git_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-11 14:44:19

Description: git操作
*/

package function

import (
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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

// 使用SSH协议将远端仓库克隆到本地
func CloneRepoViaSSH(repoPath, URL, username, repoName string, publicKeys *ssh.PublicKeys) (*git.Repository, error) {
	repoUrl := "git" + "@" + URL + ":" + username + "/" + repoName + ".git"
	repo, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:               repoUrl,
		Auth:              publicKeys,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          io.Discard, // os.Stdout会将Clone的详细过程输出到控制台，io.Discard会直接丢弃
	})

	return repo, err
}
