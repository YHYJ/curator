/*
File: define_authentication.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-11 14:30:32

Description: 身份认证
*/

package general

import (
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/gookit/color"
	cssh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// GetPublicKeysByGit 使用 go-git 自带的方法获取 ssh 公钥
//
// 参数：
//   - pemFile: 私钥文件路径
//
// 返回：
//   - ssh 公钥
//   - 错误信息
func GetPublicKeysByGit(pemFile string) (*ssh.PublicKeys, error) {
	if _, err := os.Stat(pemFile); err != nil {
		return nil, err
	}

	// 初次尝试，默认密码为空
	publicKeys, err := ssh.NewPublicKeysFromFile("git", pemFile, "")
	if err != nil {
		if err == x509.IncorrectPasswordError || strings.Contains(err.Error(), "empty password") {
			maxAttempts := 3 // 最大尝试次数
			for attempts := 0; attempts < maxAttempts; attempts++ {
				color.Printf("Enter passphrase for key '%s' (%s/%s): ", PrimaryText(pemFile), WarnText(attempts+1), NoticeText(maxAttempts))
				password, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return nil, err
				}
				color.Println() // 换行
				publicKeys, err := ssh.NewPublicKeysFromFile("git", pemFile, string(password))

				clearPassword(password)

				if err == nil {
					return publicKeys, nil
				}
			}
			return nil, fmt.Errorf("Permission denied (publickey)")
		} else {
			return nil, err
		}
	}
	return publicKeys, err
}

// GetPublicKeysBySSH 使用 crypto/ssh 模块获取 ssh 公钥
//
// 参数：
//   - pemFile: 私钥文件路径
//
// 返回：
//   - ssh 公钥
//   - 错误信息
func GetPublicKeysBySSH(pemFile string) (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	sshKey, err := os.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}
	signer, err := cssh.ParsePrivateKey([]byte(sshKey))
	if err != nil {
		return nil, err
	}
	auth = &ssh.PublicKeys{User: "git", Signer: signer}

	return auth, nil
}

// clearPassword 清除内存中的密码，以增加安全性
//
// 参数：
//   - password: 密码
func clearPassword(password []byte) {
	for i := range password {
		password[i] = 0
	}
}
