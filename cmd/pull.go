/*
File: pull.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-03-05 14:10：36

Description: 执行子命令 'pull'
*/

package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/yhyj/curator/cli"
	"github.com/yhyj/curator/general"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Fetch from and merge with another repository or local branch",
	Long:  `Pull the latest changes from the origin remote and merge into the current branch.`,
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

		cli.RollingPullRepos(configTree, sourceFlag)
	},
}

func init() {
	var source string
	pullCmd.Flags().StringVarP(&source, "source", "s", "github", "Specify the data source (github or gitea)")

	pullCmd.Flags().BoolP("help", "h", false, "help for pull command")
	rootCmd.AddCommand(pullCmd)
}
