package app

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docgen",
		Short: "Scan backend code and generate API/operation documents",
	}

	cmd.AddCommand(newScanCmd())
	cmd.AddCommand(newManualCmd())
	cmd.AddCommand(newServeCmd())
	return cmd
}
