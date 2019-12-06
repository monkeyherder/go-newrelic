package cmd

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
	"strconv"
	"time"

	"github.com/paultyng/go-newrelic/v4/api"
)

type HumanReadableAlertIncident struct {
	ID                 int               `json:"id,omitempty"`
	OpenedAt           time.Time         `json:"opened_at,omitempty"`
	ClosedAt           *time.Time         `json:"closed_at,omitempty"`
	IncidentPreference string            `json:"incident_preference,omitempty"`
	Links              api.AlertIncidentLink `json:"links"`
}

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

		format, err := cmd.Flags().GetString("format")

		if err != nil {
			fmt.Errorf("Could not get format flag: %w", err)
		}

		if format != "json" {
			newResources := []HumanReadableAlertIncident{}
			for _, incident := range resources {
				humanReadable :=
				HumanReadableAlertIncident{
					ID: incident.ID,
					OpenedAt: time.Unix(int64(incident.OpenedAt/1000), 0),
					IncidentPreference: incident.IncidentPreference,
					Links: incident.Links,
				}

				if incident.ClosedAt != 0 {
					closedTime := time.Unix(int64(incident.ClosedAt/1000),0)
					humanReadable.ClosedAt = &closedTime
				}

				newResources = append(newResources, humanReadable)
			}

			return outputList(cmd, newResources)
		}

		return outputList(cmd, resources)
	},
})

var closeIncidentsCmd = makeIncidentsCmd(cobra.Command{
	Use:  "incidents <id> [more id args]",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}

		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				return fmt.Errorf("invalid Incident ID %s: %v", arg, err)
			}
			fmt.Printf("Closing incident (%v -> %v)", arg, id)
			err = client.CloseAlertIncident(id)
			if err != nil {
				return fmt.Errorf("error closing Incident ID %s: %v", arg, err)
			}
		}

		return nil
	},
})

func init() {
	getCmd.AddCommand(getAlertIncidentsCmd)
	getAlertIncidentsCmd.Flags().BoolP("only-open", "o", false, "Excludes closed incidents if true. Default: false")
	getAlertIncidentsCmd.Flags().BoolP("exclude-violations", "x", false, "Excludes the linked violations from response if true. Default: false")

	closeCmd.AddCommand(closeIncidentsCmd)
}
