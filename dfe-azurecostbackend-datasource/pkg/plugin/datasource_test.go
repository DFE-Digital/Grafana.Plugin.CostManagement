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

type CostData struct {
	Date      time.Time
	Cost      float64
	TotalCost float64
}

func getCostData() []CostData {
	costData := []CostData{
		{Date: parseDate("30/12/2023 00:00"), Cost: 12.30, TotalCost: 12.30},
		{Date: parseDate("31/12/2023 00:00"), Cost: 12.80, TotalCost: 25.20},
		{Date: parseDate("01/01/2024 00:00"), Cost: 12.30, TotalCost: 37.50},
		{Date: parseDate("02/01/2024 00:00"), Cost: 18.60, TotalCost: 56.00},
		// ... (add the remaining data)
		{Date: parseDate("29/01/2024 00:00"), Cost: 38.50, TotalCost: 850},
	}
	return costData
}


func TestLinearRegression(t *testing.T) {
	costData := getCostData()

	var dates []float64
	var values []float64
	var sumValues []float64

	for _, data := range costData {
		days := float64(data.Date.Sub(costData[0].Date).Hours() / 24)
		dates = append(dates, days)
		values = append(values, data.Cost)
		sumValues = append(sumValues, data.TotalCost)
	}

	mValues, cValues := linearRegressionCalculation(dates, values)
	mSumValues, cSumValues := linearRegressionCalculation(dates, sumValues)

	// Assuming some expected values for the test
	expectedValuesSlope := 0.867135
	expectedValuesIntercept := 12.656628
	expectedSumValuesSlope := 28.565974
	expectedSumValuesIntercept := -9.475015

	// Check if the calculated values are close to the expected values
	assertClose(t, mValues, expectedValuesSlope, "values slope")
	assertClose(t, cValues, expectedValuesIntercept, "values intercept")
	assertClose(t, mSumValues, expectedSumValuesSlope, "sum-values slope")
	assertClose(t, cSumValues, expectedSumValuesIntercept, "sum-values intercept")
}

func assertClose(t *testing.T, actual, expected float64, name string) {
	t.Helper()
	epsilon := 1e-6
	if actual < expected-epsilon || actual > expected+epsilon {
		t.Errorf("%s: expected %.6f, got %.6f", name, expected, actual)
	}
}

// MockFetchToken is a mock implementation of the fetchToken function
func MockFetchToken(config Config) (string, error) {
	// Implement your mock logic here
	// For example, return a predefined token and no error
	return "mocked_token", nil
}

// Helper function to check if two slices are equal
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}


func TestGetCalculatedCostData(t *testing.T) {
	timeRange := backend.TimeRange{
		From: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC),
	}

	datapoints := []DatePoint{
		{Date: "2024, 01, 01", Value: 12.3},
		{Date: "2024, 01, 02", Value: 18.6},
		{Date: "2024, 01, 03", Value: 33.6},
		{Date: "2024, 01, 04", Value: 29.3},
		{Date: "2024, 01, 07", Value: 12.3},
		{Date: "2024, 01, 06", Value: 12.3},
		{Date: "2024, 01, 07", Value: 12.9},
		// Add more datapoints as needed
	}

	calcDates, calcValues, calcTotalvalues := getCalculatedCostData(timeRange, datapoints)

	expectedDates := []string{
		"2024-01-01",
		"2024-01-02",
		"2024-01-03",
		"2024-01-04",
		"2024-01-05",
		"2024-01-06",
		"2024-01-07",
		"2024-01-08",
		"2024-01-09",
		"2024-01-10",
		"2024-01-11",
		"2024-01-12",
		"2024-01-13",
		"2024-01-14",
		
	}

	expectedValues := []float64{
		22.19642857142856, 21.04999999999999, 19.90357142857142, 18.757142857142853, 17.610714285714288, 16.464285714285715, 15.31785714285715, 14.171428571428581, 13.025000000000013, 11.878571428571444, 10.732142857142875, 9.585714285714309, 8.43928571428574, 7.292857142857171,
	}

	expectedTotalValues := []float64{
		18.157142857142837, 38.642857142857125, 59.12857142857142, 79.6142857142857, 100.1, 120.58571428571429, 141.07142857142858, 161.55714285714288, 182.04285714285717, 202.52857142857147, 223.01428571428576, 243.50000000000006, 263.9857142857143, 284.47142857142865,
	}

	// Compare calculated dates with expected dates
	if !slicesEqual(calcDates, expectedDates) {
		t.Errorf("Dates do not match. Got %v, expected %v", calcDates, expectedDates)
	}

	// Compare calculated values with expected values
	isOK := true
	for i := range expectedValues {
		if calcValues[i] != expectedValues[i] {
			isOK = false
		}
	}

	if !isOK {
		t.Errorf("Values do not match. Got %v, expected %v", calcValues, expectedValues)
	}

	isOK = true
	for i := range expectedValues {
		if calcTotalvalues[i] != expectedTotalValues[i] {
			isOK = false
		}
	}

	if !isOK {
		t.Errorf("Total values do not match. Got %v, expected %v", calcTotalvalues, expectedTotalValues)
	}
}
