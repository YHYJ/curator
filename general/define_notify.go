/*
File: define_notify.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-06-17 14:32:57

Description: 消息通知
*/

package general

import "github.com/gookit/color"

var Notifier []string // 通知器

// Notification 显示通知
func Notification() {
	if len(Notifier) > 0 {
		for _, slogan := range Notifier {
			color.Notice.Tips(PrimaryText(slogan))
		}
	}
}
