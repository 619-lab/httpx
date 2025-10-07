package httpx

import "net/http"

func WithHTTPClient(h *http.Client) func(*Client) {
	return func(c *Client) {
		c.http = h
	}
}

func WithLogger(l RequestLogger) func(*Client) {
	return func(c *Client) {
		c.logger = l
	}
}
