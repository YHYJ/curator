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
		cfgFile, _ := cmd.Flags().GetString("config")
		// 解析参数
		createFlag, _ := cmd.Flags().GetBool("create")
		forceFlag, _ := cmd.Flags().GetBool("force")
		printFlag, _ := cmd.Flags().GetBool("print")

		var (
			cfgFileNotFoundMessage = "Configuration file not found (use --create to create a configuration file)" // 配置文件不存在
		)

		// 检查配置文件是否存在
		cfgFileExist := general.FileExist(cfgFile)

		// 执行配置文件操作
		if createFlag {
			if cfgFileExist {
				if forceFlag {
					if err := general.DeleteFile(cfgFile); err != nil {
						color.Error.Println(err)
						return
					}
					if err := general.CreateFile(cfgFile); err != nil {
						color.Error.Println(err)
						return
					}
					_, err := cli.WriteTomlConfig(cfgFile)
					if err != nil {
						color.Error.Println(err)
						return
					}
					color.Printf("%s %s: %s\n", general.FgWhite("Create"), general.PrimaryText(cfgFile), general.SuccessText("file overwritten"))
				} else {
					color.Printf("%s %s: %s %s\n", general.FgWhite("Create"), general.PrimaryText(cfgFile), general.WarnText("file exists"), general.SecondaryText("(use --force to overwrite)"))
				}
			} else {
				if err := general.CreateFile(cfgFile); err != nil {
					color.Error.Println(err)
					return
				}
				_, err := cli.WriteTomlConfig(cfgFile)
				if err != nil {
					color.Error.Println(err)
					return
				}
				color.Printf("%s %s: %s\n", general.FgWhite("Create"), general.PrimaryText(cfgFile), general.SuccessText("file created"))
			}
		}

		if printFlag {
			if cfgFileExist {
				configTree, err := cli.GetTomlConfig(cfgFile)
				if err != nil {
					color.Error.Println(err)
				} else {
					color.Println(general.NoteText(configTree))
				}
			} else {
				color.Error.Println(cfgFileNotFoundMessage)
			}
		}
	},
}

func init() {
	configCmd.Flags().BoolP("create", "", false, "Create a default configuration file")
	configCmd.Flags().BoolP("force", "", false, "Overwrite existing configuration files")
	configCmd.Flags().BoolP("print", "", false, "Print configuration file content")

	configCmd.Flags().BoolP("help", "h", false, "help for config command")
	rootCmd.AddCommand(configCmd)
}
