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
			          35204408
			        ],
			        "policy_id": 28069
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
				Violations: []int{35204408},
				PolicyId:   28069,
			},
		},
	}

	alertIncidents, err := c.queryAlertIncidents()
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
