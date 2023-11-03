/*
File: define_authentication.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-11 14:30:32

Description: 身份认证函数
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
	cssh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// 使用go-git自带的方法获取公钥
func GetPublicKeysByGit(pemFile string) (*ssh.PublicKeys, error) {
	_, err := os.Stat(pemFile)
	if err != nil {
		return nil, err
	}

	// 初次尝试，默认密码为空
	publicKeys, err := ssh.NewPublicKeysFromFile("git", pemFile, "")
	if err != nil {
		if err == x509.IncorrectPasswordError || strings.Contains(err.Error(), "empty password") {
			maxAttempts := 3 // 最大尝试次数
			for attempts := 0; attempts < maxAttempts; attempts++ {
				fmt.Printf("Enter passphrase for key '%s' (%d/%d): ", pemFile, attempts+1, maxAttempts)
				password, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return nil, err
				}
				fmt.Println() // 换行
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

// 使用crypto/ssh模块获取公钥
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

// 清除内存中的密码，以增加安全性
func clearPassword(password []byte) {
	for i := range password {
		password[i] = 0
	}
}
