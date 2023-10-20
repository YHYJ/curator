/*
File: define_variable.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-20 14:11:25

Description: 定义项目变量
*/

package function

// 用户名，当程序提权运行时，使用SUDO_USER变量获取提权前的用户名
var UserName = func() string {
	if GetVariable("SUDO_USER") != "" {
		return GetVariable("SUDO_USER")
	}
	return GetVariable("USER")
}()

// 用户信息
var UserInfo, _ = GetUserInfoByName(UserName)
