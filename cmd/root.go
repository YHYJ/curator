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
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yhyj/repos/general"
)

var rootCmd = &cobra.Command{
	Use:   "repos",
	Short: "Used to clone user-specified repositories",
	Long:  `repos is used to clone the specified repository.`,
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

var cfgFile = filepath.Join(general.UserInfo.HomeDir, ".config", "repos", "config.toml")

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "help for repos")

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFile, "Specify configuration file")
}
