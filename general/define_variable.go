/*
File: define_variable.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 13:46:00

Description: 操作变量（包括代码变量和环境变量）
*/

package general

import (
	"os"
	"os/user"
	"runtime"
	"strconv"
)

// ---------- 代码变量

// 输出格式变量
var (
	// 常规输出格式 1部分: <输出内容>
	Regelar1PFormat = "%s\n"
	// 常规输出格式 2部分: <输出内容1><输出内容2>
	Regelar2PFormat = "%s%s\n"

	// 标题部分输出格式 H1级别标题: <标题文字>
	TitleH1Format = "\n\x1b[36;3m%s\x1b[0m\n\n"

	// 分隔线输出格式 隐性分隔线: <分隔线>
	LineHiddenFormat = "\x1b[30m%s\x1b[0m\n"
	// 分隔线输出格式 显性分隔线: <分隔线>
	LineShownFormat = "\x1b[30;1m%s\x1b[0m\n"

	// Slice输出格式 切片遍历·2部分: <切片内容1><分隔符><切片内容2>
	SliceTraverse2PFormat = "\x1b[32;1m%s\x1b[0m%s\x1b[34m%s\x1b[0m\n"
	// Slice输出格式 带后缀的切片遍历·2部分: <切片内容1><分隔符><切片内容2><分隔符><后缀>
	SliceTraverse2PSuffixFormat = "\x1b[32;1m%s\x1b[0m%s\x1b[34m%s\x1b[0m%s%s\n"
	// Slice输出格式 带后缀的切片遍历·2部分·不换行: <切片内容1><分隔符><切片内容2><分隔符><后缀>
	SliceTraverse2PSuffixNoNewLineFormat = "\x1b[32;1m%s\x1b[0m%s\x1b[34m%s\x1b[0m%s%s"
	// Slice输出格式 带后缀的切片遍历·3部分: <切片内容1><分隔符><切片内容2><分隔符><切片内容3><分隔符><后缀>
	SliceTraverse3PSuffixFormat = "\x1b[32;1m%s\x1b[0m%s\x1b[34m%s\x1b[0m%s\x1b[33;1m%s\x1b[0m%s%s\n"
	// Slice输出格式 带后缀的切片遍历·4部分: <切片内容1><分隔符><切片内容2><分隔符><切片内容3><分隔符><切片内容4><分隔符><后缀>
	SliceTraverse4PSuffixFormat = "\x1b[32;1m%s\x1b[0m%s\x1b[34m%s\x1b[0m%s\x1b[33m%s\x1b[0m%s\x1b[35;1m%s\x1b[0m%s%s\n"

	// 成功信息输出格式 成功信息: <成功信息>
	SuccessFormat = "\x1b[32;1m%s\x1b[0m\n"
	// 成功信息输出格式 成功信息·不换行: <成功信息>
	SuccessNoNewLineFormat = "\x1b[32;1m%s\x1b[0m"
	// 成功信息输出格式 带后缀的成功信息: <成功信息><分隔符><后缀>
	SuccessSuffixFormat = "\x1b[32;1m%s\x1b[0m%s%s\n"
	// 成功信息输出格式 带后缀的成功信息·不换行: <成功信息><分隔符><后缀>
	SuccessSuffixNoNewLineFormat = "\x1b[32;1m%s\x1b[0m%s%s"

	// 提示信息输出格式 带前缀的提示信息: <提示信息>
	TipsPrefixFormat = "%s%s\x1b[32;1m%s\x1b[0m\n"
	// 提示信息输出格式 带后缀的提示信息·2部分·不换行: <提示信息1><分隔符><提示信息2><分隔符><后缀>
	Tips2PSuffixNoNewLineFormat = "\x1b[32;1m%s\x1b[0m%s\x1b[36;1m%s\x1b[0m%s%s"

	// 展示信息输出格式 展示信息: <展示信息>
	InfoFormat = "\x1b[33;1m%s\x1b[0m\n"
	// 展示信息输出格式 带前缀的展示信息: <前缀><分隔符><展示信息>
	InfoPrefixFormat = "%s%s\x1b[33;1m%s\x1b[0m\n"
	// 展示信息输出格式 带前缀的展示信息·2部分: <前缀><分隔符><展示信息1><分隔符><展示信息2>
	Info2PPrefixFormat = "%s%s\x1b[33;1m%s\x1b[0m%s\x1b[35m%s\x1b[0m\n"
	// 展示信息输出格式 带后缀的展示信息: <展示信息><分隔符><后缀>
	InfoSuffixFormat = "\x1b[33;1m%s\x1b[0m%s%s\n"
	// 展示信息输出格式 带前后缀的展示信息: <前缀><分隔符><展示信息><分隔符><后缀>
	InfoPrefixSuffixFormat = "%s%s\x1b[33;1m%s\x1b[0m%s%s\n"

	// 错误信息输出格式 基础错误: <错误信息>
	ErrorBaseFormat = "\x1b[31m%s\x1b[0m\n"
	// 错误信息输出格式 带后缀的错误: <错误信息><分隔符><后缀>
	ErrorSuffixFormat = "\x1b[31m%s\x1b[0m%s%s\n"
)

// ---------- 环境变量

// 操作系统
var Platform = runtime.GOOS

// 用户名，当程序提权运行时，使用SUDO_USER变量获取提权前的用户名
var UserName = func() string {
	if GetVariable("SUDO_USER") != "" {
		return GetVariable("SUDO_USER")
	}
	return GetVariable("USER")
}()

// 用户信息
var UserInfo, _ = GetUserInfoByName(UserName)

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

// GetVariable 获取环境变量
func GetVariable(key string) string {
	if innerMap, exists := platformChart[Platform]; exists {
		if _, variableExists := innerMap[key]; variableExists {
			key = platformChart[Platform][key]
		}
	}
	variable := os.Getenv(key)

	return variable
}

// GetHostname 获取系统 HOSTNAME
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

// SetVariable 设置环境变量
func SetVariable(key, value string) error {
	return os.Setenv(key, value)
}

// GetUserInfoByName 根据用户名获取用户信息
func GetUserInfoByName(userName string) (*user.User, error) {
	userInfo, err := user.Lookup(userName)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GetUserInfoById 根据 ID 获取用户信息
func GetUserInfoById(userId int) (*user.User, error) {
	userInfo, err := user.LookupId(strconv.Itoa(userId))
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GetCurrentUserInfo 获取当前用户信息
func GetCurrentUserInfo() (*user.User, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}
