/*
File: define_convert.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-20 12:45:37

Description: 单位和格式数据转换
*/

package general

import (
	"strings"
)

// UpperFirstChar 最大化字符串的第一个字母
//
// 参数：
//   - str: 需要处理的字符串
//
// 返回：
//   - 处理后的字符串
func UpperFirstChar(str string) string {
	if len(str) == 0 {
		return str
	}

	return strings.ToUpper(str[:1]) + str[1:]
}
