/*
File: config.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-24 16:39:13

Description: 程序子命令'config'时执行
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yhyj/clone-repos/function"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Operate configuration file",
	Long:  `Operate configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取配置文件路径
		cfgFile, _ := cmd.Flags().GetString("config")
		// 解析参数
		createFlag, _ := cmd.Flags().GetBool("create")
		forceFlag, _ := cmd.Flags().GetBool("force")
		printFlag, _ := cmd.Flags().GetBool("print")

		// 检查配置文件是否存在
		cfgFileExist := function.FileExist(cfgFile)

		// 执行配置文件操作
		if createFlag {
			if cfgFileExist {
				if forceFlag {
					function.DeleteFile(cfgFile)
					function.CreateFile(cfgFile)
					function.WriteTomlConfig(cfgFile)
					fmt.Printf("Create \x1b[33;1m%s\x1b[0m: file overwritten\n", cfgFile)
				} else {
					fmt.Printf("Create \x1b[33m%s\x1b[0m: file exists (use --force to overwrite)\n", cfgFile)
				}
			} else {
				err := function.CreateFile(cfgFile)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
					return
				}
				_, err = function.WriteTomlConfig(cfgFile)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
					return
				}
				fmt.Printf("Create \x1b[33;1m%s\x1b[0m: file created\n", cfgFile)
			}
		}

		if printFlag {
			if cfgFileExist {
				configTree, err := function.GetTomlConfig(cfgFile)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
				} else {
					fmt.Println(configTree)
				}
			} else {
				fmt.Printf("\x1b[31m%s\x1b[0m\n", "Configuration file not found (use --create to create a configuration file)")
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
