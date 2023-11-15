/*
File: define_actuator.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 15:01:47

Description: 执行系统命令的函数
*/

package general

import (
	"os"
	"os/exec"
	"strings"
)

// RunCommandGetResult 运行指定命令并获取命令输出
func RunCommandGetResult(command string, args []string) (string, error) {
	_, err := exec.LookPath(command)
	if err != nil {
		return "", err
	}

	// 定义命令
	cmd := exec.Command(command, args...)
	// 执行命令并获取命令输出
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// 类型转换
	result := strings.TrimRight(string(output), "\n")
	return result, nil
}

// RunCommand 运行指定命令（命令无输出）
func RunCommand(command string, args []string) error {
	_, err := exec.LookPath(command)
	if err != nil {
		return err
	}

	// 定义命令
	cmd := exec.Command(command, args...)
	// 定义标准输入/输出/错误
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 执行命令
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
