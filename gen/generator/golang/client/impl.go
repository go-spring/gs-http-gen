package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ClientImpl 实现最终请求的 client 封装。
type ClientImpl struct {
	Client HTTPClient
	Url    Protocol
	Json   Protocol
}

func NewClientImpl() *ClientImpl {
	return &ClientImpl{
		Client: &DefaultHTTPClient{},
		Url:    &UrlFlattenProtocol{},
		Json:   &JsonProtocol{},
	}
}

func (c *ClientImpl) Ping(ctx context.Context, req *PingReq) (*http.Response, *PingResp, error) {
	urlPath := fmt.Sprintf("%v", req)
	r, err := NewRequest(ctx, "POST", urlPath, c.Url, &req.PingReqBody)
	if err != nil {
		return nil, nil, err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return DoResponse[PingResp](c.Client, r)
}

func NewRequest(ctx context.Context, method string, urlPath string, Url Protocol, body any) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		b, err := Url.Encode(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(b)
	}
	r, err := http.NewRequestWithContext(ctx, method, urlPath, reader)
	if err != nil {
		return nil, err
	}
	return r, nil
}

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
