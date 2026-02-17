package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	maxRetries     = 3
	initialBackoff = 1 * time.Second
)

type Client struct {
	BaseURL    string
	APIToken   string
	HTTPClient *http.Client
	UserAgent  string
}

type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewClient(baseURL, apiToken string) *Client {
	return &Client{
		BaseURL:  baseURL,
		APIToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		UserAgent: "terraform-provider-localskills",
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	url := c.BaseURL + path

	var resp *http.Response
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := initialBackoff * (1 << (attempt - 1))
			tflog.Debug(ctx, "retrying request", map[string]interface{}{
				"attempt": attempt,
				"backoff": backoff.String(),
				"method":  method,
				"path":    path,
			})
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}

			// Reset body reader for retry
			if body != nil {
				jsonBytes, err := json.Marshal(body)
				if err != nil {
					return nil, fmt.Errorf("marshaling request body: %w", err)
				}
				reqBody = bytes.NewReader(jsonBytes)
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.APIToken)
		req.Header.Set("User-Agent", c.UserAgent)

		resp, lastErr = c.HTTPClient.Do(req)
		if lastErr != nil {
			continue
		}

		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			if attempt < maxRetries {
				resp.Body.Close()
				lastErr = fmt.Errorf("retryable status code: %d", resp.StatusCode)
				continue
			}
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

func DoJSON[T any](c *Client, ctx context.Context, method, path string, body interface{}) (*T, error) {
	resp, err := c.doRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiResp ApiResponse[json.RawMessage]
		msg := string(respBody)
		if json.Unmarshal(respBody, &apiResp) == nil && apiResp.Error != "" {
			msg = apiResp.Error
		}
		return nil, &ApiError{
			StatusCode: resp.StatusCode,
			Message:    msg,
		}
	}

	var apiResp ApiResponse[T]
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if !apiResp.Success {
		return nil, &ApiError{
			StatusCode: resp.StatusCode,
			Message:    apiResp.Error,
		}
	}

	return &apiResp.Data, nil
}
