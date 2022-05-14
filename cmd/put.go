/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gnana-prakash55/vault-cli/utils"
	"github.com/spf13/cobra"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := os.Getwd()

		token, err := utils.ReadToken()
		if err != nil {
			log.Fatalf("Unable to Login...")
		}

		resp, err := utils.UploadFiles(path, token)
		if err != nil {
			log.Fatalln("Unable to Upload...", err.Error())
		}

		fmt.Println(resp)
		// path, _ := os.Getwd()
		// utils.UploadFiles(path)

	},
}

func init() {
	rootCmd.AddCommand(putCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// putCmd.PersistentFlags().StringP("filename", "f", "", "Filename to Upload")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// putCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
