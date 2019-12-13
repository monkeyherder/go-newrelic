package api

import (
	"fmt"
	"net/url"
)

func (c *Client) queryAlertIncidents(onlyOpen bool, excludeViolations bool) ([]AlertIncident, error) {
	incidents := []AlertIncident{}

	reqURL, err := url.Parse("/alerts_incidents.json")
	if err != nil {
		return nil, err
	}

	qs := reqURL.Query()
	if onlyOpen {
		qs.Set("only_open", "true")
	}
	if excludeViolations {
		qs.Set("exclude_violations", "true")
	}

	reqURL.RawQuery = qs.Encode()

	nextPath := reqURL.String()

	for nextPath != "" {
		resp := struct {
			Incidents []AlertIncident `json:"incidents,omitempty"`
		}{}

		nextPath, err = c.Do("GET", nextPath, nil, &resp)
		if err != nil {
			return nil, err
		}

		incidents = append(incidents, resp.Incidents...)
	}

	return incidents, nil
}

func (c *Client) updateAlertIncident(id int, verb string) error {
	path := fmt.Sprintf("/alerts_incidents/%v/%v.json", id, verb)
	_, err := c.Do("PUT", path, nil, nil)
	return err
}

// ListAlertIncidents returns all alert incidents
func (c *Client) ListAlertIncidents(only_open bool, exclude_violations bool) ([]AlertIncident, error) {
	return c.queryAlertIncidents(only_open, exclude_violations)
}

func (c *Client) AcknowledgeAlertIncident(id int) error {
	return c.updateAlertIncident(id, "acknowledge")
}

func (c *Client) CloseAlertIncident(id int) error {
	return c.updateAlertIncident(id, "close")
}
