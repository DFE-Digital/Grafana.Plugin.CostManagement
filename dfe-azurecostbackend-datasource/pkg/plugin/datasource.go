package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// Config represents the configuration structure.
type Config struct {
	AzureCostSubscriptionUrl string `json:"AzureCostSubscriptionUrl"`
	Password                 string `json:"password"`
	ClientID                 string `json:"clientId"`
	TenantID                 string `json:"tenantId"`
	SubscriptionID           string `json:"subscriptionId"`
	Region                   string `json:"region"`
	TokenURL                 string `json:"tokenUrl"`
}

// CostResponse represents the response structure
type CostResponse struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Location   interface{} `json:"location"`
	SKU        interface{} `json:"sku"`
	ETag       interface{} `json:"eTag"`
	Properties Properties  `json:"properties"`
}

// Properties represents the 'properties' structure
type Properties struct {
	NextLink string          `json:"nextLink"`
	Columns  []Column        `json:"columns"`
	Rows     [][]interface{} `json:"rows"`
}

// Column represents the 'column' structure
type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// DatePoint represents a point in time with a corresponding value
type DatePoint struct {
	Date string `json:"date"`
	//Date  time.Time `json:"date"`
	Value float64 `json:"value"`
}

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

func getSetting(settings backend.DataSourceInstanceSettings, requiredSetting string) (string, error) {
	settingValue, found := settings.DecryptedSecureJSONData[requiredSetting]
	if !found {
		return "", fmt.Errorf("could not locate " + requiredSetting)
	}
	log.Println(settingValue)

	return settingValue, nil
}

// NewDatasource creates a new datasource instance...
func NewDatasource(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {

	Password, found := settings.DecryptedSecureJSONData["Password"]
	if !found {
		return nil, fmt.Errorf("could not locate Password in NewDatasource ")
	}
	log.Println(Password)

	ClientID, found := settings.DecryptedSecureJSONData["ClientID"]
	if !found {
		return nil, fmt.Errorf("could not locate ClientID in NewDatasource")
	}
	log.Println(ClientID)

	TenantID, found := settings.DecryptedSecureJSONData["TenantID"]
	if !found {
		return nil, fmt.Errorf("could not locate TenantID in NewDatasource")
	}
	log.Println(TenantID)

	SubscriptionID, found := settings.DecryptedSecureJSONData["SubscriptionID"]
	if !found {
		return nil, fmt.Errorf("could not locate SubscriptionID in NewDatasource")
	}
	log.Println(SubscriptionID)

	Region, found := settings.DecryptedSecureJSONData["Region"]
	if !found {
		return nil, fmt.Errorf("could not locate Region in NewDatasource")
	}
	log.Println(Region)

	// Initialize the datasource with the configuration.
	config := Config{
		AzureCostSubscriptionUrl: "https://management.azure.com:443/subscriptions/",
		Password:                 Password,
		ClientID:                 ClientID,
		TenantID:                 TenantID,
		SubscriptionID:           SubscriptionID,
		Region:                   Region, //"UK South",
		TokenURL:                 "",
	}

	// Set up logging to a file
	logFile, err := os.Create("token_log.txt")
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	//
	opts, err := settings.HTTPClientOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("http client options: %w", err)
	}

	// Important: Reuse the same client for each query to avoid using all available connections on a host.
	opts.ForwardHTTPHeaders = true

	// Customize URL or add query parameters
	cl, err := httpclient.New(opts)
	if err != nil {
		return nil, fmt.Errorf("httpclient new: %w", err)
	}
	//

	log.Println("Datasource initialized")

	return &Datasource{config: config, httpClient: cl}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
// type Datasource struct{}
type Datasource struct {
	httpClient *http.Client
	config     Config
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	QueryText string `json:"queryText"`
	Constant  float64  `json:"constant"`
}

func (d *Datasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/developers/plugin-tools/introduction/data-frames
	frame := data.NewFrame("response")

	//MAB Added
	log.Print("Starting Datasource")
	log.Println("URL:", d.config.AzureCostSubscriptionUrl)
	log.Println("ResourceId:", qm.QueryText)

	// Call the fetchToken function
	token, err := fetchToken(d.config)
	if err != nil {
		log.Println("Error fetching token:", err)
		return response
	}

	// Use the token as needed
	log.Println("Fetched token:", token)

	start, end := getCurrentYearDates()

	timeRange := query.TimeRange
	if !timeRange.From.IsZero() && !timeRange.To.IsZero() {
		start = timeRange.From.Format("2006-01-02")
		end = timeRange.To.Format("2006-01-02")
	}

	costs, err := getCosts(token, d.config, start, end, qm.QueryText)
	if err != nil {
		log.Println("Error getting costs:", err)
		return response
	}

	//log.Printf("Costs: %+v", costs)
	datepoints, err := convertCostsToDatePoint(costs)
	if err != nil {
		log.Println("Error getting costs:", err)
		return response
	}

	log.Printf("Datapoint count: %d", len(datepoints))

	// // Add fields for "time" and "values"
	timeField := data.NewField("time", nil, make([]time.Time, len(datepoints)))
	valuesField := data.NewField("values", nil, make([]float64, len(datepoints)))
	sumField := data.NewField("sum-values", nil, make([]float64, len(datepoints)))

	rollingTotal := 0.0

	// Populate the fields with DatePoint values
	for i, dp := range datepoints {

		date, err := time.Parse("2006-01-02", dp.Date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		rollingTotal = rollingTotal + dp.Value

		timeField.Set(i, date)
		valuesField.Set(i, dp.Value)
		sumField.Set(i, rollingTotal)
	}

	// Add fields to the frame
	frame.Fields = append(frame.Fields, timeField, valuesField, sumField)

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Data source is working"

	if rand.Int()%2 == 0 {
		status = backend.HealthStatusError
		message = "health check randomized error"
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

// Api calls
func fetchToken(config Config) (string, error) {

	log.Print("Fetching Token")

	// URL for the token endpoint
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", config.TenantID)
	if len(config.TokenURL) > 1 {
		tokenURL = config.TokenURL
	}

	// Build the request parameters
	requestParams := url.Values{}
	requestParams.Set("grant_type", "client_credentials")
	requestParams.Set("client_id", config.ClientID)
	requestParams.Set("client_secret", config.Password)
	requestParams.Set("scope", "https://management.azure.com/.default")

	// Make a POST request to the token endpoint
	resp, err := http.PostForm(tokenURL, requestParams)
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("Too many requests: %v", resp.Status)
	}

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status: %v", resp.Status)
		return "", fmt.Errorf("unexpected response status: %v", resp.Status)
	}

	log.Print("Token Status is OK")

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the JSON response to get the access token
	token, err := parseAccessToken(body)
	if err != nil {
		return "", err
	}

	// Print the token
	log.Printf("Token is: %s", token)

	// Return the token
	return token, nil
}

// parseAccessToken parses the JSON response to get the access token
func parseAccessToken(body []byte) (string, error) {
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON response: %v", err)
	}

	accessToken, ok := response["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found in JSON response")
	}

	return accessToken, nil
}

// Fetch Costs
func getCosts(token string, config Config, start string, end string, resourceid string) (CostResponse, error) {
	url := config.SubscriptionID + "/providers/Microsoft.CostManagement/query?api-version=2023-03-01"

	if len(resourceid) > 2 {
		url = config.SubscriptionID + "/resourceGroups/" + resourceid + "/providers/Microsoft.CostManagement/query?api-version=2023-03-01"
	}

	log.Println("CostUrl:", url)

	bodyParameters := map[string]interface{}{
		"type":      "Usage",
		"timeframe": "Custom",
		"timeperiod": map[string]string{
			"from": start,
			"to":   end,
		},
		"dataset": map[string]interface{}{
			"granularity": "Daily",
			"aggregation": map[string]interface{}{
				"PreTaxCost": map[string]interface{}{
					"name":     "PreTaxCost",
					"function": "Sum",
				},
			},
		},
	}

	body, err := json.Marshal(bodyParameters)
	if err != nil {
		return CostResponse{}, fmt.Errorf("Error marshalling request body: %v", err)
	}

	requestURL := config.AzureCostSubscriptionUrl + url
	if len(config.TokenURL) > 1 {
		requestURL = config.TokenURL
	}
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(body))
	if err != nil {
		return CostResponse{}, fmt.Errorf("Error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CostResponse{}, fmt.Errorf("Error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CostResponse{}, fmt.Errorf("Error fetching cost. Status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return CostResponse{}, fmt.Errorf("Error reading response body: %v", err)
	}

	var costResponse CostResponse
	err = json.Unmarshal(responseBody, &costResponse)
	if err != nil {
		return CostResponse{}, fmt.Errorf("Error unmarshalling response body: %v", err)
	}

	return costResponse, nil
}

// FetchCosts fetches Azure costs and returns an array of DatePoint
func convertCostsToDatePoint(costs CostResponse) ([]DatePoint, error) {
	// Call the fetchToken method and handle the result

	// Convert result to DatePoint array
	datePoints := make([]DatePoint, 0)
	for _, row := range costs.Properties.Rows {
		// Assuming the value is always a float64 and date is in the second column
		value, ok := row[0].(float64)
		if !ok {
			return nil, fmt.Errorf("Invalid value format")
		}
		dateStr, ok := row[1].(float64)
		if !ok {
			return nil, fmt.Errorf("Invalid date format")
		}
		dateResult, err := convertToStandardDateFormat(strconv.Itoa(int(dateStr)))
		if err != nil {
			return nil, err
		}

		newDatePoint := DatePoint{
			Date:  dateResult,
			Value: value,
		}

		datePoints = append(datePoints, newDatePoint)
	}

	return datePoints, nil
}

// ConvertToStandardDateFormat converts the input date string to standard format
func convertToStandardDateFormat(inputDate string) (string, error) {
	if len(inputDate) == 8 {
		year := inputDate[:4]
		month := inputDate[4:6]
		day := inputDate[6:8]
		return fmt.Sprintf("%s-%s-%s", year, month, day), nil
	}

	// Handle invalid input length
	fmt.Println("Invalid date format. Expected YYYYMMDD.")
	return inputDate, fmt.Errorf("Invalid date format. Expected YYYYMMDD.")
}

func getCurrentYearDates() (string, string) {
	// Get the current year
	currentYear := time.Now().Year()

	// Get the 1st of January of the current year
	firstOfJanuary := time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	firstOfJanuaryFormatted := firstOfJanuary.Format("2006-01-02")

	// Get the 31st of December of the current year
	thirtyFirstOfDecember := time.Date(currentYear, time.December, 31, 0, 0, 0, 0, time.UTC)
	thirtyFirstOfDecemberFormatted := thirtyFirstOfDecember.Format("2006-01-02")

	return firstOfJanuaryFormatted, thirtyFirstOfDecemberFormatted
}