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
	Short: "Used to clone user-specified repositories",
	Long:  `Clone-repos is used to clone the specified repository.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile = function.UserInfo.HomeDir + "/.config" + "/clone-repos/config.toml"

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "help for Clone-repos")

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFile, "Config file")
}
