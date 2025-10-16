package client

import (
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

func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, []byte, error) {
	panic("not implemented")
}
