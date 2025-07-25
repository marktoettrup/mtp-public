package pwdb

import (
	"net/http"
)

type Client struct {
	endpoint string
	apiKey   string
	hc       *http.Client
}

type Option func(*Client)

func New(opts ...Option) *Client {
	c := &Client{
		hc: &http.Client{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithEndpoint(endpoint string) Option {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

func WithAPIKey(apiKey string) Option {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}
