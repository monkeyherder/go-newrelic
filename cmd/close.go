package cmd

import "github.com/spf13/cobra"

var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "Actions that close NR entities",
}

func init() {
	RootCmd.AddCommand(closeCmd)
}
