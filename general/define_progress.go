/*
File: define_progress.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-05-29 15:53:32

Description: 定义进度相关
*/

package general

import (
	"time"

	"github.com/briandowns/spinner"
)

var WaitSpinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond) // 等待动画
