package httpx_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/619-lab/httpx"
)

type printLogger struct{}

func (printLogger) Start(req *http.Request) any {
	fmt.Println("➡️ Sending request:", req.Method, req.URL.String())
	return nil
}

func (printLogger) End(ctx any, resp *http.Response, err error) {
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}
	fmt.Println("✅ Got response:", resp.Status)
}

func optHBO(c *httpx.Client) {
	c.Trick = func(r *http.Request) {
		r.Header.Set("HBO-Access-Key", "blah blah")
	}
}

func TestHttpx(t *testing.T) {
	client := httpx.NewClient("https://httpbingo.org", httpx.WithLogger(printLogger{}), optHBO)

	var result map[string]any

	if err := client.Get("/get", map[string]string{"q": "hello"}, &result, httpx.ReqNoLog()); err != nil {
		panic(err)
	}

	if err := client.Get("/get", map[string]string{"q": "world"}, &result); err != nil {
		panic(err)
	}

	fmt.Println("Response data:", result["url"])
}
