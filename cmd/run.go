/*
File: run.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-21 15:45:25

Description: 程序子命令'run'时执行
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yhyj/clone-repos/cli"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start Cloning",
	Long:  `Start cloning.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		sourceFlag, _ := cmd.Flags().GetString("source")

		cli.RollingCloneRepos(cfgFile, sourceFlag)
	},
}

func init() {
	var source string
	runCmd.Flags().StringVarP(&source, "source", "s", "github", "Specify the data source (github or gitea)")

	runCmd.Flags().BoolP("help", "h", false, "help for run command")
	rootCmd.AddCommand(runCmd)
}
