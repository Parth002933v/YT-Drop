package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "short test of how to use it",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("test...")
			
		},
	}
	return cmd
}
