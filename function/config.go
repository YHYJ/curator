/*
File: config.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-24 16:41:33

Description: 子命令`config`的实现
*/

package function

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
)

// 判断文件是不是toml文件
func isTomlFile(filePath string) bool {
	if strings.HasSuffix(filePath, ".toml") {
		return true
	}
	return false
}

// 读取toml配置文件
func GetTomlConfig(filePath string) (*toml.Tree, error) {
	if !FileExist(filePath) {
		return nil, fmt.Errorf("open %s: no such file or directory", filePath)
	}
	if !isTomlFile(filePath) {
		return nil, fmt.Errorf("open %s: is not a toml file", filePath)
	}
	tree, err := toml.LoadFile(filePath)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// 写入toml配置文件
func WriteTomlConfig(filePath string) (int64, error) {
	// 获取指定用户信息
	userInfo, err := GetUserInfo(1000)
	if err != nil {
		return 0, err
	}
	// 定义一个map[string]interface{}类型的变量并赋值
	exampleConf := map[string]interface{}{
		"ssh": map[string]interface{}{
			"private_key_file": userInfo.HomeDir + "/.ssh/id_rsa",
		},
		"storage": map[string]interface{}{
			"path": userInfo.HomeDir + "/Documents/Repos",
		},
		"git": map[string]interface{}{
			"url": "git@github.com:YHYJ",
			"repos": []string{
				"checker",
				"clone-repos",
				"DataCollector",
				"DataWizard",
				"Devicer",
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
				"OPC2Socket",
				"rolling",
				"scleaner",
				"skynet",
				"Sniffer",
				"System",
				"Test",
				"TimeNote",
				"WeatherStation",
				"YHYJ",
			},
		},
	}
	if !FileExist(filePath) {
		return 0, fmt.Errorf("open %s: no such file or directory", filePath)
	}
	if !isTomlFile(filePath) {
		return 0, fmt.Errorf("open %s: is not a toml file", filePath)
	}
	// 把exampleConf转换为*toml.Tree
	tree, err := toml.TreeFromMap(exampleConf)
	if err != nil {
		return 0, err
	}
	// 打开一个文件并获取io.Writer接口
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	return tree.WriteTo(file)
}
