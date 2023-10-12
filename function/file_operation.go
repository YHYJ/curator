/*
File: file_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-24 16:41:17

Description: 操作文件
*/

package function

import (
	"io"
	"os"
	"strings"
)

// 判断文件是否存在
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断文件夹是否为空，包括隐藏文件
func FolderEmpty(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return true
	}
	defer file.Close()

	_, err = file.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}

// 创建文件，如果其父目录不存在则创建父目录
func CreateFile(filePath string) error {
	if FileExist(filePath) {
		return nil
	}
	// 截取filePath的父目录
	parentPath := filePath[:strings.LastIndex(filePath, "/")]
	if err := os.MkdirAll(parentPath, os.ModePerm); err != nil {
		return err
	}
	_, err := os.Create(filePath)
	return err
}

// 删除文件
func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}
