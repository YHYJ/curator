/*
File: git_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-11 14:44:19

Description: git操作
*/

package general

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
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

// 获取本地仓库[本地|远程]分支信息
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

// 本地仓库根据远程分支refs/remotes/origin/<remoteBranchName>创建本地分支refs/heads/<localBranchName>
func CreateLocalBranch(repo *git.Repository, branchs []fs.FileInfo) []string {
	var errList []string //使用一个Slice存储所有错误信息以美化输出
	for _, branch := range branchs {
		// 修改.git/config，增加新的分支配置
		branchReferenceName := plumbing.NewBranchReferenceName(branch.Name()) // 构建本地分支Reference名，格式：refs/heads/<localBranchName>
		repo.CreateBranch(&config.Branch{                                     // 分支配置写入.git/config
			Name:   branch.Name(),
			Remote: "origin",
			Merge:  branchReferenceName,
		})

		// 创建一个新的Reference
		remote := "origin"                                                             // 远程名称
		newBranchReferenceName := plumbing.ReferenceName(branchReferenceName.String()) // refs/heads/test
		remoteReferenceName := plumbing.NewRemoteReferenceName(remote, branch.Name())  // 构建远程分支Reference名，格式：refs/remotes/origin/<remoteBranchName>
		remoteReferenceData, err := repo.Reference(remoteReferenceName, true)          // 根据远程分支Reference名获取其Hash值，格式：1a8f900411d35a620407ce07902aecadfc782ded refs/remotes/origin/test
		if err != nil {
			errList = append(errList, err.Error())
			continue
		}
		newReference := plumbing.NewHashReference(newBranchReferenceName, remoteReferenceData.Hash()) // 基于Hash创建新的Reference
		if err = repo.Storer.SetReference(newReference); err != nil {                                 // 写入新Reference
			errList = append(errList, err.Error())
			continue
		}
	}
	return errList
}

// 获取本地仓库子模块信息
func GetLocalRepoSubmoduleInfo(worktree *git.Worktree) (git.Submodules, error) {
	submodules, err := worktree.Submodules()
	if err != nil {
		return nil, err
	}

	return submodules, nil
}
