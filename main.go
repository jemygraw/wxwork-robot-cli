package main

import (
	"fmt"
	"github.com/jemygraw/wxwork-robot-cli/cmds"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd cobra.Command

func init() {
	rootCmd.AddCommand(&cmds.AddCmd)
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
