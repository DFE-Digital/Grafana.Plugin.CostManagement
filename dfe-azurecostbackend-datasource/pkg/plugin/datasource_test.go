package plugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestQueryData(t *testing.T) {
	ds := Datasource{}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A"},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}
}

func TestConvertCostsToDatePoint(t *testing.T) {
	// Mock CostResponse for testing
	costs := CostResponse{
		Properties: Properties{
			Rows: [][]interface{}{
				{10.0, float64(20220101)},
				{20.0, float64(20220102)},
			},
		},
	}

	datePoints, err := convertCostsToDatePoint(costs)
	if err != nil {
		t.Errorf("Error converting costs to DatePoint: %v", err)
	}

	expectedDatePoints := []DatePoint{
		{Date: "2022-01-01", Value: 10.0},
		{Date: "2022-01-02", Value: 20.0},
	}

	if len(datePoints) != len(expectedDatePoints) {
		t.Errorf("Expected %d date points, got %d", len(expectedDatePoints), len(datePoints))
	}

	for i, dp := range datePoints {
		if dp != expectedDatePoints[i] {
			t.Errorf("Expected DatePoint %+v, got %+v", expectedDatePoints[i], dp)
		}
	}
}

func TestFetchToken(t *testing.T) {
	config := Config{
		TenantID:  "your_tenant_id",
		ClientID:  "your_client_id",
		Password:  "your_client_secret",
	}

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request parameters
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		// if r.URL.String() != "https://login.microsoftonline.com/your_tenant_id/oauth2/v2.0/token" {
		// 	t.Errorf("Unexpected request URL: %s", r.URL.String())
		// }

		// Simulate the response based on the request parameters
		if r.FormValue("client_id") != "your_client_id" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Return a sample access token
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access_token": "sample_access_token"}`))
	}))

	// Close the server when the test finishes
	defer server.Close()

	// Override the token URL with the mock server URL
	config.TokenURL = server.URL

	// Call the fetchToken function
	token, err := fetchToken(config)

	// Check the result
	if err != nil {
		t.Errorf("fetchToken returned an error: %v", err)
	}
	if token != "sample_access_token" {
		t.Errorf("fetchToken returned an unexpected token: %s", token)
	}
}