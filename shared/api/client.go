package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

// A client for making API requests
type ApiClient struct {
	// The base address and route for API calls
	apiUrl *url.URL

	// An HTTP client for sending requests
	client *http.Client

	// Logger to print debug messages to
	logger *slog.Logger
}

// Create a new API client
func NewApiClient(logger *slog.Logger, serviceName string, port uint) (*ApiClient, error) {
	// Create the address
	address := fmt.Sprintf("http://%s:%d/%s", serviceName, port, ApiRoute)

	// Parse the address
	apiUrl, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("error parsing API address: %w", err)
	}

	// Create the client
	client := &http.Client{}

	return &ApiClient{apiUrl, client, logger}, nil
}

// Get a config parameter from the API
func (c *ApiClient) GetParam(paramID string) (*ApiResponse[GetParamResponse], error) {
	args := map[string]string{
		"id": paramID,
	}
	return sendGetRequest[GetParamResponse](c, "param", args)
}

// Set a config parameter via the API
func (c *ApiClient) SetParam(paramID string, value string) (*ApiResponse[SuccessResponse], error) {
	body := PostParamRequest{
		ID:    paramID,
		Value: value,
	}
	return sendPostRequest[SuccessResponse](c, "param", body)
}

// Submit a GET request to the API server
func sendGetRequest[DataType any](client *ApiClient, method string, args map[string]string) (*ApiResponse[DataType], error) {
	if args == nil {
		args = map[string]string{}
	}

	// Create the request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", client.apiUrl, method), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Encode the params into a query string
	values := url.Values{}
	for name, value := range args {
		values.Add(name, value)
	}
	req.URL.RawQuery = values.Encode()

	// Debug log
	client.logger.Debug("API Request",
		slog.String("method", http.MethodGet),
		slog.String("query", req.URL.String()),
	)

	// Run the request
	resp, err := client.client.Do(req)
	return handleResponse[DataType](client.logger, resp, method, err)
}

// Submit a POST request to the API server
func sendPostRequest[DataType any](client *ApiClient, method string, body any) (*ApiResponse[DataType], error) {
	// Serialize the body
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error serializing request body for POST %s: %w", method, err)
	}

	// Create the request
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", client.apiUrl, method), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Debug log
	client.logger.Debug("API Request",
		slog.String("method", http.MethodPost),
		slog.String("path", method),
		slog.String("body", string(bodyBytes)),
	)

	// Run the request
	resp, err := client.client.Do(req)
	return handleResponse[DataType](client.logger, resp, method, err)
}

// Processes a response to a request
func handleResponse[DataType any](logger *slog.Logger, resp *http.Response, path string, err error) (*ApiResponse[DataType], error) {
	if err != nil {
		return nil, fmt.Errorf("error requesting %s: %w", path, err)
	}

	// Read the body
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response body for %s: %w", path, err)
	}

	// Handle 404s specially since they won't have a JSON body
	if resp.StatusCode == http.StatusNotFound {
		logger.Debug("API Response (raw)",
			slog.String("code", resp.Status),
			slog.String("body", string(bytes)),
		)
		return nil, fmt.Errorf("route '%s' not found", path)
	}

	// Deserialize the response into the provided type
	var parsedResponse ApiResponse[DataType]
	err = json.Unmarshal(bytes, &parsedResponse)
	if err != nil {
		logger.Debug("API Response (raw)",
			slog.String("code", resp.Status),
			slog.String("body", string(bytes)),
		)
		return nil, fmt.Errorf("error deserializing response to %s: %w", path, err)
	}

	// Check if the request failed
	if resp.StatusCode != http.StatusOK {
		logger.Debug("API Response",
			slog.String("path", path),
			slog.String("code", resp.Status),
			slog.String("err", parsedResponse.Error),
		)
		return nil, errors.New(parsedResponse.Error)
	}

	// Debug log
	logger.Debug("API Response",
		slog.String("body", string(bytes)),
	)

	return &parsedResponse, nil
}
