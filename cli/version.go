/*
File: version.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-11 13:30:05

Description: 子命令 'version' 的实现
*/

package cli

import (
	"github.com/gookit/color"
	"github.com/yhyj/curator/general"
)

// PrintVersionInfo 打印版本信息
//
// 参数：
//   - only: 是否只打印版本号
func PrintVersionInfo(only bool) {
	// 获取版本信息
	programInfo := general.ProgramInfo()
	// 根据参数决定输出
	if only {
		color.Printf("%s\n", programInfo["Version"])
	} else {
		color.Printf("%s %s\n", general.LightText(programInfo["Name"]), general.LightText(programInfo["Version"]))
		color.Printf("%s %s\n", general.SecondaryText("Project:"), general.SecondaryText(programInfo["Project"]))
		color.Printf("%s %s\n", general.SecondaryText("Build rev:"), general.SecondaryText(programInfo["GitCommitHash"]))
		color.Printf("%s %s\n", general.SecondaryText("Built on:"), general.SecondaryText(programInfo["BuildTime"]))
		color.Printf("%s %s\n", general.SecondaryText("Built by:"), general.SecondaryText(programInfo["BuildBy"]))
	}
}
