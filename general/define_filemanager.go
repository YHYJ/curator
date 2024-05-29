/*
File: define_filemanager.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-24 16:41:17

Description: 文件管理
*/

package general

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gookit/color"
)

// ReadFileLine 读取文件指定行
//
// 参数：
//   - file: 文件路径
//   - line: 行号
//
// 返回：
//   - 指定行的内容
func ReadFileLine(file string, line int) string {
	// 打开文件
	text, err := os.Open(file)
	if err != nil {
		color.Danger.Println(err)
	}
	defer text.Close()

	// 创建一个扫描器对象按行遍历
	scanner := bufio.NewScanner(text)
	// 行计数
	count := 1
	// 逐行读取，输出指定行
	for scanner.Scan() {
		if line == count {
			return scanner.Text()
		}
		count++
	}
	return ""
}

// ReadFileKey 读取文件包含关键字的行
//
// 参数：
//   - file: 文件路径
//   - key: 关键字
//
// 返回：
//   - 包含关键字的行的内容
func ReadFileKey(file, key string) string {
	// 打开文件
	text, err := os.Open(file)
	if err != nil {
		color.Danger.Println(err)
	}
	defer text.Close()

	// 创建一个扫描器对象按行遍历
	scanner := bufio.NewScanner(text)
	// 逐行读取，输出指定行
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), key) {
			return scanner.Text()
		}
	}
	return ""
}

// ReadFileCount 获取文件包含关键字的行的计数
//
// 参数：
//   - file: 文件路径
//   - key: 关键字
//
// 返回：
//   - 包含关键字的行的数量
func ReadFileCount(file, key string) int {
	// 打开文件
	text, err := os.Open(file)
	if err != nil {
		color.Danger.Println(err)
	}
	defer text.Close()

	// 创建一个扫描器对象按行遍历
	scanner := bufio.NewScanner(text)
	// 计数器
	count := 0
	// 逐行读取，输出指定行
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), key) {
			count++
		}
	}
	return count
}

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

// GetAbsPath 获取指定文件的绝对路径
//
// 参数：
//   - filePath: 文件路径
//
// 返回：
//   - 文件的绝对路径
func GetAbsPath(filePath string) string {
	// 获取绝对路径
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return ""
	} else {
		return absPath
	}
}

// FileEmpty 判断文件是否为空
//
//   - 无法判断文件夹
//
// 参数：
//   - filePath: 文件路径
//
// 返回：
//   - 文件为空返回 true，否则返回 false
func FileEmpty(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return true
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return true
	}
	return fi.Size() == 0
}

// FolderEmpty 判断文件夹是否为空
//
//   - 包括隐藏文件
//
// 参数：
//   - dirPath: 文件夹路径
//
// 返回：
//   - 文件夹为空返回 true，否则返回 false
func FolderEmpty(dirPath string) bool {
	file, err := os.Open(dirPath)
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

// CreateFile 创建文件，包括其父目录
//
// 参数：
//   - filePath: 文件路径
//
// 返回：
//   - 错误信息
func CreateFile(filePath string) error {
	if FileExist(filePath) {
		return nil
	}
	// 创建父目录
	parentPath := filepath.Dir(filePath)
	if err := os.MkdirAll(parentPath, os.ModePerm); err != nil {
		return err
	}
	// 创建文件
	if _, err := os.Create(filePath); err != nil {
		return err
	}

	return nil
}

// CreateDir 创建文件夹
//
// 参数：
//   - dirPath: 文件夹路径
//
// 返回：
//   - 错误信息
func CreateDir(dirPath string) error {
	if FileExist(dirPath) {
		return nil
	}
	return os.MkdirAll(dirPath, os.ModePerm)
}

// GoToDir 进到指定文件夹
//
// 参数：
//   - dirPath: 文件夹路径
//
// 返回：
//   - 错误信息
func GoToDir(dirPath string) error {
	return os.Chdir(dirPath)
}

// WriteFile 写入内容到文件
//
// 参数：
//   - filePath: 文件路径
//   - content: 内容
//
// 返回：
//   - 错误信息
func WriteFile(filePath string, content string) error {
	// 文件存在
	if FileExist(filePath) {
		if FileEmpty(filePath) { // 文件内容为空
			// 打开文件并写入内容
			file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				return err
			} else {
				_, err := file.WriteString(content)
				if err != nil {
					return err
				}
			}
		} else { // 文件内容不为空
			return fmt.Errorf("File %s is not empty", filePath)
		}
	} else {
		// 文件不存在，创建文件
		if err := CreateFile(filePath); err != nil {
			return err
		}
		// 打开文件并写入内容
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			return err
		} else {
			_, err := file.WriteString(content)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteFile 删除文件
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
	return os.Remove(filePath)
}

// CompareFile 并发比较两个文件是否相同
//
// 参数：
//   - file1Path: 文件1路径
//   - file2Path: 文件2路径
//
// 返回：
//   - 文件相同返回 true，出错或不同返回 false
func CompareFile(file1Path string, file2Path string) (bool, error) {
	// 尝试打开文件
	file1, err := os.Open(file1Path)
	if err != nil {
		return false, err
	}
	defer file1.Close()
	file2, err := os.Open(file2Path)
	if err != nil {
		return false, err
	}
	defer file2.Close()

	// 获取文件大小
	file1Info, err := file1.Stat()
	if err != nil {
		return false, err
	}
	file2Info, err := file2.Stat()
	if err != nil {
		return false, err
	}
	file1Size := file1Info.Size()
	file2Size := file2Info.Size()

	// 如果文件大小不同则直接返回
	if file1Size != file2Size {
		return false, nil
	}

	// 文件大小相同则（按块）比较文件内容
	const chunkSize = 1024 * 1024                             // 每次比较的块大小（1MB）
	numChunks := int((file1Size + chunkSize - 1) / chunkSize) // 计算文件需要被分成多少块

	equal := true                // 文件是否相同的标志位
	var wg sync.WaitGroup        // wg 用于等待所有的 goroutine 执行完毕，然后关闭 errCh 通道
	errCh := make(chan error, 1) // errCh 用于接收 goroutine 执行过程中返回的错误

	for i := 0; i < numChunks; i++ { // 同时比较多个块
		wg.Add(1)
		go func(chunkIndex int) {
			defer wg.Done()

			// 计算当前块的偏移量和大小
			offset := int64(chunkIndex) * chunkSize
			size := chunkSize
			if offset+int64(size) > file1Size {
				size = int(file1Size - offset)
			}

			// 创建两个大小为 size 的 buffer
			buffer1 := make([]byte, size)
			buffer2 := make([]byte, size)

			// 从文件中读取指定大小的内容到 buffer
			_, err := file1.ReadAt(buffer1, offset)
			if err != nil && err != io.EOF {
				errCh <- err
				return
			}

			// 从文件中读取指定大小的内容到 buffer
			_, err = file2.ReadAt(buffer2, offset)
			if err != nil && err != io.EOF {
				errCh <- err
				return
			}

			// 比较两个 buffer 是否相同
			if !bytesEqual(buffer1, buffer2) {
				equal = false
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return false, err
		}
	}

	return equal, nil
}

// bytesEqual 比较两个文件的内容
//
// 参数：
//   - b1: 文件1内容
//   - b2: 文件2内容
//
// 返回：
//   - 相同返回 true，不同返回 false
func bytesEqual(b1 []byte, b2 []byte) bool {
	if len(b1) != len(b2) {
		return false
	}

	for i := 0; i < len(b1); i++ {
		if b1[i] != b2[i] {
			return false
		}
	}

	return true
}
