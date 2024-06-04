/*
File: config.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-24 16:39:13

Description: 执行子命令 'config'
*/

package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/yhyj/curator/cli"
	"github.com/yhyj/curator/general"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Operate configuration file",
	Long:  `Manipulate the program's configuration files, including generating and printing.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取配置文件路径
		configFile, _ := cmd.Flags().GetString("config")
		// 解析参数
		createFlag, _ := cmd.Flags().GetBool("create")
		forceFlag, _ := cmd.Flags().GetBool("force")
		openFlag, _ := cmd.Flags().GetBool("open")
		printFlag, _ := cmd.Flags().GetBool("print")

		var (
			noticeSlogan []string // 提示标语
		)
		// 检查参数
		if !createFlag && !printFlag && !openFlag {
			cmd.Help()
			noticeSlogan = append(noticeSlogan, "Please refer to the above help information")
		}

		// 创建配置文件流程
		if createFlag {
			cli.CreateConfigFile(configFile, forceFlag)
		}

		// 打开配置文件流程
		if openFlag {
			cli.OpenConfigFile(configFile)
		}

		// 打印配置文件流程
		if printFlag {
			cli.PrintConfigFile(configFile)
		}

		// 输出标语
		if len(noticeSlogan) > 0 {
			color.Println()
			for _, slogan := range noticeSlogan {
				color.Notice.Tips(general.PrimaryText(slogan))
			}
		}
	},
}

func init() {
	configCmd.Flags().BoolP("create", "", false, "Create a default configuration file")
	configCmd.Flags().BoolP("force", "", false, "Overwrite existing configuration files")
	configCmd.Flags().BoolP("open", "", false, "Open the configuration file with the default editor")
	configCmd.Flags().BoolP("print", "", false, "Print configuration file content")

	configCmd.Flags().BoolP("help", "h", false, "help for config command")
	rootCmd.AddCommand(configCmd)
}
