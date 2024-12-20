/*
File: define_actuator.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 15:01:47

Description: 执行系统命令
*/

package general

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

// RunCommandToOS 运行命令，将命令的 Stdin, Stdout 和 Stderr 定向到系统标准输入、标准输出和标准错误
//
// 参数：
//   - command: 命令
//   - args: 命令参数（每个以空格分隔的参数作为切片的一个元素）
//
// 返回：
//   - 错误信息
func RunCommandToOS(command string, args []string) error {
	// 检查命令是否存在，添加了对 `sudo` 命令的规避，`sudo`命令应独立检测
	if command == "sudo" {
		command = args[0]
		args = args[1:]
		if _, err := exec.LookPath(command); err != nil {
			return err
		}
	}

	// 定义命令
	cmd := exec.Command(command, args...)

	// 将命令的 Stdin, Stdout 和 Stderr 定向到系统标准输出和标准错误
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	err := cmd.Run()

	return err
}

// RunCommandToBuffer 运行命令，将命令的 Stdout 和 Stderr 定向到字节缓冲区
//
//   - 命令的 Stdout 和 Stderr 末尾自带的换行符已去除
//
// 参数：
//   - command: 命令
//   - args: 命令参数（每个以空格分隔的参数作为切片的一个元素）
//
// 返回：
//   - Stdout 缓冲区内容
//   - Stderr 缓冲区内容
//   - 错误信息
func RunCommandToBuffer(command string, args []string) (string, string, error) {
	// 检查命令是否存在，添加了对 `sudo` 命令的规避，`sudo`命令应独立检测
	if command == "sudo" {
		command = args[0]
		args = args[1:]
		if _, err := exec.LookPath(command); err != nil {
			return "", "", err
		}
	}

	// 定义命令
	cmd := exec.Command(command, args...)

	// 创建字节缓冲区
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// 将命令的 Stdout 和 Stderr 定向到字节缓冲区
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()

	// 去除缓冲区字符串末尾的换行符
	modifiedStdout := strings.TrimRightFunc(stdout.String(), unicode.IsSpace)
	modifiedStderr := strings.TrimRightFunc(stderr.String(), unicode.IsSpace)
	// stdXXX.Truncate(stdXXX.Len() - 1) 在 stdXXX 长度为0时需要多一步处理，否则 Out of range

	return modifiedStdout, modifiedStderr, err
}
