package cli

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func PostRequest(address string, route string, contentType string, _ string, timeout int64) (*http.Response, error) {
	endpoint, err := url.JoinPath(address, route)
	if err != nil {
		return nil, fmt.Errorf("Error joining endpoint path: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string("")))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %w", err)
	}

	return resp, nil
}
