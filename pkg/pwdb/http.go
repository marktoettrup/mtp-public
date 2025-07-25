package pwdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) newRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not parse endpoint URL %s: %w", c.endpoint, err)
	}
	u = u.JoinPath(endpoint)

	var reqBody io.Reader
	if body != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("could not encode body: %w", err)
		}
		reqBody = &buf
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("APIKey", c.apiKey)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) newRawRequest(ctx context.Context, method, endpoint string, query map[string]string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not parse endpoint URL %s: %w", c.endpoint, err)
	}
	u = u.JoinPath(endpoint)

	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("APIKey", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "multipart/form-data")
	}

	return req, nil
}

func (c *Client) do(req *http.Request, response interface{}) error {
	res, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, res.Body); err != nil {
		return fmt.Errorf("could not copy response to buffer: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code %d: %s", res.StatusCode, buf.String())
	}

	if err := json.NewDecoder(&buf).Decode(response); err != nil {
		return fmt.Errorf("could not decode response: %w: %s", err, buf.String())
	}

	return nil
}

func (c *Client) doRaw(req *http.Request) ([]byte, error) {
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, res.Body); err != nil {
		return nil, fmt.Errorf("could not copy response to buffer: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d: %s", res.StatusCode, buf.String())
	}

	return buf.Bytes(), nil
}
