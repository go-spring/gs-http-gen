package jsonutil

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// JSONv2Decoder ...
type JSONv2Decoder struct {
	*jsontext.Decoder
}

// NewJSONv2Decoder ...
func NewJSONv2Decoder(d *jsontext.Decoder) Decoder {
	return &JSONv2Decoder{d}
}

// Unmarshal ...
func (d *JSONv2Decoder) Unmarshal(b []byte, i any) error {
	return json.Unmarshal(b, i)
}

// toKind ...
func toKind(k jsontext.Kind) Kind {
	switch k {
	case 'n':
		return 'n'
	case 'f':
		return 'f'
	case 't':
		return 't'
	case '"':
		return '"'
	case '0':
		return '0'
	case '{':
		return '{'
	case '}':
		return '}'
	case '[':
		return '['
	case ']':
		return ']'
	default:
		return InvalidKind
	}
}

// PeekKind ...
func (d *JSONv2Decoder) PeekKind() Kind {
	return toKind(d.Decoder.PeekKind())
}

// ReadToken ...
func (d *JSONv2Decoder) ReadToken() (string, Kind, error) {
	token, err := d.Decoder.ReadToken()
	if err != nil {
		return "", 0, err
	}
	return token.String(), toKind(token.Kind()), nil
}

// ReadValue ...
func (d *JSONv2Decoder) ReadValue() ([]byte, error) {
	return d.Decoder.ReadValue()
}

// SkipValue ...
func (d *JSONv2Decoder) SkipValue() error {
	return d.Decoder.SkipValue()
}
