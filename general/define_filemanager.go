/*
File: define_filemanager.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-24 16:41:17

Description: 文件管理
*/

package general

import (
	"io"
	"os"
	"path/filepath"
)

// FileExist 判断文件是否存在
//
// 参数：
//   - filePath: 文件路径
//
// 返回：
//   - 文件存在返回 true，否则返回 false
func FileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// FolderEmpty 判断文件夹是否为空
//
//   - 包括隐藏文件
//
// 参数：
//   - dir: 文件夹路径
//
// 返回：
//   - 文件夹为空返回 true，否则返回 false
func FolderEmpty(dir string) bool {
	text, err := os.Open(dir)
	if err != nil {
		return true
	}
	defer text.Close()

	if _, err = text.Readdir(1); err == io.EOF {
		return true
	}
	return false
}

// CreateFile 创建文件，包括其父目录
//
// 参数：
//   - file: 文件路径
//
// 返回：
//   - 错误信息
func CreateFile(file string) error {
	if FileExist(file) {
		return nil
	}
	// 创建父目录
	parentPath := filepath.Dir(file)
	if err := os.MkdirAll(parentPath, os.ModePerm); err != nil {
		return err
	}
	// 创建文件
	if _, err := os.Create(file); err != nil {
		return err
	}

	return nil
}

// DeleteFile 删除文件，如果目标是文件夹则包括其下所有文件
//
// 参数：
//   - filePath: 文件路径
//
// 返回：
//   - 错误信息
func DeleteFile(filePath string) error {
	if !FileExist(filePath) {
		return nil
	}
	return os.RemoveAll(filePath)
}
