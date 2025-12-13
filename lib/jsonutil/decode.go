package jsonutil

import (
	"encoding/base64"
	"encoding/json/jsontext"
	"errors"
	"strconv"

	"github.com/lvan100/golib/errutil"
)

var ErrNull = errors.New("null")

// DecodeBool ...
func DecodeBool(d *jsontext.Decoder) (bool, error) {
	value, err := d.ReadValue()
	if err != nil {
		return false, err
	}
	switch value.Kind() {
	case 't':
		return true, nil
	case 'f':
		return false, nil
	case 'n':
		return false, ErrNull
	case '"':
		var s string
		s, err = strconv.Unquote(value.String())
		if err != nil {
			return false, err
		}
		return strconv.ParseBool(s)
	default:
		return false, errutil.Explain(err, "invalid JSON: expected boolean")
	}
}

// DecodeBoolPtr ...
func DecodeBoolPtr(d *jsontext.Decoder) (*bool, error) {
	b, err := DecodeBool(d)
	if err != nil {
		if errors.Is(err, ErrNull) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// DecodeInt ...
func DecodeInt[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) (T, error) {
	value, err := d.ReadValue()
	if err != nil {
		return 0, err
	}
	switch value.Kind() {
	case '0':
		i, err := strconv.ParseInt(value.String(), 10, 64)
		if err != nil {
			return 0, err
		}
		return T(i), nil
	case 'n':
		return 0, ErrNull
	case '"':
		var s string
		s, err = strconv.Unquote(value.String())
		if err != nil {
			return 0, err
		}
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(i), nil
	default:
		return 0, errutil.Explain(err, "invalid JSON: expected integer")
	}
}

// DecodeIntPtr ...
func DecodeIntPtr[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) (*T, error) {
	i, err := DecodeInt[T](d)
	if err != nil {
		if errors.Is(err, ErrNull) {
			return nil, nil
		}
		return nil, err
	}
	return &i, nil
}

// DecodeUint ...
func DecodeUint[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) (T, error) {
	value, err := d.ReadValue()
	if err != nil {
		return 0, err
	}
	switch value.Kind() {
	case '0':
		u, err := strconv.ParseUint(value.String(), 10, 64)
		if err != nil {
			return 0, err
		}
		return T(u), nil
	case 'n':
		return 0, ErrNull
	case '"':
		var s string
		s, err = strconv.Unquote(value.String())
		if err != nil {
			return 0, err
		}
		u, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(u), nil
	default:
		return 0, errutil.Explain(err, "invalid JSON: expected integer")
	}
}

// DecodeUintPtr ...
func DecodeUintPtr[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) (*T, error) {
	u, err := DecodeUint[T](d)
	if err != nil {
		if errors.Is(err, ErrNull) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// DecodeFloat ...
func DecodeFloat[T float32 | float64](d *jsontext.Decoder) (T, error) {
	value, err := d.ReadValue()
	if err != nil {
		return 0, err
	}
	switch value.Kind() {
	case '0':
		f, err := strconv.ParseFloat(value.String(), 64)
		if err != nil {
			return 0, err
		}
		return T(f), nil
	case 'n':
		return 0, ErrNull
	case '"':
		var s string
		s, err = strconv.Unquote(value.String())
		if err != nil {
			return 0, err
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return T(f), nil
	default:
		return 0, errutil.Explain(err, "invalid JSON: expected float")
	}
}

// DecodeFloatPtr ...
func DecodeFloatPtr[T float32 | float64](d *jsontext.Decoder) (*T, error) {
	f, err := DecodeFloat[T](d)
	if err != nil {
		if errors.Is(err, ErrNull) {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

// DecodeString decodes a string value from the given JSON decoder.
func DecodeString(d *jsontext.Decoder) (string, error) {
	value, err := d.ReadValue()
	if err != nil {
		return "", err
	}
	switch value.Kind() {
	case 'n':
		return "", ErrNull
	case '"':
		return strconv.Unquote(value.String())
	default:
		return "", errutil.Explain(err, "invalid JSON: expected string")
	}
}

// DecodeStringPtr ...
func DecodeStringPtr(d *jsontext.Decoder) (*string, error) {
	s, err := DecodeString(d)
	if err != nil {
		if errors.Is(err, ErrNull) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// DecodeBytes decodes a byte array value from the given JSON decoder.
func DecodeBytes(d *jsontext.Decoder) ([]byte, error) {
	value, err := d.ReadValue()
	if err != nil {
		return nil, err
	}
	switch value.Kind() {
	case 'n':
		return nil, nil
	case '"':
		s, err := strconv.Unquote(value.String())
		if err != nil {
			return nil, err
		}
		return base64.StdEncoding.DecodeString(s)
	default:
		return nil, errutil.Explain(err, "invalid JSON: expected []byte")
	}
}

// DecodeKey ...
func DecodeKey(d *jsontext.Decoder) (string, error) {
	key, err := DecodeString(d)
	if err != nil {
		if errors.Is(err, ErrNull) {
			return "", errutil.Explain(err, "invalid JSON: expected key")
		}
		return "", err
	}
	return key, nil
}

// DecodeObjectBegin ...
func DecodeObjectBegin(d *jsontext.Decoder) error {
	token, err := d.ReadToken()
	if err != nil {
		return err
	}
	if token.Kind() != '{' {
		return errutil.Explain(err, "invalid JSON: expected object")
	}
	return nil
}

// DecodeObjectEnd ...
func DecodeObjectEnd(d *jsontext.Decoder) error {
	token, err := d.ReadToken()
	if err != nil {
		return err
	}
	if token.Kind() != '}' {
		return errutil.Explain(err, "invalid JSON: expected end of object")
	}
	return nil
}

// DecodeArrayBegin ...
func DecodeArrayBegin(d *jsontext.Decoder) error {
	token, err := d.ReadToken()
	if err != nil {
		return err
	}
	if token.Kind() != '[' {
		return errutil.Explain(err, "invalid JSON: expected array")
	}
	return nil
}

// DecodeArrayEnd ...
func DecodeArrayEnd(d *jsontext.Decoder) error {
	token, err := d.ReadToken()
	if err != nil {
		return err
	}
	if token.Kind() != ']' {
		return errutil.Explain(err, "invalid JSON: expected end of array")
	}
	return nil
}
