package cmd

import (
	"YTDownloaderCli/internal/yt"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {

	cmd := &cobra.Command{
		Short: "used to download YT videos",
		Run: func(cmd *cobra.Command, args []string) {
			runApplication()
		},
	}
	cmd.AddCommand(newRun())
	return cmd
}

func runApplication() {
	yt.Start()
}
