/*
File: run_command.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 15:01:47

Description: 执行系统命令的函数
*/

package function

import (
	"os"
	"os/exec"
	"strings"
)

// 运行指定命令并获取命令输出
func RunCommandGetResult(command string, args []string) (result string, err error) {
	_, err = exec.LookPath(command)
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
	result = strings.TrimRight(string(output), "\n")
	return result, nil
}

// 运行指定命令（命令无输出）
func RunCommand(command string, args []string) (err error) {
	_, err = exec.LookPath(command)
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
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
