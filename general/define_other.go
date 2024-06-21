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
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gookit/color"
)

// RealLength 去除控制字符和图标的附加字符，获取文本实际长度
//
// 参数：
//   - text: 文本
//
// 返回：
//   - 实际长度
func RealLength(text string) int {
	combinedRegex := regexp.MustCompile(`\x1b\[[0-9;]*m|\x{FE0E}|\x{FE0F}`)
	return utf8.RuneCountInString(combinedRegex.ReplaceAllString(text, ""))
}

// PrintDelimiter 打印分隔符
//
// 参数：
//   - length: 分隔符长度
func PrintDelimiter(length int) {
	dashes := strings.Repeat("-", length)     // 组装分隔符
	color.Printf("%s\n", FgBlackText(dashes)) // 美化输出
}

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

// GetInput 获取用户输入
//
// 参数：
//   - tips: 提示信息
//   - default: 用户未输入时的默认值
//
// 返回：
//   - 用户输入（去掉了最后的换行符）
//   - 错误信息
func GetInput(tips string, defaultValue string) (string, error) {
	color.Printf("%s %s: ", tips, SecondaryText("(", defaultValue, ")"))

	// 从标准输入中读取用户的回答
	reader := bufio.NewReader(os.Stdin)
	originalValue, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	value := func() string {
		if len(originalValue) <= 1 {
			return defaultValue
		}
		return strings.TrimSuffix(originalValue, "\n")
	}()

	return value, nil
}

// Capitalize 将字符串的首字母转换为大写
//
// 参数：
//   - s: 字符串
//
// 返回：
//   - 首字母大写后的字符串
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	// 将首字母转换为大写
	firstRune := []rune(s)[0]
	capitalizedFirstRune := unicode.ToUpper(firstRune)

	// 拼接首字母和剩余部分
	return string(capitalizedFirstRune) + s[1:]
}
