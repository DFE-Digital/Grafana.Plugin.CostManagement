package plugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	"encoding/json"
	"reflect"

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

func TestGetCosts(t *testing.T) {
	// Create a new instance of the Config struct
	config := Config{
		SubscriptionID:            "your-subscription-id",
		AzureCostSubscriptionUrl:  "https://your-cost-management-url.com/",
	}

	costResponse := createMockCostResponse()
	responseJSON, err := json.Marshal(costResponse)
	if err != nil {
		t.Errorf("Error marshalling response: %v", err)
	}
	
	// Create a new instance of the httptest.Server struct
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method and URL
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		// if r.URL.String() != "/your-subscription-id/providers/Microsoft.CostManagement/query?api-version=2023-03-01" {
		// 	t.Errorf("Expected URL /your-subscription-id/providers/Microsoft.CostManagement/query?api-version=2023-03-01, got %s", r.URL.String())
		// }

		// Write a mock response
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}))

	// Close the server when the test is done
	defer server.Close()

	// Call the getCosts function with the mock token and config
	token := "your-mock-token"
	config.TokenURL = server.URL
	start, end := getCurrentYearDates()
	costs, err := getCosts(token, config, start, end, "resourceid")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	id := costs.ID
	if (id == ""){
		t.Errorf("Expected response %v, got %v", id, costs.ID)
	}

	// Check the response
	expected := costResponse
	if !reflect.DeepEqual(costs, expected) {
		t.Errorf("Expected response %v, got %v", expected, costs)
	}
}

func TestGetForecast(t *testing.T) {
	// Create a new instance of the Config struct
	config := Config{
		SubscriptionID:            "your-subscription-id",
		AzureCostSubscriptionUrl:  "https://your-cost-management-url.com/",
	}

	costResponse := createMockCostResponse()
	responseJSON, err := json.Marshal(costResponse)
	if err != nil {
		t.Errorf("Error marshalling response: %v", err)
	}
	
	// Create a new instance of the httptest.Server struct
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method and URL
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		// if r.URL.String() != "/your-subscription-id/providers/Microsoft.CostManagement/query?api-version=2023-03-01" {
		// 	t.Errorf("Expected URL /your-subscription-id/providers/Microsoft.CostManagement/query?api-version=2023-03-01, got %s", r.URL.String())
		// }

		// Write a mock response
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}))

	// Close the server when the test is done
	defer server.Close()

	// Call the getCosts function with the mock token and config
	token := "your-mock-token"
	config.TokenURL = server.URL
	start, end := getCurrentYearDates()
	costs, err := getForecast(token, config, start, end, "resourceid")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	id := costs.ID
	if (id == ""){
		t.Errorf("Expected response %v, got %v", id, costs.ID)
	}

	// Check the response
	expected := costResponse
	if !reflect.DeepEqual(costs, expected) {
		t.Errorf("Expected response %v, got %v", expected, costs)
	}
}

func createMockCostResponse() CostResponse {
	// Define the properties of the CostResponse object
	properties := Properties{
		NextLink: "nextLink",
		Columns: []Column{
			{Name: "column1", Type: "type1"},
			{Name: "column2", Type: "type2"},
		},
		Rows: [][]interface{}{
			{1.0, "row1"},
			{2.0, "row2"},
		},
	}

	// Create a new instance of the CostResponse struct
	costResponse := CostResponse{
		ID:         "id1",
		Name:       "name1",
		Type:       "type1",
		Location:   "location1",
		SKU:        "sku1",
		ETag:       "eTag1",
		Properties: properties,
	}

	return costResponse
}

var timeNow = time.Now

func TestGetCurrentYearDates(t *testing.T) {
	// Mock current date for testing
	mockDate := time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC)
	// Save the original time function and replace it with a mock
	originalTimeNow := timeNow
	timeNow = func() time.Time { return mockDate }
	defer func() { timeNow = originalTimeNow }()

	// Call the function to get the formatted dates
	firstOfJanuary, thirtyFirstOfDecember := getCurrentYearDates()

	// Expected results for the mock date
	expectedFirstOfJanuary := strconv.Itoa(time.Now().Year()) +  "-01-01"
	expectedThirtyFirstOfDecember := strconv.Itoa(time.Now().Year()) + "-12-31"

	// Check if the actual results match the expected results
	if firstOfJanuary != expectedFirstOfJanuary {
		t.Errorf("First of January: expected %s, got %s", expectedFirstOfJanuary, firstOfJanuary)
	}

	if thirtyFirstOfDecember != expectedThirtyFirstOfDecember {
		t.Errorf("Thirty-First of December: expected %s, got %s", expectedThirtyFirstOfDecember, thirtyFirstOfDecember)
	}
}