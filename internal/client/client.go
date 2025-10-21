package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	apiKeyHeader   = "X-TierZero-Org-Api-Key"
)

// Client is the HTTP client for the TierZero API
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new TierZero API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set(apiKeyHeader, c.APIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
		}
	}

	return respBody, nil
}

// APIError represents an error response from the API
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}
