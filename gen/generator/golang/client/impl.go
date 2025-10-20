package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-spring/gs-http-gen/gen/generator/golang/httputil"
)

// ClientImpl 实现最终请求的 client 封装。
type ClientImpl struct {
	Client httputil.HTTPClient
}

func (c *ClientImpl) Ping(ctx context.Context, req *PingReq) (*http.Response, *PingResp, error) {
	urlPath := fmt.Sprintf("/v1/ping?name=%s", req.Name)
	r, err := httputil.NewRequest(ctx, "POST", urlPath, httputil.Form, &req.PingReqBody)
	if err != nil {
		return nil, nil, err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("XXX-Header", "123")
	return httputil.JsonResponse[PingResp](c.Client, r)
}
