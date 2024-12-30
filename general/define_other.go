/*
File: define_other.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-11-24 13:35:18

Description: 处理一些杂事
*/

package general

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/gookit/color"
)

// Delay 延时
//
// 参数：
//   - second: 延时秒数
func Delay(second float32) {
	time.Sleep(time.Duration(second*1000) * time.Millisecond)
}

// AreYouSure 获取用户二次确认
//
// 参数：
//   - question: 问题
//   - defaultAnswer: 默认回答，true 或 false
//
// 返回：
//   - true/false
//   - 错误信息
func AreYouSure(question string, defaultAnswer bool) (bool, error) {
	var (
		viewAnswers []string                                 // 显示用可选答案
		answersMap  = map[string]bool{"y": true, "n": false} // 可选答案和实际返回值的映射
		reader      = bufio.NewReader(os.Stdin)              // 标准输入
	)

	// 根据 defaultAnswer 设置显示用的可选答案
	if defaultAnswer == true {
		viewAnswers = []string{"Y", "n"}
	} else {
		viewAnswers = []string{"y", "N"}
	}
	viewAnswersConsortium := strings.Join(viewAnswers, "/") // 显示用可选答案的组合体

	for {
		// 输出问题
		color.Printf("%s %s: ", question, SecondaryText("(", viewAnswersConsortium, ")"))

		// 从标准输入中读取用户的回答
		userRawAnswer, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		// 去除用户回答中的换行符
		userAnswer := strings.TrimSpace(strings.TrimSuffix(userRawAnswer, "\n"))

		// 检查用户回答是否符合要求
		for answer, result := range answersMap {
			if strings.EqualFold(userAnswer, answer) {
				return result, nil
			} else if userAnswer == "" { // 如果用户回答为空，返回默认回答
				return defaultAnswer, nil
			}
		}
	}
}
