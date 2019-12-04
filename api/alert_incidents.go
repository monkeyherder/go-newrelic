package api

import (
	"fmt"
	"net/url"
	"strconv"
)

func (c *Client) queryAlertIncidents(policyID int) ([]AlertIncident, error) {
	conditions := []AlertIncident{}

	reqURL, err := url.Parse("/alerts_conditions.json")
	if err != nil {
		return nil, err
	}

	qs := reqURL.Query()
	qs.Set("policy_id", strconv.Itoa(policyID))

	reqURL.RawQuery = qs.Encode()

	nextPath := reqURL.String()

	for nextPath != "" {
		resp := struct {
			Incidents []AlertIncident `json:"conditions,omitempty"`
		}{}

		nextPath, err = c.Do("GET", nextPath, nil, &resp)
		if err != nil {
			return nil, err
		}

		for _, c := range resp.Incidents {
			c.PolicyID = policyID
		}

		conditions = append(conditions, resp.Incidents...)
	}

	return conditions, nil
}

// GetAlertIncident gets information about an alert condition given an ID and policy ID.
func (c *Client) GetAlertIncident(policyID int, id int) (*AlertIncident, error) {
	conditions, err := c.queryAlertIncidents(policyID)
	if err != nil {
		return nil, err
	}

	for _, condition := range conditions {
		if condition.ID == id {
			return &condition, nil
		}
	}

	return nil, ErrNotFound
}

// ListAlertIncidents returns alert conditions for the specified policy.
func (c *Client) ListAlertIncidents(policyID int) ([]AlertIncident, error) {
	return c.queryAlertIncidents(policyID)
}

// CreateAlertIncident creates an alert condition given the passed configuration.
func (c *Client) CreateAlertIncident(condition AlertIncident) (*AlertIncident, error) {
	policyID := condition.PolicyID

	req := struct {
		Incident AlertIncident `json:"condition"`
	}{
		Incident: condition,
	}

	resp := struct {
		Incident AlertIncident `json:"condition,omitempty"`
	}{}

	u := &url.URL{Path: fmt.Sprintf("/alerts_conditions/policies/%v.json", policyID)}
	_, err := c.Do("POST", u.String(), req, &resp)
	if err != nil {
		return nil, err
	}

	resp.Incident.PolicyID = policyID

	return &resp.Incident, nil
}

// UpdateAlertIncident updates an alert condition with the specified changes.
func (c *Client) UpdateAlertIncident(condition AlertIncident) (*AlertIncident, error) {
	policyID := condition.PolicyID
	id := condition.ID

	req := struct {
		Incident AlertIncident `json:"condition"`
	}{
		Incident: condition,
	}

	resp := struct {
		Incident AlertIncident `json:"condition,omitempty"`
	}{}

	u := &url.URL{Path: fmt.Sprintf("/alerts_conditions/%v.json", id)}
	_, err := c.Do("PUT", u.String(), req, &resp)
	if err != nil {
		return nil, err
	}

	resp.Incident.PolicyID = policyID

	return &resp.Incident, nil
}

// DeleteAlertIncident removes the alert condition given the specified ID and policy ID.
func (c *Client) DeleteAlertIncident(policyID int, id int) error {
	u := &url.URL{Path: fmt.Sprintf("/alerts_conditions/%v.json", id)}
	_, err := c.Do("DELETE", u.String(), nil, nil)
	return err
}
