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
				"conditions": [
					{
						"id": 1234,
						"type": "browser_metric",
						"name": "End User Apdex (Low)",
						"enabled": false,
						"entities": ["126408", "127809"],
						"metric": "end_user_apdex",
						"condition_scope": "application"
					}
				]
			}
			`))
	}))

	expected := []AlertIncident{
		{
			ID:       1234,
			Type:     "browser_metric",
			Name:     "End User Apdex (Low)",
			Enabled:  false,
			Entities: []string{"126408", "127809"},
			Metric:   "end_user_apdex",
			Scope:    "application",
		},
	}

	policyID := 123
	alertIncidents, err := c.queryAlertIncidents(policyID)
	if err != nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if alertIncidents == nil {
		t.Log(err)
		t.Fatal("GetAlertIncident error")
	}
	if diff := cmp.Diff(alertIncidents, expected); diff != "" {
		t.Fatalf("Alert conditions not parsed correctly: %s", diff)
	}
}
