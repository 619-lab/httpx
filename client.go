package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	url    string
	http   *http.Client
	logger RequestLogger
	Trick  func(r *http.Request)
}

func NewClient(url string, options ...func(*Client)) *Client {
	cln := Client{
		url:  url,
		http: defaultClient,
	}
	for _, option := range options {
		option(&cln)
	}
	return &cln
}

func (c *Client) Post(path string, queryParam map[string]string, body any, v any, reqOpts ...requestOption) error {
	return c.request(http.MethodPost, c.url+path, queryParam, body, v, reqOpts...)
}

func (c *Client) Get(path string, queryParams map[string]string, v any, reqOpts ...requestOption) error {
	return c.request(http.MethodGet, c.url+path, queryParams, nil, v, reqOpts...)
}

func (c *Client) Patch(path string, queryParams map[string]string, body any, v any, reqOpts ...requestOption) error {
	return c.request(http.MethodPatch, c.url+path, queryParams, body, v, reqOpts...)
}

func (c *Client) Put(path string, reqOpts ...requestOption) error {
	return c.request(http.MethodPut, c.url+path, nil, nil, nil, reqOpts...)
}

func (c *Client) Delete(url string, reqOpts ...requestOption) error {
	return c.request(http.MethodDelete, url, nil, nil, nil, reqOpts...)
}

func (c *Client) request(method string, url_ string, queryParams map[string]string, body any, v any, reqOpts ...requestOption) error {

	rc := &requestContext{}
	for _, o := range reqOpts {
		o(rc)
	}

	var b bytes.Buffer
	if body != nil {
		switch rc.bodyType {
		case bodyForm:
			switch fc := body.(type) {
			case string:
				b.WriteString(fc)
			case url.Values:
				b.WriteString(fc.Encode())
			default:
				return fmt.Errorf("form body must be string or url.Values")
			}
		default:
			if err := json.NewEncoder(&b).Encode(body); err != nil {
				return fmt.Errorf("json encode failed: %w", err)
			}
		}
	}

	r, err := http.NewRequest(method, url_, &b)
	if err != nil {
		return err
	}

	switch rc.bodyType {
	case bodyForm:
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	default:
		r.Header.Add("Content-Type", "application/json; charset=utf-8")
	}

	if c.Trick != nil {
		c.Trick(r)
	}

	q := r.URL.Query()
	for k, v := range queryParams {
		if k != "" && v != "" {
			q.Add(k, v)
		}
	}
	r.URL.RawQuery = q.Encode()

	var ctx any
	if !rc.disableLog && c.logger != nil {
		ctx = c.logger.Start(r)
	}

	// requestOption
	for hk, hv := range rc.headers {
		r.Header.Set(hk, hv)
	}

	if rc.modifier != nil {
		rc.modifier(r)
	}

	//==========================================================================

	response, err := c.http.Do(r)
	if !rc.disableLog && c.logger != nil {
		c.logger.End(ctx, response, err)
	}

	if err != nil {
		return fmt.Errorf("HTTPClient.Do: %w", err)
	}
	defer response.Body.Close()
	if v != nil {
		err = json.NewDecoder(response.Body).Decode(v)
		if err != nil {
			return fmt.Errorf("decode error: %w", err)
		}
	}
	return nil
}
