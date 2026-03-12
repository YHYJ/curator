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
	"time"

	"github.com/gookit/color"
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

// 配置项
var (
	// 允许用户修改的配置项
	// 使用默认值的配置项
	rsaFile        = filepath.Join(UserInfo.HomeDir, ".ssh", "id_rsa")
	storagePath    = filepath.Join(UserInfo.HomeDir, "Documents", "Repos")
	scriptRunQueue = []string{"create-git-hook.sh"}
	githubUrl      = "github.com"
	githubUsername = "YHYJ"
	giteaUrl       = "git.yj1516.top"
	giteaUsername  = "YJ"
	repos          = []string{
		"checker",
		"curator",
		"eniac",
		"kbdstage",
		"manager",
		"msgcenter",
		"mqttc",
		"rolling",
		"scleaner",
		"skynet",
		"trash",
		"wocker",
		"www",
		"LearningCenter",
		"LogWrapper",
		"Modules",
		"MyBlogs",
		"MyDocker",
		"MyDockerfile",
		"MyRaspberry",
		"MyShell",
		"MyWiki",
		"Sniffer",
		"System",
		"Test",
		"YHYJ",
	}
)

// 配置
var appConfig = Config{
	Git: GitConfig{
		GithubUrl:      githubUrl,
		GithubUsername: githubUsername,
		GiteaUrl:       giteaUrl,
		GiteaUsername:  giteaUsername,
		Repos:          repos,
	},
	Script: ScriptConfig{
		RunQueue: scriptRunQueue,
	},
	SSH: SSHConfig{
		RsaFile: rsaFile,
	},
	Storage: StorageConfig{
		Path: storagePath,
	},
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
	// 打开配置文件
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 写入注释
	manual := color.Sprintf("##\n## %s - %s\n## Generaled on %s\n##\n\n", Name, Version, time.Now().Format("2006-01-02 15:04:05"))
	n, err := file.WriteString(manual)
	if err != nil {
		return int64(n), err
	}

	// 创建编码器并设置顺序保留
	encoder := toml.NewEncoder(file)
	encoder.Order(toml.OrderPreserve)

	if err := encoder.Encode(appConfig); err != nil {
		return int64(n), err
	}

	stat, _ := file.Stat()
	return stat.Size(), nil
}
