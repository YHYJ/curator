/*
File: define_flag.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-05-29 15:48:58

Description: 定义符号
*/

package general

var (
	RunFlag     = "🐙"  // 运行状态符号 - 运行中
	LatestFlag  = "🌟"  // 运行状态符号 - 已是最新
	SuccessFlag = "✅"  // 运行状态符号 - 成功
	WarningFlag = "⚠️" // 运行状态符号 - 警告
	ErrorFlag   = "❌"  // 运行状态符号 - 失败
)

var (
	CursorOnFlag   = "👉"  // 选择器符号 - 光标在
	CursorOffFlag  = "  " // 选择器符号 - 光标不在
	SelectedFlag   = "•"  // 选择器符号 - 已选中
	UnselectedFlag = " "  // 选择器符号 - 已选中
	SelectAllFlag  = "⭐️" // 选择器符号 - 全选
)

var (
	Separator1st = "=" // 分隔符 - 1级
	Separator2st = "-" // 分隔符 - 2级
	Separator3st = "·" // 分隔符 - 3级
)

var (
	BranchFlag    = "🌿" // Git 符号 - 分支
	SubmoduleFlag = "📦" // Git 符号 - 子模块
)

var (
	JoinerIng    = "├──" // 条目连接符号 - 中间条目
	JoinerFinish = "└──" // 条目连接符号 - 最后条目
)
