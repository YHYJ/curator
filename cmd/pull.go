/*
File: pull.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-03-05 14:10：36

Description: 程序子命令'pull'时执行
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yhyj/curator/cli"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull changes from a remote repository",
	Long:  `Pull the latest changes from the origin remote and merge into the current branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		sourceFlag, _ := cmd.Flags().GetString("source")

		cli.RollingPullRepos(cfgFile, sourceFlag)
	},
}

func init() {
	var source string
	pullCmd.Flags().StringVarP(&source, "source", "s", "github", "Specify the data source (github or gitea)")

	pullCmd.Flags().BoolP("help", "h", false, "help for pull command")
	rootCmd.AddCommand(pullCmd)
}
