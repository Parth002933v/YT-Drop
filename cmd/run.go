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

			// p := tea.NewProgram(yt.InitialModel())
			// if _, err := p.Run(); err != nil {
			// 	log.Fatal(err)
			// }

		},
	}
	return cmd
}
