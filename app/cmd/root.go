package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"tg_bot/rest_api"
)

func Root(f *rest_api.Server) *cobra.Command {
	root := cobra.Command{
		Use:   "api",
		Short: "api application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Type a 'api help' for usage details")
		},
	}

	root.AddCommand(
		start(f),
	)
	return &root
}
