package client

import (
	"context"
	"net/http"
)

// MockClientImpl 可以实现 mock 功能的 client 中间件
type MockClientImpl struct {
	ClientInterface
}

func (c *MockClientImpl) Ping(ctx context.Context, req *PingReq) (*http.Response, *PingResp, error) {
	// ...
	return c.ClientInterface.Ping(ctx, req)
}
