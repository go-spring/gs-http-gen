package client

import (
	"encoding/json"
	"net/http"
	"net/url"
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

type Protocol interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}

type UrlFlattenProtocol struct{}

type FormValues interface {
	FormValues() url.Values
}

func (p *UrlFlattenProtocol) Encode(v any) ([]byte, error) {
	if v, ok := v.(FormValues); ok {
		return []byte(v.FormValues().Encode()), nil
	}
	return nil, nil
}

func (p *UrlFlattenProtocol) Decode(data []byte, v any) error {
	panic("not implemented")
}

type JsonProtocol struct{}

func (p *JsonProtocol) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (p *JsonProtocol) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
