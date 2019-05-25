package cmd

import (
	"github.com/mbamber/tvo/twitch"

	"github.com/spf13/cobra"
)

var (
	After string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the VODs for the given users.",
	Long:  "List the VODs for the given users.",
	Args:  cobra.MinimumNArgs(1),
	RunE:  listTwitchVODs,
}

func listTwitchVODs(cmd *cobra.Command, args []string) error {
	return twitch.ListTwitchVODs(args)
}
