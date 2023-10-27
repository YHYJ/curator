/*
File: version.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 13:22:05

Description: 子命令`version`功能函数
*/

package function

import "fmt"

// 程序信息
const (
	Name    = "Clone-repos"
	Version = "v0.6.2"
	Path    = "github.com/yhyj/Clone-repos"
)

func ProgramInfo(only bool) string {
	programInfo := fmt.Sprintf("%s\n", Version)
	if !only {
		programInfo = fmt.Sprintf("%s version %s\n", Name, Version)
	}
	return programInfo
}
