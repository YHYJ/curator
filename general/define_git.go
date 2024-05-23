/*
File: define_git.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-11 14:44:19

Description: git 操作
*/

package general

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/gookit/color"
)

var remoteName = "origin" // 远程名称

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

// PullRepo 拉取远端仓库的更改到本地
//
// 参数：
//   - repo: 本地仓库对象
//
// 返回：
//   - 仓库的 git 工作树对象
//   - 拉取前本地最新 Commit 的 Hash 值
//   - 拉取后本地最新 Commit 的 Hash 值
//   - 错误信息
func PullRepo(repo *git.Repository, publicKeys *ssh.PublicKeys) (worktree *git.Worktree, leftCommit, rightCommit *object.Commit, err error) {
	// 获取本地仓库的 worktree
	worktree, err = repo.Worktree()
	if err != nil {
		return nil, nil, nil, err
	}

	// 获取拉取前的最新 Commit 的 Hash 值
	leftRef, err := repo.Head()
	if err != nil {
		return worktree, nil, nil, err
	}
	leftCommit, err = repo.CommitObject(leftRef.Hash())
	if err != nil {
		return worktree, nil, nil, err
	}

	// 拉取远端仓库的更改
	err = worktree.Pull(&git.PullOptions{
		Auth:          publicKeys,
		RemoteName:    remoteName,
		ReferenceName: leftRef.Name(),
	})
	if err != nil {
		return worktree, nil, nil, err
	}

	// 获取拉取后的最新 Commit 的 Hash 值
	RightRef, err := repo.Head()
	if err != nil {
		return worktree, nil, nil, err
	}
	rightCommit, err = repo.CommitObject(RightRef.Hash())
	if err != nil {
		return worktree, nil, nil, err
	}

	return worktree, leftCommit, rightCommit, nil
}

// IsLocalRepo 检测是不是本地仓库，是的话返回本地仓库对象及其 HEAD 指向的引用
//
// 参数：
//   - path: 本地仓库路径
//
// 返回：
//   - 是否本地仓库
//   - 本地仓库对象
//   - HEAD 引用
func IsLocalRepo(path string) (bool, *git.Repository, *plumbing.Reference) {
	// 能打开就是本地仓库
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, nil, nil
	}

	// 获取 HEAD 引用
	headRef := GetRepoHeadRef(repo)

	return true, repo, headRef
}

// GetRepoHeadRef 获取本地仓库对象 HEAD 指向的引用
//
// 参数：
//   - repo: 本地仓库对象
//
// 返回：
//   - HEAD 引用
func GetRepoHeadRef(repo *git.Repository) *plumbing.Reference {
	// 获取 HEAD 引用
	headRef, err := repo.Head()
	if err != nil {
		return nil
	}

	return headRef
}

// GetRepoBranchInfo 获取本地仓库的[本地|远程]分支信息
//
// 参数：
//   - worktree: 仓库的 git 工作树对象
//   - isSubmodule: 该仓库是否作为子模块
//   - submoduleName: 当该仓库作为子模块时需要仓库名
//   - which: 'local' or 'remote'，指定要获取的是本地分支还是远程分支
//
// 返回：
//   - 分支信息
//   - 错误信息
func GetRepoBranchInfo(worktree *git.Worktree, isSubmodule bool, submoduleName string, which string) ([]fs.FileInfo, error) {
	var branchDir string
	switch which {
	case "local":
		switch isSubmodule {
		case false:
			branchDir = ".git/refs/heads"
		case true:
			branchDir = color.Sprintf(".git/modules/%s/refs/heads", submoduleName)
		}
	case "remote":
		switch isSubmodule {
		case false:
			branchDir = color.Sprintf(".git/refs/remotes/%s", remoteName)
		case true:
			branchDir = color.Sprintf(".git/modules/%s/refs/remotes/%s", submoduleName, remoteName)
		}
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
//   - 远程分支 refs/remotes/${remote}/<remoteBranchName>
//   - 本地分支 refs/heads/<localBranchName>
//
// 参数：
//   - repo: 本地仓库对象
//   - branchs: 远程分支信息
//
// 返回：
//   - 错误信息切片
func CreateLocalBranch(repo *git.Repository, branchs []fs.FileInfo) []string {
	// 使用一个 Slice 存储所有错误信息以美化输出
	var errList []string

	for _, branch := range branchs {
		// 修改 .git/config ，增加新的分支配置
		branchReferenceName := plumbing.NewBranchReferenceName(branch.Name()) // 构建本地分支 Reference 名，格式： refs/heads/<localBranchName>
		repo.CreateBranch(&config.Branch{                                     // 分支配置写入 .git/config
			Name:   branch.Name(),
			Remote: remoteName,
			Merge:  branchReferenceName,
		})

		// 创建一个新的 Reference
		newBranchReferenceName := plumbing.ReferenceName(branchReferenceName.String())    // refs/heads/test
		remoteReferenceName := plumbing.NewRemoteReferenceName(remoteName, branch.Name()) // 构建远程分支 Reference 名，格式： refs/remotes/${remote}/<remoteBranchName>
		remoteReferenceData, err := repo.Reference(remoteReferenceName, true)             // 根据远程分支 Reference 名获取其 Hash 值，格式： 1a8f900411d35a620407ce07902aecadfc782ded refs/remotes/${remote}/test
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

// ModifyGitConfig 修改 .git/config 文件，确保 [remote "origin"] 的 url 字段是以 'git@' 开头，并添加两行 pushurl
//
// 参数：
//   - configFile: .git/config 文件路径
//   - originalLink: 需要替换的原始链接
//   - newLink: 替换上去的新链接
//
// 返回：
//   - 错误信息
func ModifyGitConfig(configFile, originalLink, newLink string) error {
	// 以读写模式打开文件
	file, err := os.OpenFile(configFile, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取文件
	scanner := bufio.NewScanner(file) // 创建一个扫描器来读取文件内容
	var lines []string                // 存储读取到的行

	// 正则匹配（主仓库和子模块的匹配规则一样）
	regexPattern := `.*url\s*=\s*.*[:\/].*\.git` // 定义正则匹配规则
	regex := regexp.MustCompile(regexPattern)    // 创建正则表达式
	matched := false                             // 是否匹配到，用于限制只匹配一次

	// 需要新增的行
	pushUrl1 := "" // 第一行 pushurl
	pushUrl2 := "" // 第二行 pushurl

	// 逐行读取文件内容
	for scanner.Scan() {
		line := scanner.Text()

		// 检索一次模糊匹配的行
		if !matched && regex.MatchString(line) {
			// 第一次匹配：将可能存在的 "ssh://" 删除，并在"/"多于1个时将第1个替换为":"
			// 该次匹配是专对子模块的 .git/config 的处理
			line = strings.Replace(line, "ssh://", "", 1)
			if strings.Count(line, "/") >= 2 {
				line = strings.Replace(line, "/", ":", 1)
			}
			lines = append(lines, line)
			// 第二次匹配：创建2行 "pushurl"
			// 该次匹配是对于 .git/config 的通用处理
			pushUrl1 = strings.ReplaceAll(line, "url", "pushurl")
			pushUrl2 = strings.ReplaceAll(pushUrl1, originalLink, newLink)
			lines = append(lines, pushUrl1)
			lines = append(lines, pushUrl2)
			matched = true
		} else {
			lines = append(lines, line)
		}
	}

	// 将修改后的内容写回文件
	file.Truncate(0) // 清空文件内容
	file.Seek(0, 0)  // 移动光标到文件开头
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, _ = writer.WriteString(line + "\n")
	}
	writer.Flush()

	return nil
}

// GetDefaultBranchName 获取默认分支名
//
// 参数：
//   - repo: 本地仓库对象
//   - publicKeys: ssh 公钥
//
// 返回：
//   - 默认分支名
//   - 错误信息切片
func GetDefaultBranchName(repo *git.Repository, publicKeys *ssh.PublicKeys) (string, []string) {
	var defaultBranchName string
	// 使用一个 Slice 存储所有错误信息以美化输出
	var errList []string

	// 获取默认分支名
	remotes, _ := repo.Remotes() // 远程仓库信息
	for _, remote := range remotes {
		references, err := remote.List(&git.ListOptions{Auth: publicKeys}) // 远程引用信息
		if err != nil {
			errList = append(errList, "Failed to list references: "+err.Error())
			continue
		}
		for _, reference := range references {
			if reference.Name().Short() == "HEAD" { // 寻找 HEAD 引用
				// 输出 HEAD 分支名称
				defaultBranchName = reference.Target().Short()
				break
			}
		}
	}

	return defaultBranchName, errList
}

// CheckoutBranch 切换到指定分支
//
// 参数：
//   - worktree: 仓库的 git 工作树对象
//   - branchName: 分支名
//
// 返回：
//   - 错误信息
func CheckoutBranch(worktree *git.Worktree, branchName string) error {
	// 获取分支引用
	branch := plumbing.ReferenceName("refs/heads/" + branchName)

	// 切换分支
	err := worktree.Checkout(&git.CheckoutOptions{
		Branch: branch,
		Force:  false, // 如果有未提交的更改，不强制切换分支（否则会丢弃本地更改）
	})
	return err
}
