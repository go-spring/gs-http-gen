package client

import (
	"encoding/json"
	"net/url"
	"strconv"
)

var (
	Form Protocol = &FormProtocol{}
	Json Protocol = &JsonProtocol{}
)

// Protocol defines a common interface for encoding and decoding data.
type Protocol interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}

// FormProtocol serializes data into application/x-www-form-urlencoded format.
type FormProtocol struct{}

// Encode converts a struct or map into a form-encoded byte slice.
func (p *FormProtocol) Encode(i any) ([]byte, error) {
	if i == nil {
		return []byte(""), nil
	}

	b, err := Json.Encode(i)
	if err != nil {
		return nil, err
	}

	var m map[string]json.RawMessage
	if err = Json.Decode(b, &m); err != nil {
		return nil, err
	}

	u := url.Values{}
	for k, v := range m {
		if v[0] == '"' {
			s, err := strconv.Unquote(string(v))
			if err != nil {
				return nil, err
			}
			u.Set(k, s)
		} else {
			u.Set(k, string(v))
		}
	}

	return []byte(u.Encode()), nil
}

// Decode parses form-encoded data into the given struct or map.
func (p *FormProtocol) Decode(data []byte, v any) error {
	panic("not implemented")
}

// JsonProtocol implements the JSON encoding protocol.
type JsonProtocol struct{}

// Encode serializes the input value into JSON format.
func (p *JsonProtocol) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Decode deserializes JSON data into the given value.
func (p *JsonProtocol) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
