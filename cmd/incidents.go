package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"

	"github.com/paultyng/go-newrelic/v4/api"
)

type HumanReadableAlertIncident struct {
	ID                 int
	OpenedAt           time.Time
	ClosedAt           *time.Time
	TimeOpen           time.Duration
	IncidentPreference string
	Links              api.AlertIncidentLink
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
			return fmt.Errorf("Could not get only-open flag: %w", err)
		}

		exclude_violations, err := cmd.Flags().GetBool("exclude-violations")

		if err != nil {
			return fmt.Errorf("Could not get only-open flag: %w", err)
		}

		resources, err := client.ListAlertIncidents(only_open, exclude_violations)
		if err != nil {
			return err
		}

		format, err := cmd.Flags().GetString("format")

		if err != nil {
			return fmt.Errorf("Could not get format flag: %w", err)
		}

		if format != "json" {
			formattedResources := formatIncidents(resources)

			return outputList(cmd, formattedResources)
		}

		return outputList(cmd, resources)
	},
})

func formatIncidents(resources []api.AlertIncident) []HumanReadableAlertIncident {
	currentEpochNanos := time.Now().UnixNano()

	newResources := []HumanReadableAlertIncident{}
	for _, incident := range resources {
		humanReadable :=
			HumanReadableAlertIncident{
				ID:                 incident.ID,
				OpenedAt:           time.Unix(int64(incident.OpenedAt/1000), 0),
				IncidentPreference: incident.IncidentPreference,
				Links:              incident.Links,
			}

		if incident.ClosedAt != 0 {
			closedTime := time.Unix(int64(incident.ClosedAt/1000), 0)
			humanReadable.ClosedAt = &closedTime
			humanReadable.TimeOpen = time.Duration(int64(incident.ClosedAt - incident.OpenedAt))
		}

		humanReadable.TimeOpen = time.Duration((currentEpochNanos/1000000 - int64(incident.OpenedAt)) * 1000000)
		newResources = append(newResources, humanReadable)
	}

	return newResources
}

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
