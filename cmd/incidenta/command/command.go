package command

import "github.com/spf13/cobra"

func NewRootCmd(_ []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "incidenta",
		SilenceUsage: true,
	}

	rootCmd.AddCommand(
		newWebCmd(),
		newDemoCmd(),
	)

	return rootCmd
}
