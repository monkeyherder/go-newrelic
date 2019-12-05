package api

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestListAlertIncidents(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"incidents": [
			    {
			      "id": 42,
			      "opened_at": 1575502560942,
			      "incident_preference": "PER_CONDITION",
			      "links": {
			        "violations": [
			          123456789
			        ],
			        "policy_id": 12345
				  }
				},
				{
					"id": 24,
					"opened_at": 1575506284796,
					"closed_at": 1575506342161,
					"incident_preference": "PER_POLICY",
					"links": {
						"violations": [
						987654321
						],
						"policy_id": 54321
					}
				}
				]
			}
			`))
	}))

	expected := []AlertIncident{
		{
			ID:                 42,
			OpenedAt:           1575502560942,
			IncidentPreference: "PER_CONDITION",
			Links: AlertIncidentLink{
				Violations: []int{123456789},
				PolicyId:   12345,
			},
		},
		{
			ID:                 24,
			OpenedAt:           1575506284796,
			ClosedAt:           1575506342161,
			IncidentPreference: "PER_POLICY",
			Links: AlertIncidentLink{
				Violations: []int{987654321},
				PolicyId:   54321,
			},
		},
	}

	alertIncidents, err := c.ListAlertIncidents()
	if err != nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if alertIncidents == nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if diff := cmp.Diff(alertIncidents, expected); diff != "" {
		t.Fatalf("Alert incidents not parsed correctly: %s", diff)
	}
}

func TestOpenListAlertIncidents(t *testing.T) {
	c := newTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			{
				"incidents": [
			    {
			      "id": 42,
			      "opened_at": 1575502560942,
			      "incident_preference": "PER_CONDITION",
			      "links": {
			        "violations": [
			          123456789
			        ],
			        "policy_id": 12345
				  }
				}
				]
			}
			`))
	}))

	expected := []AlertIncident{
		{
			ID:                 42,
			OpenedAt:           1575502560942,
			IncidentPreference: "PER_CONDITION",
			Links: AlertIncidentLink{
				Violations: []int{123456789},
				PolicyId:   12345,
			},
		},
	}

	alertIncidents, err := c.ListOpenAlertIncidents()
	if err != nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if alertIncidents == nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if diff := cmp.Diff(alertIncidents, expected); diff != "" {
		t.Fatalf("Alert incidents not parsed correctly: %s", diff)
	}
}