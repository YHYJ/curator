/*
File: version.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-18 13:21:10

Description: 执行子命令 'version'
*/

package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/yhyj/curator/general"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print program version",
	Long:  `Print program version and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		onlyFlag, _ := cmd.Flags().GetBool("only")

		programInfo := general.ProgramInfo()

		if onlyFlag {
			color.Printf("%s\n", programInfo["Version"])
		} else {
			color.Printf("%s %s\n", general.LightText(programInfo["Name"]), general.LightText(programInfo["Version"]))
			color.Printf("%s %s\n", general.SecondaryText("Project:"), general.SecondaryText(programInfo["Project"]))
			color.Printf("%s %s\n", general.SecondaryText("Build rev:"), general.SecondaryText(programInfo["GitCommitHash"]))
			color.Printf("%s %s\n", general.SecondaryText("Built on:"), general.SecondaryText(programInfo["BuildTime"]))
			color.Printf("%s %s\n", general.SecondaryText("Built by:"), general.SecondaryText(programInfo["BuildBy"]))
		}
	},
}

func init() {
	versionCmd.Flags().BoolP("only", "", false, "Only print the version number, like 'v0.0.1'")

	versionCmd.Flags().BoolP("help", "h", false, "help for version command")
	rootCmd.AddCommand(versionCmd)
}
