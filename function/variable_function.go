/*
File: variable_function.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 13:46:00

Description: 执行变量操作的函数
*/

package function

import "os"

var platformChart = map[string]map[string]string{
	"linux": {
		"HOME": "HOME",
		"PWD":  "PWD",
		"USER": "USER",
	},
	"darwin": {
		"HOME": "HOME",
		"PWD":  "PWD",
		"USER": "USER",
	},
	"windows": {
		"HOME": "USERPROFILE",
		"PWD":  "PWD",
		"USER": "USERNAME",
	},
}

var platform = os.Getenv("GOOS")

func getVariable(key string) string {
	return platformChart[platform][key]
}
