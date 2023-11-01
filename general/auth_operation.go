/*
File: auth_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-11 14:30:32

Description: 身份认证
*/

package general

import (
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	cssh "golang.org/x/crypto/ssh"
)

// 使用go-git自带的方法获取公钥
func GetPublicKeysByGit(pemFile, password string) (*ssh.PublicKeys, error) {
	_, err := os.Stat(pemFile)
	if err != nil {
		return nil, err
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", pemFile, password)
	if err != nil {
		return nil, err
	}
	return publicKeys, nil
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
