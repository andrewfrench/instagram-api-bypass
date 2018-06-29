package main

import (
	"fmt"

	"github.com/andrewfrench/instagram-api-bypass/pkg/account"

	"github.com/spf13/cobra"
	"encoding/json"
)

var accountCmd = &cobra.Command{
	Use: "account [username]",
	Short: "Get account information",
	Long: "Provide the username to return detailed public account information.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("Requires one argument: username\n%s\n", cmd.UsageString())
			return
		}

		username := args[0]
		acc, err := account.Get(username)
		if err != nil {
			fmt.Printf("Unable to get account [%s]: %v\n", username, err)
			return
		}

		accountJson, err := json.MarshalIndent(acc, "", "  ")
		if err != nil {
			fmt.Printf("Unable to represent account [%s]: %v\n", username, err)
			return
		}

		fmt.Printf("%s\n", accountJson)

		return
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}