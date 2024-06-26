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

// AskUser 询问用户
//
// 参数：
//   - question: 问题
//   - standardAnswers: 期望回答的切片（最后一个选项是默认值），例如 [y, N] 代表期望输入 y 或 n，最后一个选项是默认值（大写为了提示用户其为默认值）
//
// 返回：
//   - 回答
//   - 错误信息
func AskUser(question string, standardAnswers []string) (string, error) {
	viewAnswers := strings.Join(standardAnswers, "/")
	color.Printf("%s %s: ", question, SecondaryText("(", viewAnswers, ")"))

	// 默认回答
	var answer = strings.ToLower(standardAnswers[len(standardAnswers)-1])

	// 从标准输入中读取用户的回答
	reader := bufio.NewReader(os.Stdin)
	userRawAnswer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// 用户回答不为空则符合要求，转为小写
	userAnswer := strings.TrimSpace(strings.TrimSuffix(userRawAnswer, "\n"))
	if len(userAnswer) != 0 {
		answer = strings.ToLower(userAnswer)
	}

	// 检测输入是否符合要求，不符合则返回默认值
	for _, standardAnswer := range standardAnswers {
		if answer == standardAnswer {
			return answer, nil
		}
	}

	return answer, nil
}
