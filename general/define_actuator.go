/*
File: define_actuator.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 15:01:47

Description: 执行系统命令
*/

package general

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunCommandGetResult 运行命令并返回命令的输出
//
// 参数：
//   - command: 命令
//   - args: 命令参数
//
// 返回：
//   - 命令的输出
//   - 错误信息
func RunCommandGetResult(command string, args []string) (string, error) {
	if _, err := exec.LookPath(command); err != nil {
		return "", err
	}

	// 定义命令
	cmd := exec.Command(command, args...)
	// 执行命令并获取命令输出
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// 类型转换消除乱码和格式问题
	result := strings.TrimRight(string(output), "\n")
	return result, nil
}

// RunCommand 运行命令不返回命令的输出
//
// 参数：
//   - command: 命令
//   - args: 命令参数
//
// 返回：
//   - 错误信息
func RunCommand(command string, args []string) error {
	if _, err := exec.LookPath(command); err != nil {
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

// RunScript 运行 shell 脚本
//
// 参数：
//   - filePath: 脚本所在目录
//   - scriptName: 脚本名
//
// 返回：
//   - 错误信息
func RunScript(filePath, scriptName string) error {
	// 判断是否存在脚本文件，存在则运行脚本，不存在则忽略
	if FileExist(filepath.Join(filePath, scriptName)) {
		// 进到指定目录
		if err := os.Chdir(filePath); err != nil {
			return err
		}
		// 运行脚本
		bashArgs := []string{scriptName}
		if err := RunCommand("bash", bashArgs); err != nil {
			return err
		}
	}
	return nil
}
