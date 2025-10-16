package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ClientImpl 实现最终请求的 client 封装。
type ClientImpl struct {
	HTTPClient
}

func (c *ClientImpl) Ping(ctx context.Context, req *PingReq) (*http.Response, *PingResp, error) {
	b, err := json.Marshal(req.PingReqBody)
	if err != nil {
		return nil, nil, err
	}
	urlPath := fmt.Sprintf("%v", req)
	r, err := http.NewRequestWithContext(ctx, "POST", urlPath, bytes.NewReader(b))
	if err != nil {
		return nil, nil, err
	}
	resp, b, err := c.HTTPClient.Do(r)
	if err != nil {
		return nil, nil, err
	}
	var ret *PingResp
	err = json.Unmarshal(b, &ret)
	if err != nil {
		return nil, nil, err
	}
	return resp, ret, nil
}
