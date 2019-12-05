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

func (c *Client) postAlertIncident(id int, verb string) error {
	path := fmt.Sprintf("/alerts_incidents/%v/%v.json", id, verb)
	_, err := c.Do("POST", path, nil, nil)
	return err
}

// ListAlertIncidents returns all alert incidents
func (c *Client) ListAlertIncidents() ([]AlertIncident, error) {
	return c.queryAlertIncidents(false, false)
}

// ListOpenAlertIncidents returns open alert incidents
func (c *Client) ListOpenAlertIncidents() ([]AlertIncident, error) {
	return c.queryAlertIncidents(true, false)
}

func (c *Client) AcknowledgeAlertIncident(id int) error {
	return c.postAlertIncident(id, "acknowledge")
}

func (c *Client) CloseAlertIncident(id int) error {
	return c.postAlertIncident(id, "close")
}
