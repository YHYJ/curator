/*
File: root.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 15:16:00

Description:
*/

package function

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	gssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pelletier/go-toml"
	"golang.org/x/crypto/ssh"
)

func getTomlConfig(filename string) (*toml.Tree, error) {
	tree, err := toml.LoadFile(filename)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func getSshKeyAuth(privateSshKeyFile string) transport.AuthMethod {
	var auth transport.AuthMethod
	sshKey, _ := ioutil.ReadFile(privateSshKeyFile)
	signer, _ := ssh.ParsePrivateKey([]byte(sshKey))
	auth = &gssh.PublicKeys{User: "git", Signer: signer}
	return auth
}

func RollingCLoneRepos(confile string) {
	// 加载配置文件
	conf, err := getTomlConfig(confile)
	if err != nil {
		fmt.Printf("\x1b[36;1m%s\x1b[0m\n", err)
	} else {
		// 获取配置项
		private_key_file := conf.Get("ssh.private_key_file")
		path := conf.Get("storage.path").(string)
		url := conf.Get("git.url").(string)
		repos := conf.Get("git.repos").([]interface{})
		auth := getSshKeyAuth(private_key_file.(string))
		// 开始克隆
		for _, repo := range repos {
			_, err := git.PlainClone(path+"/"+repo.(string), false, &git.CloneOptions{
				URL:               url + "/" + repo.(string) + ".git",
				Auth:              auth,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
				Progress:          os.Stdout,
			})
			if err != nil {
				fmt.Printf("\x1b[36;1m%s %s\x1b[0m\n", "Clone "+repo.(string)+" faild: ", err)
			} else {
				fmt.Printf("\x1b[36;1m%s\x1b[0m\n", "Clone "+repo.(string)+" success")
			}
		}
	}
}
