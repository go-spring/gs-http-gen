package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type DefaultHTTPClient struct {
	*http.Client
}

func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

type ExampleClient struct{}

type UrlPath interface {
	PathParams() map[string]string
}
type PingReq struct{}

func (r *PingReq) PathParams() []string {
	return []string{}
}

type PingResp struct{}

func (c *ExampleClient) Ping(ctx context.Context, req *PingReq) (http.Response, *PingResp, error) {
	getQueryString(req)
}

func buildURL(uri string, req UrlPath) (string, error) {
	if !strings.Contains(uri, "{") && !strings.Contains(uri, "*") {
		return uri, nil
	}

	// 获取路径参数
	pathParams := req.PathParams()

	// 替换路径中的参数，格式为 {param}
	result := uri
	for key, value := range pathParams {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// 处理通配符参数 *
	if strings.Contains(result, "*") && len(params)%2 == 0 && len(params) > 0 {
		// 通配符参数应该是最后一个参数
		wildcardValue := params[len(params)-1]
		result = strings.ReplaceAll(result, "*", wildcardValue)
	}

	return result, nil
}

func getQueryString(req any) (string, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	var m map[string]json.RawMessage
	if err = json.Unmarshal(b, &m); err != nil {
		return "", err
	}
	var buf strings.Builder
	for k, v := range m {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		if v[0] == '"' {
			s, err := strconv.Unquote(string(v))
			if err != nil {
				return "", err
			}
			buf.Write(s)
		} else {
			buf.Write(v)
		}
	}
	return buf.String(), nil
}

type Request struct {
	Req  *http.Request
	body io.Reader
	Err  error
}

// GetForm returns a GET request with form data.
func GetForm(ctx context.Context, url string, data url.Values) *Request {
	return form(ctx, "GET", url, data)
}

// PostForm returns a POST request with form data.
func PostForm(ctx context.Context, url string, data url.Values) *Request {
	return form(ctx, "POST", url, data)
}

func form(ctx context.Context, method string, url string, data url.Values) *Request {
	var body io.Reader
	if data != nil {
		body = strings.NewReader(data.Encode())
	}
	r := newRequest(ctx, method, url, body)
	r.Header("Content-Type", "application/x-www-form-urlencoded")
	return r
}

// PostJSON returns a POST request with JSON data.
func PostJSON(ctx context.Context, url string, data any) *Request {
	var body io.Reader
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return &Request{Err: err}
		}
		body = bytes.NewReader(b)
	}
	r := newRequest(ctx, "POST", url, body)
	r.Header("Content-Type", "application/json")
	return r
}

func GetFile(ctx context.Context, method string, url string) *Request {
	panic("not implemented")
}

func PostFile(ctx context.Context, url string, file io.Reader) *Request {
	panic("not implemented")
}

func newRequest(ctx context.Context, method, url string, body io.Reader) *Request {
	r, err := http.NewRequestWithContext(ctx, method, url, body)
	return &Request{Req: r, Err: err}
}

func (r *Request) Header(key, value string) *Request {
	r.Req.Header.Set(key, value)
	return r
}

func (r *Request) Cookie(key, value string) *Request {
	r.Req.AddCookie(&http.Cookie{Name: key, Value: value})
	return r
}

func (r *Request) Send(ctx context.Context, client Client) (*http.Response, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	_, err := http.NewRequestWithContext(ctx, r.Req.Method, r.Req.URL.String(), r.body)
	if err != nil {
		return nil, err
	}
	return client.Do(r.Req)
}
