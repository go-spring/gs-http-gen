package client

import (
	"context"
	"net/http"

	"github.com/lvan100/golib/httputil"
)

// ClientInterface 通过接口实现中间件机制
type ClientInterface interface {
	Ping(ctx context.Context, req *PingReq) (*http.Response, *PingResp, error)
}

// Client 通过结构体追踪构建过程
type Client struct {
	ClientInterface
}

// NewClient 使用默认方案创建客户端
func NewClient(config map[string]any) *Client {
	return &Client{&ClientImpl{
		Client: &httputil.DefaultClient{
			Client: http.DefaultClient,
		},
	}}
}
