/*
File: git_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-11 14:44:19

Description: git操作
*/

package function

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

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

// 检测是不是本地仓库，是的话返回*git.Repository对象
func IsLocalRepo(path string) (bool, *git.Repository) {
	// 能打开就是本地仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, nil
	}
	return true, repo
}

// 输出本地仓库[本地|远程]分支信息
func GetRepoBranchInfo(worktree *git.Worktree, which string) ([]fs.FileInfo, error) {
	var branchDir string
	switch which {
	case "local":
		branchDir = ".git/refs/heads"
	case "remote":
		branchDir = ".git/refs/remotes/origin"
	default:
		return nil, fmt.Errorf("Parameter error: %s", "optional value of which is 'local' or 'remote'")
	}
	branchs, err := worktree.Filesystem.ReadDir(branchDir)
	if err != nil {
		return nil, err
	}

	return branchs, nil
}

// 输出本地仓库子模块信息
func GetLocalRepoSubmoduleInfo(worktree *git.Worktree) (git.Submodules, error) {
	submodules, err := worktree.Submodules()
	if err != nil {
		return nil, err
	}

	return submodules, nil
}
