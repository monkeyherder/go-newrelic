package cmd

import (
	"github.com/imdario/mergo"
	"fmt"
	"github.com/spf13/cobra"
)

func makeIncidentsCmd(dst cobra.Command) *cobra.Command {
	src := cobra.Command{
		Use: "incidents",
	}
	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}
	return &dst
}

var getAlertIncidentsCmd = makeIncidentsCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}

		only_open, err := cmd.Flags().GetBool("only-open")

		if err != nil {
			fmt.Errorf("Could not get only-open flag: %w", err)
		}

		exclude_violations, err := cmd.Flags().GetBool("exclude-violations")

		if err != nil {
			fmt.Errorf("Could not get only-open flag: %w", err)
		}

		resources, err := client.ListAlertIncidents(only_open, exclude_violations)
		if err != nil {
			return err
		}
		return outputList(cmd, resources)
	},
})

func init() {
	getCmd.AddCommand(getAlertIncidentsCmd)
	getAlertIncidentsCmd.Flags().BoolP("only-open", "o", false, "Excludes closed incidents if true. Default: false")
	getAlertIncidentsCmd.Flags().BoolP("exclude-violations", "x", false, "Excludes the linked violations from response if true. Default: false")

}
