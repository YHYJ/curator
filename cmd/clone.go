/*
File: clone.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-21 15:45:25

Description: 执行子命令 'clone'
*/

package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/yhyj/curator/cli"
	"github.com/yhyj/curator/general"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone the specified repository",
	Long:  `Clone the repository specified in the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取配置文件路径
		configFile, _ := cmd.Flags().GetString("config")
		// 解析参数
		sourceFlag, _ := cmd.Flags().GetString("source")

		// 读取配置文件
		configTree, err := general.GetTomlConfig(configFile)
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
			return
		}

		// 使用指定的数据源进行克隆
		cli.RollingCloneRepos(configTree, sourceFlag)
	},
}

func init() {
	var source string
	cloneCmd.Flags().StringVarP(&source, "source", "s", "github", "Specify the data source (github or gitea)")

	cloneCmd.Flags().BoolP("help", "h", false, "help for clone command")
	rootCmd.AddCommand(cloneCmd)
}
