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
	"github.com/yhyj/clone-repos/function"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start Cloning",
	Long:  `Start cloning.`,
	Run: func(cmd *cobra.Command, args []string) {
		function.RollingCLoneRepos(cfgFile)
	},
}

func init() {
	runCmd.Flags().BoolP("help", "h", false, "help for run command")
	rootCmd.AddCommand(runCmd)
}
