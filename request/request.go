package request

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Error struct {
	Message any
	Status  int
}

func (e *Error) Error() string {
	switch msg := e.Message.(type) {
	case string:
		return msg
	case []byte:
		return string(msg)
	default:
		return fmt.Sprintf("%v", e.Message)
	}
}

type Request struct {
	body     io.Reader
	endpoint *url.URL
	headers  map[string]string
	queries  map[string]string
}

type Response struct {
	Body   []byte
	Status int
}

func New() *Request {
	return &Request{headers: make(map[string]string), queries: make(map[string]string)}
}

func (r *Request) SetEndpoint(endpoint *url.URL) *Request {
	r.endpoint = endpoint
	return r
}

func (r *Request) SetBody(body io.Reader) *Request {
	r.body = body
	return r
}

func (r *Request) AppendHeader(key string, value string) *Request {
	r.headers[key] = value
	return r
}

func (r *Request) AppendURLQuery(key string, value string) *Request {
	r.queries[key] = value
	return r
}

func (r *Request) Post(ctx context.Context) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.endpoint.String(), r.body)
	if err != nil {
		// nolint:wrapcheck
		return nil, err
	}

	for k, v := range r.headers {
		req.Header.Add(k, v)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		// nolint:wrapcheck
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// nolint:wrapcheck
		return nil, err
	}

	return &Response{Body: body, Status: resp.StatusCode}, nil
}

func (r *Request) Get(ctx context.Context) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.endpoint.String(), r.body)
	if err != nil {
		// nolint:wrapcheck
		return nil, err
	}

	for k, v := range r.headers {
		req.Header.Add(k, v)
	}

	q := req.URL.Query()
	for k, v := range r.queries {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		// nolint:wrapcheck
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// nolint:wrapcheck
		return nil, err
	}

	return &Response{Body: body, Status: resp.StatusCode}, nil
}
