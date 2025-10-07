package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *Client) Post(path string, queryParam map[string]string, body any, v any) error {
	return c.request(http.MethodPost, c.url+path, queryParam, body, v)
}

func (c *Client) Get(path string, queryParams map[string]string, v any) error {
	return c.request(http.MethodGet, c.url+path, queryParams, nil, v)
}

func (c *Client) Patch(path string, queryParams map[string]string, body any, v any) error {
	return c.request(http.MethodPatch, c.url+path, queryParams, body, v)
}

func (c *Client) Put(path string) error {
	return c.request(http.MethodPut, c.url+path, nil, nil, nil)
}

func (c *Client) Delete(url string) error {
	return c.request(http.MethodDelete, url, nil, nil, nil)
}

func (c *Client) request(method string, url string, queryParams map[string]string, body any, v any) error {

	var b bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&b).Encode(body); err != nil {
			return fmt.Errorf("encoding: error: %w", err)
		}
	}

	r, err := http.NewRequest(method, url, &b)
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/json; charset=utf-8")
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
	if c.logger != nil {
		ctx = c.logger.Start(r)
	}

	response, err := c.http.Do(r)
	if c.logger != nil {
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
