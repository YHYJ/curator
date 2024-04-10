/*
File: config.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-24 16:41:33

Description: 子命令 'config' 的实现
*/

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/yhyj/curator/general"
)

// 用于转换 Toml 配置树的结构体
type Config struct {
	Git     GitConfig     `toml:"git"`
	Script  ScriptConfig  `toml:"script"`
	SSH     SSHConfig     `toml:"ssh"`
	Storage StorageConfig `toml:"storage"`
}
type GitConfig struct {
	GithubUrl      string   `toml:"github_url"`
	GithubUsername string   `toml:"github_username"`
	GiteaUrl       string   `toml:"gitea_url"`
	GiteaUsername  string   `toml:"gitea_username"`
	Repos          []string `toml:"repos"`
}
type ScriptConfig struct {
	NameList []string `toml:"name_list"`
}
type SSHConfig struct {
	RsaFile string `toml:"rsa_file"`
}
type StorageConfig struct {
	Path string `toml:"path"`
}

// isTomlFile 检测文件是不是 toml 文件
//
// 参数：
//   - filePath: 待检测文件路径
//
// 返回：
//   - 是 toml 文件返回 true，否则返回 false
func isTomlFile(filePath string) bool {
	if strings.HasSuffix(filePath, ".toml") {
		return true
	}
	return false
}

// GetTomlConfig 读取 toml 配置文件
//
// 参数：
//   - filePath: toml 配置文件路径
//
// 返回：
//   - toml 配置树
//   - 错误信息
func GetTomlConfig(filePath string) (*toml.Tree, error) {
	if !general.FileExist(filePath) {
		return nil, fmt.Errorf("Open %s: no such file or directory", filePath)
	}
	if !isTomlFile(filePath) {
		return nil, fmt.Errorf("Open %s: is not a toml file", filePath)
	}
	tree, err := toml.LoadFile(filePath)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// LoadConfigToStruct 将 Toml 配置树加载到结构体
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
//
// 返回：
//   - 结构体
//   - 错误信息
func LoadConfigToStruct(configTree *toml.Tree) (*Config, error) {
	var config Config
	if err := configTree.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// WriteTomlConfig 写入 toml 配置文件
//
// 参数：
//   - filePath: toml 配置文件路径
//
// 返回：
//   - 写入的字节数
//   - 错误信息
func WriteTomlConfig(filePath string) (int64, error) {
	// 根据系统不同决定某些参数
	var (
		scriptNameList = []string{} // 脚本名列表
	)
	if general.Platform == "linux" {
		scriptNameList = []string{
			"create-hook-link.sh",
		}
	} else if general.Platform == "darwin" {
		scriptNameList = []string{
			"create-hook-link.sh",
		}
	} else if general.Platform == "windows" {
	}
	// 定义一个 map[string]interface{} 类型的变量并赋值
	exampleConf := map[string]interface{}{
		"ssh": map[string]interface{}{
			"rsa_file": filepath.Join(general.UserInfo.HomeDir, ".ssh", "id_rsa"),
		},
		"storage": map[string]interface{}{
			"path": filepath.Join(general.UserInfo.HomeDir, "Documents", "Repos"),
		},
		"script": map[string]interface{}{
			"name_list": scriptNameList,
		},
		"git": map[string]interface{}{
			"github_url":      "github.com",
			"github_username": "YHYJ",
			"gitea_url":       "git.yj1516.top",
			"gitea_username":  "YJ",
			"repos": []string{
				"checker",
				"curator",
				"eniac",
				"kbdstage",
				"LearningCenter",
				"LogWrapper",
				"manager",
				"Modules",
				"MyBlogs",
				"MyDocker",
				"MyDockerfile",
				"MyRaspberry",
				"MyShell",
				"MyWiki",
				"rolling",
				"scleaner",
				"skynet",
				"Sniffer",
				"System",
				"Test",
				"trash",
				"www",
				"YHYJ",
			},
		},
	}
	// 检测配置文件是否存在
	if !general.FileExist(filePath) {
		return 0, fmt.Errorf("Open %s: no such file or directory", filePath)
	}
	// 检测配置文件是否是 toml 文件
	if !isTomlFile(filePath) {
		return 0, fmt.Errorf("Open %s: is not a toml file", filePath)
	}
	// 把 exampleConf 转换为 *toml.Tree 类型
	tree, err := toml.TreeFromMap(exampleConf)
	if err != nil {
		return 0, err
	}
	// 打开一个文件并获取 io.Writer 接口
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	return tree.WriteTo(file)
}
