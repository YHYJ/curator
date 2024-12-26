/*
File: define_toml.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-11 14:16:09

Description: 操作 TOML 配置文件
*/

package general

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
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
	RunQueue []string `toml:"run_queue"`
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
	if !FileExist(filePath) {
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
	// 存储库初次克隆到本地后自动执行的脚本队列
	var scriptRunQueue = []string{"create-git-hook.sh"}

	// 定义一个 map[string]interface{} 类型的变量并赋值
	exampleConf := map[string]interface{}{
		"ssh": map[string]interface{}{
			"rsa_file": filepath.Join(UserInfo.HomeDir, ".ssh", "id_rsa"),
		},
		"storage": map[string]interface{}{
			"path": filepath.Join(UserInfo.HomeDir, "Documents", "Repos"),
		},
		"script": map[string]interface{}{
			"run_queue": scriptRunQueue,
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
				"msgcenter",
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
				"wocker",
				"www",
				"YHYJ",
			},
		},
	}
	// 检测配置文件是否存在
	if !FileExist(filePath) {
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
