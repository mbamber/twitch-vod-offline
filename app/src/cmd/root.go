package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "tvo",
	Short:   "tvo (Twitch VOD Offline) is a tool for retrieving Twitch VODs.",
	Long:    "A tool for retrieving Twitch VODs for a set of provided users.",
	Version: "0.0.1",
}

func Execute() {
	rootCmd.SilenceUsage = true
	rootCmd.AddCommand(listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
