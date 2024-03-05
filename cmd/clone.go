/*
File: clone.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-21 15:45:25

Description: 程序子命令'clone'时执行
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yhyj/repos/cli"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone the specified repository",
	Long:  `Clone the repository specified in the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		sourceFlag, _ := cmd.Flags().GetString("source")

		cli.RollingCloneRepos(cfgFile, sourceFlag)
	},
}

func init() {
	var source string
	cloneCmd.Flags().StringVarP(&source, "source", "s", "github", "Specify the data source (github or gitea)")

	cloneCmd.Flags().BoolP("help", "h", false, "help for clone command")
	rootCmd.AddCommand(cloneCmd)
}
