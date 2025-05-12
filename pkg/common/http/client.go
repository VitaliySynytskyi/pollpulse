package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a wrapper around the standard http.Client with additional features
type Client struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

// NewClient creates a new HTTP client
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		headers: make(map[string]string),
	}
}

// SetHeader sets a header for all requests
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// SetAuthToken sets the Authorization header with a Bearer token
func (c *Client) SetAuthToken(token string) {
	c.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
}

// Request makes an HTTP request with the given method, path, and payload
func (c *Client) Request(ctx context.Context, method, path string, payload interface{}, result interface{}) error {
	url := c.baseURL + path

	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type if we have a payload
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set any custom headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error status codes
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	// Unmarshal the response if a result container was provided
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.Request(ctx, http.MethodGet, path, nil, result)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, payload interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, path, payload, result)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, payload interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPut, path, payload, result)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, path, nil, result)
} 