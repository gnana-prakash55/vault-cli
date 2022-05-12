/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// newrepoCmd represents the newrepo command
var newrepoCmd = &cobra.Command{
	Use:   "newrepo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// repoName, _ := cmd.Flags().GetString("repo")

	},
}

func init() {
	rootCmd.AddCommand(newrepoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	newrepoCmd.PersistentFlags().StringP("repo", "r", "", "Repository Name")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newrepoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
