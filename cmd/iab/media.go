package main

import (
	"encoding/json"
	"fmt"

	"github.com/andrewfrench/instagram-api-bypass/pkg/media"

	"github.com/spf13/cobra"
)

var mediaCmd = &cobra.Command{
	Use: "media [media-shortcode]",
	Short: "Get media information",
	Long: "Provide the media shortcode to get detailed public media information.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("Requires one argument: media shortcode\n%s\n", cmd.UsageString())
			return
		}

		mediaShortcode := args[1]
		med, err := media.Get(mediaShortcode)
		if err != nil {
			fmt.Printf("Unable to get media [%s]: %v\n", mediaShortcode, err)
			return
		}

		mediaJson, err := json.MarshalIndent(med, "", "  ")
		if err != nil {
			fmt.Printf("Unable to represent media [%s]: %v\n", mediaShortcode, err)
			return
		}

		fmt.Printf("%s\n", mediaJson)

		return
	},
}

func init() {
	rootCmd.AddCommand(mediaCmd)
}
