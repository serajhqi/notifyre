package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "notifyre",
		Short: "Telegram notification gateway",
		Long:  "notifyre forwards authenticated HTTP requests to a Telegram channel via a bot.",
	}

	root.AddCommand(
		newServeCmd(),
		newSendCmd(),
		newSnippetCmd(),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
