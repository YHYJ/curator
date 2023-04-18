/*
File: root.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 13:16:00

Description: 程序未带子命令或参数时执行
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yhyj/clone-repos/function"
)

var rootCmd = &cobra.Command{
	Use:   "clone-repos",
	Short: "用于克隆用户指定仓库",
	Long:  `Clone-repos用于克隆指定用户的指定仓库`,
	Run: func(cmd *cobra.Command, args []string) {
		function.RollingCLoneRepos(cfgFile)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var varHome = function.GetVariable("HOME")
var cfgFile = varHome+"/.config/clone-repos/config.toml"

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "Help for Rolling")

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFile, "config file")
}
