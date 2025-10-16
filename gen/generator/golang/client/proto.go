package client

import (
	"encoding/json"
)

var (
	Url  Protocol = &UrlProtocol{}
	Json Protocol = &JsonProtocol{}
)

// Protocol 协议接口
type Protocol interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}

// UrlProtocol 默认的 url 编码协议
type UrlProtocol struct{}

func (p *UrlProtocol) Encode(v any) ([]byte, error) {
	panic("not implemented")
}

func (p *UrlProtocol) Decode(data []byte, v any) error {
	panic("not implemented")
}

// JsonProtocol json 编码协议
type JsonProtocol struct{}

func (p *JsonProtocol) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (p *JsonProtocol) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
