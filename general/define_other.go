/*
File: define_other.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-11-24 13:35:18

Description: 处理一些杂事
*/

package general

import "time"

// Delay 延时
//
// 参数：
//   - second: 延时秒数
func Delay(second float32) {
	time.Sleep(time.Duration(second*1000) * time.Millisecond)
}
