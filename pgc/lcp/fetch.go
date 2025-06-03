package lcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lightweight-cache-proxy-service/internal/apis/secrets"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	Token      string
	httpClient http.Client
}

type Response[T CacheData] struct {
	Data    T
	Updated time.Time
}

func FetchCache[T CacheData](client *Client) (Response[T], error) {
	var zero Response[T]

	if client.Token == "" {
		return zero, errors.New("no token provided in client")
	}

	var cacheName string

	switch any(*new(T)).(type) {
	case []GitHubRepository:
		cacheName = "github"
	default:
		return zero, fmt.Errorf("unsupported cache type: %T", *new(T))
	}

	url, err := url.JoinPath(secrets.ENV.CorePath, cacheName)
	if err != nil {
		return zero, fmt.Errorf("failed to join path for URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return zero, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+client.Token)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return zero, fmt.Errorf("%s — response body: %q", msg, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, fmt.Errorf("failed to read response body: %w", err)
	}

	var result Response[T]
	if err := json.Unmarshal(body, &result); err != nil {
		return zero, fmt.Errorf("failed to unmarshal response JSON: %w", err)
	}

	return result, nil
}
