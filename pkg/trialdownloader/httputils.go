package trialdownloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getOne(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	return doOne(ctx, client, http.MethodGet, url, nil)
}

func postOne(ctx context.Context, client *http.Client, url string, body string) ([]byte, error) {
	return doOne(ctx, client, http.MethodPost, url, strings.NewReader(body))
}

func doOne(ctx context.Context, client *http.Client, method string, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("fetch page: building request: %w", err)
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	if method == http.MethodPost {
		req.Header.Add("content-type", "application/x-www-form-urlencoded")
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("fetch page: request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch page: unexpected status: %v (%v)", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetch page body: %w", err)
	}

	return data, nil
}
