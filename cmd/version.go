package cmd

import (
	"github.com/spf13/cobra"
	"runtime"
)

func NewVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Check the version info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("- YT Drop %s\n", "1.0.6")
			cmd.Printf("- os/type: %s\n", runtime.GOOS)
			cmd.Printf("- os/arch: %s\n", runtime.GOARCH)
			cmd.Printf("- go/version: %s\n", runtime.Version())
			cmd.Printf("\n check out the link for all the new releases: %s\n", "https://github.com/Parth002933v/YT-Drop/releases")
			return nil
		},
	}
}
