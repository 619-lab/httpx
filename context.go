package httpx

import (
	"net/http"
)

type requestContext struct {
	headers    map[string]string
	disableLog bool
	modifier   func(*http.Request) error
	bodyType   bodyType

	maxResponseBytes int64
}

type bodyType int

const (
	bodyJSON bodyType = iota
	bodyForm
)
