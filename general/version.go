/*
File: version.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 13:22:05

Description: 版本信息
*/

package general

import (
	"strconv"
	"time"
)

const (
	Name    string = "Curator"                 // 程序名
	Version string = "v1.1.3"                  // 程序版本
	Project string = "github.com/yhyj/curator" // 项目地址
)

var (
	GitCommitHash string = "Unknown" // Git 提交 Hash
	BuildTime     string = "Unknown" // 编译时间
	BuildBy       string = "Unknown" // 编译者
)

// ProgramInfo 返回程序信息
//
// 返回：
//   - 程序信息
func ProgramInfo() map[string]string {
	programInfo := make(map[string]string)

	// 解析构建时间
	BuildTimeConverted := "Unknown"
	timestamp, err := strconv.ParseInt(BuildTime, 10, 64)
	if err == nil {
		BuildTimeConverted = time.Unix(timestamp, 0).String()
	}

	programInfo["Name"] = Name
	programInfo["Version"] = Version
	programInfo["Project"] = Project
	programInfo["GitCommitHash"] = GitCommitHash
	programInfo["BuildTime"] = BuildTimeConverted
	programInfo["BuildBy"] = BuildBy

	return programInfo
}
