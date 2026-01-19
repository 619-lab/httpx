package httpx_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/619-lab/httpx"
)

// go test -v -run=TestPPT
func TestPPT(t *testing.T) {
	client := httpx.NewClient("https://api-m.sandbox.paypal.com", httpx.WithLogger(printLogger{}))
	client.Trick = func(r *http.Request) {
		r.SetBasicAuth("CLIENT_ID",
			"CLIENT_SECRET")
	}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	resp := make(map[string]any)
	err := client.Post("/v1/oauth2/token", nil, data, &resp, httpx.ReqForm())
	if err != nil {
		fmt.Println("eerrr", err)
	}
	fmt.Println(resp)
}
