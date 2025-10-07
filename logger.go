package httpx

import "net/http"

type RequestLogger interface {
	Start(req *http.Request) any
	End(ctx any, resp *http.Response, err error)
}
