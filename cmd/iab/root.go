package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "iab",
	Short: "Instagram API Bypass",
	Long: "A tool for accessing public account and media information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", cmd.UsageString())
	},
}
