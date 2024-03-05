/*
File: define_git.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-11 14:44:19

Description: git 操作
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

// CloneRepoViaSSH 使用 SSH 协议将远端仓库克隆到本地
//
// 参数：
//   - repoPath: 本地仓库路径
//   - URL: 远端仓库地址（仅包括主地址，例如：github.com）
//   - username: 远端仓库用户名
//   - repoName: 远端仓库名称
//   - publicKeys: ssh 公钥
//
// 返回：
//   - 本地仓库对象
//   - 错误信息
func CloneRepoViaSSH(repoPath, URL, username, repoName string, publicKeys *ssh.PublicKeys) (*git.Repository, error) {
	repoUrl := "git" + "@" + URL + ":" + username + "/" + repoName + ".git"
	repo, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:               repoUrl,
		Auth:              publicKeys,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          io.Discard, // os.Stdout 会将 Clone 的详细过程输出到控制台，io.Discard 会直接丢弃
	})

	return repo, err
}

// IsLocalRepo 检测是不是本地仓库，是的话返回本地仓库对象
//
// 参数：
//   - path: 本地仓库路径
//
// 返回：
//   - 是否本地仓库
//   - 本地仓库对象
func IsLocalRepo(path string) (bool, *git.Repository) {
	// 能打开就是本地仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, nil
	}
	return true, repo
}

// GetRepoBranchInfo 获取本地仓库的[本地|远程]分支信息
//
// 参数：
//   - worktree: 仓库的 git 工作树对象
//   - which: 'local' or 'remote'，指定要获取的是本地分支还是远程分支
//
// 返回：
//   - 分支信息
//   - 错误信息
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

// CreateLocalBranch 本地仓库根据远程分支创建本地分支
//
//   - 远程分支 refs/remotes/origin/<remoteBranchName>
//   - 本地分支 refs/heads/<localBranchName>
//
// 参数：
//   - repo: 本地仓库对象
//   - branchs: 远程分支信息
//
// 返回：
//   - 错误信息切片
func CreateLocalBranch(repo *git.Repository, branchs []fs.FileInfo) []string {
	var errList []string  // 使用一个 Slice 存储所有错误信息以美化输出
	var remote = "origin" // 远程名称
	for _, branch := range branchs {
		// 修改 .git/config ，增加新的分支配置
		branchReferenceName := plumbing.NewBranchReferenceName(branch.Name()) // 构建本地分支 Reference 名，格式： refs/heads/<localBranchName>
		repo.CreateBranch(&config.Branch{                                     // 分支配置写入 .git/config
			Name:   branch.Name(),
			Remote: remote,
			Merge:  branchReferenceName,
		})

		// 创建一个新的 Reference
		newBranchReferenceName := plumbing.ReferenceName(branchReferenceName.String()) // refs/heads/test
		remoteReferenceName := plumbing.NewRemoteReferenceName(remote, branch.Name())  // 构建远程分支 Reference 名，格式： refs/remotes/origin/<remoteBranchName>
		remoteReferenceData, err := repo.Reference(remoteReferenceName, true)          // 根据远程分支 Reference 名获取其 Hash 值，格式： 1a8f900411d35a620407ce07902aecadfc782ded refs/remotes/origin/test
		if err != nil {
			errList = append(errList, err.Error())
			continue
		}
		newReference := plumbing.NewHashReference(newBranchReferenceName, remoteReferenceData.Hash()) // 基于 Hash 创建新的 Reference
		if err = repo.Storer.SetReference(newReference); err != nil {                                 // 写入新 Reference
			errList = append(errList, err.Error())
			continue
		}
	}
	return errList
}

// GetLocalRepoSubmoduleInfo 获取本地仓库子模块信息
//
// 参数：
//   - worktree: 仓库的 git 工作树对象
//
// 返回：
//   - 子模块信息
//   - 错误信息
func GetLocalRepoSubmoduleInfo(worktree *git.Worktree) (git.Submodules, error) {
	submodules, err := worktree.Submodules()
	if err != nil {
		return nil, err
	}

	return submodules, nil
}
