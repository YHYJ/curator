/*
File: version.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 13:21:10

Description: 程序子命令'version'时执行
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yhyj/clone-repos/function"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印程序版本",
	Long:  `显示版本信息并退出`,
	Run: func(cmd *cobra.Command, args []string) {
		programInfo := function.ProgramInfo()
		fmt.Printf(programInfo)
	},
}

func init() {
	versionCmd.Flags().BoolP("help", "h", false, "Help for version")
	rootCmd.AddCommand(versionCmd)
}
