package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	Category string  `json:"category"`
	Score    float64 `json:"score"`
}

type MLClient struct {
	baseURL string
	http    *http.Client
}

func NewMLClient(baseURL string, timeout time.Duration) *MLClient {
	return &MLClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *MLClient) Classify(ctx context.Context, content string) (Result, error) {
	body, err := json.Marshal(map[string]string{"content": content})
	if err != nil {
		return Result{}, fmt.Errorf("failed to marshal json: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/classify", bytes.NewReader(body))
	if err != nil {
		return Result{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("ml request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("ml services returned status %d", resp.StatusCode)
	}

	var out Result
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return Result{}, fmt.Errorf("failed to decode ml response: %w", err)
	}

	return out, nil
}
