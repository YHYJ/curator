/*
File: define_message.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-05-29 15:51:14

Description: 定义输出信息及其格式
*/

package general

var (
	MultiSelectTips  = "Please select from the %s below (multi-select)\n"  // 提示词 - 多选
	SingleSelectTips = "Please select from the %s below (single-select)\n" // 提示词 - 单选
	QuietTips        = "Press '%s' to quit\n"                              // 提示词 - 退出
	SelectAllTips    = "Select All"                                        // 提示词 - 全选
)

var (
	OverWriteTips = "%s file already exists, do you want to overwrite it?" // 提示词 - 文件已存在是否覆写
)
