package httpx

import (
	"net/http"
)

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

/**
per_request

Sometimes, even requests from the same system can have different behavior.For example:
like a client added logger, but some requests still don't need log.
*/

type requestOption func(ctx *requestContext)

func ReqNoLog() requestOption {
	return func(rc *requestContext) {
		rc.disableLog = true
	}
}

func ReqHeader(key, value string) requestOption {
	return func(rc *requestContext) {
		if rc.headers == nil {
			rc.headers = make(map[string]string)
		}
		rc.headers[key] = value
	}
}

func ReqModifier(fn func(*http.Request) error) requestOption {
	return func(rc *requestContext) {
		rc.modifier = fn
	}
}
