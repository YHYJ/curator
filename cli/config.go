/*
File: config.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-24 16:41:33

Description: 子命令 'config' 的实现
*/

package cli

import (
	"github.com/gookit/color"
	"github.com/yhyj/curator/general"
)

// CreateConfigFile 创建配置文件
//
// 参数：
//   - configFile: 配置文件路径
//   - reWrite: 是否覆写
func CreateConfigFile(configFile string, reWrite bool) {
	// 检查配置文件是否存在
	fileExist := general.FileExist(configFile)

	// 检测并创建配置文件
	if fileExist {
		if reWrite {
			if err := general.DeleteFile(configFile); err != nil {
				fileName, lineNo := general.GetCallerInfo()
				color.Danger.Printf("Delete file error (%s:%d): %s\n", fileName, lineNo+1, err)
				return
			}
			if err := general.CreateFile(configFile); err != nil {
				fileName, lineNo := general.GetCallerInfo()
				color.Danger.Printf("Create file error (%s:%d): %s\n", fileName, lineNo+1, err)
				return
			}
			_, err := general.WriteTomlConfig(configFile)
			if err != nil {
				fileName, lineNo := general.GetCallerInfo()
				color.Danger.Printf("Write config error (%s:%d): %s\n", fileName, lineNo+1, err)
				return
			}
			color.Printf("%s %s: %s\n", general.FgWhiteText("Create"), general.PrimaryText(configFile), general.SuccessText("file overwritten"))
		} else {
			color.Printf("%s %s: %s %s\n", general.FgWhiteText("Create"), general.PrimaryText(configFile), general.WarnText("file exists"), general.SecondaryText("(use --force to overwrite)"))
		}
	} else {
		if err := general.CreateFile(configFile); err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Danger.Printf("Create file error (%s:%d): %s\n", fileName, lineNo+1, err)
			return
		}
		_, err := general.WriteTomlConfig(configFile)
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Danger.Printf("Write config error (%s:%d): %s\n", fileName, lineNo+1, err)
			return
		}
		color.Printf("%s %s: %s\n", general.FgWhiteText("Create"), general.PrimaryText(configFile), general.SuccessText("file created"))
	}
}

// OpenConfigFile 打开配置文件
//
// 参数：
//   - configFile: 配置文件路径
func OpenConfigFile(configFile string) {
	// 检查配置文件是否存在
	fileExist := general.FileExist(configFile)

	if fileExist {
		editor := general.GetVariable("EDITOR")
		if editor == "" {
			editor = "vim"
			err := general.RunCommand(editor, []string{configFile})
			if err != nil {
				editor = "vi"
				err = general.RunCommand(editor, []string{configFile})
				if err != nil {
					fileName, lineNo := general.GetCallerInfo()
					color.Danger.Printf("Open config file error (%s:%d): %s\n", fileName, lineNo+1, err)
				}
			}
		} else {
			err := general.RunCommand(editor, []string{configFile})
			if err != nil {
				fileName, lineNo := general.GetCallerInfo()
				color.Danger.Printf("Open config file error (%s:%d): %s\n", fileName, lineNo+1, err)
			}
		}
	}
}

// PrintConfigFile 打印配置文件内容
//
// 参数：
//   - configFile: 配置文件路径
func PrintConfigFile(configFile string) {
	// 检查配置文件是否存在
	fileExist := general.FileExist(configFile)

	var (
		configFileNotFoundMessage = "Configuration file not found (use --create to create a configuration file)" // 配置文件不存在
	)

	if fileExist {
		configTree, err := general.GetTomlConfig(configFile)
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Danger.Printf("Get config error (%s:%d): %s\n", fileName, lineNo+1, err)
		} else {
			color.Println(general.PrimaryText(configTree))
		}
	} else {
		color.Danger.Println(configFileNotFoundMessage)
	}
}
