/*
File: root.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 13:16:00

Description: 执行程序
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yhyj/curator/general"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "curator",
	Short: "My code repository curator",
	Long:  `Responsible for managing my code repository.`,
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

func init() {
	rootCmd.PersistentFlags().String("config", general.ConfigFile, "Specify configuration file")

	rootCmd.Flags().BoolP("help", "h", false, "help for curator")
}
