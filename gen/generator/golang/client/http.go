package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// HTTPClient 用户实现此接口可以实现符合自己要求的 http 执行器
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, []byte, error)
}

// DefaultHTTPClient 默认的 http 执行器
type DefaultHTTPClient struct {
	*http.Client
}

// Do 执行 http 请求，返回 *http.Response 对象和返回值数据
func (c *DefaultHTTPClient) Do(r *http.Request) (*http.Response, []byte, error) {
	resp, err := c.Client.Do(r)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(b))
	return resp, b, nil
}

// ----------------------------------------------------------------------------

// NewRequest 创建 http 请求，需要基础的 ctx、method、urlPath、protocol、body 参数
func NewRequest(ctx context.Context, method string, urlPath string, p Protocol, body any) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		b, err := p.Encode(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(b)
	}
	return http.NewRequestWithContext(ctx, method, urlPath, reader)
}

// DoResponse 执行 http 请求，返回 *http.Response 、返回值数据
func DoResponse[RespType any](c HTTPClient, r *http.Request) (*http.Response, *RespType, error) {
	resp, b, err := c.Do(r)
	if err != nil {
		return nil, nil, err
	}
	var ret *RespType
	err = json.Unmarshal(b, &ret)
	if err != nil {
		return nil, nil, err
	}
	return resp, ret, nil
}
