package client

import (
	"context"
	"fmt"
	"net/http"
)

// ClientImpl 实现最终请求的 client 封装。
type ClientImpl struct {
	Client HTTPClient
}

func (c *ClientImpl) Ping(ctx context.Context, req *PingReq) (*http.Response, *PingResp, error) {
	urlPath := fmt.Sprintf("%v", req)
	r, err := NewRequest(ctx, "POST", urlPath, Url, &req.PingReqBody)
	if err != nil {
		return nil, nil, err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("XXX-XXX", fmt.Sprint(req.PingReqBody))
	return DoResponse[PingResp](c.Client, r)
}
