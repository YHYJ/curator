/*
File: variable_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 13:46:00

Description: 操作变量
*/

package function

import (
	"os"
	"os/user"
	"runtime"
	"strconv"
)

// 用来处理不同系统之间的变量名差异
var platformChart = map[string]map[string]string{
	"windows": {
		"HOME":     "USERPROFILE",  // 用户主目录路径
		"USER":     "USERNAME",     // 当前登录用户名
		"SHELL":    "ComSpec",      // 默认shell或命令提示符路径
		"PWD":      "CD",           // 当前工作目录路径
		"HOSTNAME": "COMPUTERNAME", // 计算机主机名
	},
}

var platform = runtime.GOOS

// 获取环境变量
func GetVariable(key string) string {
	if innerMap, exists := platformChart[platform]; exists {
		if _, variableExists := innerMap[key]; variableExists {
			key = platformChart[platform][key]
		}
	}
	variable := os.Getenv(key)

	return variable
}

// 获取不在环境变量中的HOSTNAME
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

// 设置环境变量
func SetVariable(key, value string) error {
	return os.Setenv(key, value)
}

// 根据用户名获取用户信息
func GetUserInfoByName(userName string) (*user.User, error) {
	userInfo, err := user.Lookup(userName)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// 根据ID获取用户信息
func GetUserInfoById(userId int) (*user.User, error) {
	userInfo, err := user.LookupId(strconv.Itoa(userId))
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// 获取当前用户信息
func GetCurrentUserInfo() (*user.User, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}
