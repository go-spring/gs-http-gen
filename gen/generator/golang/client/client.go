package client

import (
	"context"
	"io"
	"net/http"
)

type Client struct {
}

func NewClient(m map[string]string) (*Client, error) {
	return nil, nil
}

type Request struct {
	Req  *http.Request
	body io.Reader
	Err  error
}

func NewRequest(method string, url string) *Request {
	r, err := http.NewRequest(method, url, nil)
	return &Request{Req: r, Err: err}
}

func (r *Request) Form(data map[string]string) *Request {
	r.Req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

func (r *Request) JSON(data any) *Request {
	r.Req.Header.Set("Content-Type", "application/json")
	return r
}

func (r *Request) File() *Request {
	panic("not implemented")
}

func (r *Request) Header(key, value string) *Request {
	r.Req.Header.Set(key, value)
	return r
}

func (r *Request) Cookie(key, value string) *Request {
	r.Req.AddCookie(&http.Cookie{Name: key, Value: value})
	return r
}

func (r *Request) Send(ctx context.Context) (*http.Response, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	_, err := http.NewRequestWithContext(ctx, r.Req.Method, r.Req.URL.String(), r.body)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// HEAD ...
func (c *Client) HEAD(url string) *Request {
	return NewRequest("HEAD", url)
}

// GET ...
func (c *Client) GET(url string) *Request {
	return NewRequest("GET", url)
}

// POST ...
func (c *Client) POST(url string) *Request {
	return NewRequest("POST", url)
}

// PUT ...
func (c *Client) PUT(url string) *Request {
	return NewRequest("PUT", url)
}

// PATCH ...
func (c *Client) PATCH(url string) *Request {
	return NewRequest("PATCH", url)
}

// DELETE ...
func (c *Client) DELETE(url string) *Request {
	return NewRequest("DELETE", url)
}
