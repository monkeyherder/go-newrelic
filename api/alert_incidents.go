package api

import (
	"net/url"
)

func (c *Client) queryAlertIncidents() ([]AlertIncident, error) {
	incidents := []AlertIncident{}

	reqURL, err := url.Parse("/alerts_incidents.json")
	if err != nil {
		return nil, err
	}

	qs := reqURL.Query()

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

// ListAlertIncidents returns all alert incidents
func (c *Client) ListAlertIncidents() ([]AlertIncident, error) {
	return c.queryAlertIncidents()
}
