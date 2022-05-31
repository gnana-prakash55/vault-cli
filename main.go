/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"os"

	"github.com/gnana-prakash55/vault-cli/cmd"
)

func main() {
	os.Setenv("URL", "http://primevault.tech")
	cmd.Execute()
}
