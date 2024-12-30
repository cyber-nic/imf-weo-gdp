package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func makeIMFUrl(path string) string {
	// Define the URL for the POST request
	return fmt.Sprintf("https://betadata.imf.org/%s", path)
}

func main() {
	// Step 1: Submit the POST request and get the UUID
	uuid, err := submitQuery()
	if err != nil {
		log.Fatalf("Error fetching UUID: %v", err)
	}
	fmt.Printf("Received UUID: %s\n", uuid)

	body, err := fetchData(uuid)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	defer body.Close()

	filename := "data.csv"
	if err := saveToFile(body, filename); err != nil {
		log.Fatalf("Error saving data to file: %v", err)
	}

	fmt.Printf("Data saved to %s\n", filename)
}

// submitQuery requests the World Economic Outlook (WEO) data from the IMF API and returns the UUID for the async request.
func submitQuery() (string, error) {
	// Define the URL for the POST request
	url := makeIMFUrl("platform/rest/v2/engine/data/sync/submit")

	// Calculate time range
	now := time.Now()
	pastTenYears := now.AddDate(-10, 0, 0).UnixMilli()
	futureFiveYears := now.AddDate(5, 0, 0).UnixMilli()

	// Define the JSON payload
	payload := []byte(fmt.Sprintf(`{
		"agencyID": "IMF.RES",
		"attributes": "all",
		"filters": [
			{"componentCode": "TIME_PERIOD", "operator": "ge", "value": "%d"},
			{"componentCode": "TIME_PERIOD", "operator": "lt", "value": "%d"},
			{"componentCode": "INDICATOR", "operator": "eq", "value": "NGDP_D"}
		],
		"headerConfig": {"languages": ["en"]},
		"includeHistory": "false",
		"messageVersion": "2.0.0",
		"outputFormat": "CSV",
		"resourceID": "WEO",
		"startsWithBom": true,
		"version": "4.0.0",
		"_type": "SdmxTableDataQueryPlusV2",
		"columns": [
			{"componentId": "DATASET"},
			{"componentId": "SERIES_CODE"},
			{"componentId": "OBS_MEASURE"},
			{"componentId": "COUNTRY", "value": "name", "header": ".Excel"},
			{"componentId": "INDICATOR", "value": "name", "header": ".Excel"},
			{"componentId": "FREQUENCY", "value": "name", "header": ".Excel"}
		],
		"applyFormatting": true,
		"rows": [{"componentId": "OBS_VALUE"}],
		"viewMode": "TIMESERIES_PER_ROW"
	}`, pastTenYears, futureFiveYears))

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set the headers
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-dissemination-channel", "Portals")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read and validate the response body as a UUID
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Convert body to string and trim any whitespace
	uuidStr := string(body)
	uuidStr = strings.TrimSpace(uuidStr)

	// Validate the UUID using the google/uuid library
	if _, err := uuid.Parse(uuidStr); err != nil {
		return "", fmt.Errorf("invalid UUID in response: %w", err)
	}

	return uuidStr, nil
}

// fetchData performs a GET request using the UUID and returns the response body as a ReadCloser.
func fetchData(uuid string) (io.ReadCloser, error) {
	// Define the URL for the GET request
	url := makeIMFUrl(fmt.Sprintf("api/platform/v2/engine/data/sync/ott/%s", uuid))

	// Create a new HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	return resp.Body, nil
}

// saveToFile writes data from an io.Reader to a specified file.
func saveToFile(reader io.Reader, filename string) error {
	// Open the output file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Write the data to the file
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
