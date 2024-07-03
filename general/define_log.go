/*
File: define_log.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-06-13 15:04:44

Description: 处理日志
*/

package general

import (
	"path/filepath"
	"runtime"
	"strings"
)

// GetCallerInfo 获取调用者信息
//
// 返回：
//   - 调用者所在文件名（不带后缀）
//   - 调用者所在行号
func GetCallerInfo() (string, int) {
	// runtime.Caller 的参数 skip 是要上升的堆栈数，0 表示 Caller 的调用者，1 表示上层的调用者
	_, fullFilePath, line, ok := runtime.Caller(1)
	if !ok {
		return "", 0
	}
	file := strings.Split(filepath.Base(fullFilePath), ".")[0]
	return file, line
}
