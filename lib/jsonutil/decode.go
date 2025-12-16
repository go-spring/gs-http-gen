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

type Number interface {
	int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 |
	float32 | float64
}

func DecodeNumber[T Number](d *jsontext.Decoder, parseFn func(string) (T, error)) (T, error) {
	var v T
	value, err := d.ReadValue()
	if err != nil {
		return v, err
	}
	switch value.Kind() {
	case '0':
		return parseFn(value.String())
	case 'n':
		return v, ErrNull
	case '"':
		var s string
		s, err = strconv.Unquote(value.String())
		if err != nil {
			return v, err
		}
		return parseFn(s)
	default:
		return v, errutil.Explain(err, "invalid JSON: expected integer")
	}
}

// DecodeIntV2 ...
func DecodeIntV2[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) (T, error) {
	return DecodeNumber(d, func(s string) (T, error) {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
}

// DecodeInt ...
func DecodeInt[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) (T, error) {
	return DecodeIntV2[T](d)
	//value, err := d.ReadValue()
	//if err != nil {
	//	return 0, err
	//}
	//switch value.Kind() {
	//case '0':
	//	i, err := strconv.ParseInt(value.String(), 10, 64)
	//	if err != nil {
	//		return 0, err
	//	}
	//	return T(i), nil
	//case 'n':
	//	return 0, ErrNull
	//case '"':
	//	var s string
	//	s, err = strconv.Unquote(value.String())
	//	if err != nil {
	//		return 0, err
	//	}
	//	i, err := strconv.ParseInt(s, 10, 64)
	//	if err != nil {
	//		return 0, err
	//	}
	//	return T(i), nil
	//default:
	//	return 0, errutil.Explain(err, "invalid JSON: expected integer")
	//}
}

// DecodeUintV2 ...
func DecodeUintV2[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) (T, error) {
	return DecodeNumber(d, func(s string) (T, error) {
		u, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(u), nil
	})
}

// DecodeUint ...
func DecodeUint[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) (T, error) {
	return DecodeUintV2[T](d)
	//value, err := d.ReadValue()
	//if err != nil {
	//	return 0, err
	//}
	//switch value.Kind() {
	//case '0':
	//	u, err := strconv.ParseUint(value.String(), 10, 64)
	//	if err != nil {
	//		return 0, err
	//	}
	//	return T(u), nil
	//case 'n':
	//	return 0, ErrNull
	//case '"':
	//	var s string
	//	s, err = strconv.Unquote(value.String())
	//	if err != nil {
	//		return 0, err
	//	}
	//	u, err := strconv.ParseUint(s, 10, 64)
	//	if err != nil {
	//		return 0, err
	//	}
	//	return T(u), nil
	//default:
	//	return 0, errutil.Explain(err, "invalid JSON: expected integer")
	//}
}

// DecodeFloatV2 ...
func DecodeFloatV2[T float32 | float64](d *jsontext.Decoder) (T, error) {
	return DecodeNumber(d, func(s string) (T, error) {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return T(f), nil
	})
}

// DecodeFloat ...
func DecodeFloat[T float32 | float64](d *jsontext.Decoder) (T, error) {
	return DecodeFloatV2[T](d)
	//value, err := d.ReadValue()
	//if err != nil {
	//	return 0, err
	//}
	//switch value.Kind() {
	//case '0':
	//	f, err := strconv.ParseFloat(value.String(), 64)
	//	if err != nil {
	//		return 0, err
	//	}
	//	return T(f), nil
	//case 'n':
	//	return 0, ErrNull
	//case '"':
	//	var s string
	//	s, err = strconv.Unquote(value.String())
	//	if err != nil {
	//		return 0, err
	//	}
	//	f, err := strconv.ParseFloat(s, 64)
	//	if err != nil {
	//		return 0, err
	//	}
	//	return T(f), nil
	//default:
	//	return 0, errutil.Explain(err, "invalid JSON: expected float")
	//}
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

// DecodeBools ...
func DecodeBools(d *jsontext.Decoder) ([]bool, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []bool
	for {
		if d.PeekKind() == ']' {
			break
		}
		i, err := DecodeBool(d)
		if err != nil {
			if errors.Is(err, ErrNull) {
				return nil, errutil.Explain(nil, "null value is not allowed")
			}
			return nil, err
		}
		v = append(v, i)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeBoolPtrs ...
func DecodeBoolPtrs(d *jsontext.Decoder) ([]*bool, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []*bool
	for {
		if d.PeekKind() == ']' {
			break
		}
		i, err := DecodeBoolPtr(d)
		if err != nil {
			return nil, err
		}
		v = append(v, i)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeArray ...
func DecodeArray[T Number | bool | string](d *jsontext.Decoder, parseFn func(string) (T, error)) ([]T, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []T
	for {
		if d.PeekKind() == ']' {
			break
		}
		value, err := d.ReadValue()
		if err != nil {
			return v, err
		}
		i, err := parseFn(value.String())
		if err != nil {
			if errors.Is(err, ErrNull) {
				return nil, errutil.Explain(nil, "null value is not allowed")
			}
			return nil, err
		}
		v = append(v, i)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeIntsV2 ...
func DecodeIntsV2[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) ([]T, error) {
	return DecodeArray[T](d, func(s string) (T, error) {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
}

// DecodeInts ...
func DecodeInts[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) ([]T, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []T
	for {
		if d.PeekKind() == ']' {
			break
		}
		i, err := DecodeInt[T](d)
		if err != nil {
			if errors.Is(err, ErrNull) {
				return nil, errutil.Explain(nil, "null value is not allowed")
			}
			return nil, err
		}
		v = append(v, i)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeUintsV2 ...
func DecodeUintsV2[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) ([]T, error) {
	return DecodeArray[T](d, func(s string) (T, error) {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
}

// DecodeUints ...
func DecodeUints[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) ([]T, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []T
	for {
		if d.PeekKind() == ']' {
			break
		}
		i, err := DecodeUint[T](d)
		if err != nil {
			if errors.Is(err, ErrNull) {
				return nil, errutil.Explain(nil, "null value is not allowed")
			}
			return nil, err
		}
		v = append(v, i)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeFloatsV2 ...
func DecodeFloatsV2[T float32 | float64](d *jsontext.Decoder) ([]T, error) {
	return DecodeArray[T](d, func(s string) (T, error) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
}

// DecodeFloats ...
func DecodeFloats[T float32 | float64](d *jsontext.Decoder) ([]T, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []T
	for {
		if d.PeekKind() == ']' {
			break
		}
		i, err := DecodeFloat[T](d)
		if err != nil {
			if errors.Is(err, ErrNull) {
				return nil, errutil.Explain(nil, "null value is not allowed")
			}
		}
		v = append(v, i)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeStringsV2 ...
func DecodeStringsV2(d *jsontext.Decoder) ([]string, error) {
	return DecodeArray[string](d, func(s string) (string, error) {
		return s, nil
	})
}

// DecodeStrings ...
func DecodeStrings(d *jsontext.Decoder) ([]string, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []string
	for {
		if d.PeekKind() == ']' {
			break
		}
		s, err := DecodeString(d)
		if err != nil {
			if errors.Is(err, ErrNull) {
				return nil, errutil.Explain(nil, "null value is not allowed")
			}
			return nil, err
		}
		v = append(v, s)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeIntPtrs ...
func DecodeIntPtrs[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) ([]*T, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []*T
	for {
		if d.PeekKind() == ']' {
			break
		}
		i, err := DecodeIntPtr[T](d)
		if err != nil {
			return nil, err
		}
		v = append(v, i)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeUintPtrs ...
func DecodeUintPtrs[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) ([]*T, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []*T
	for {
		if d.PeekKind() == ']' {
			break
		}
		i, err := DecodeUintPtr[T](d)
		if err != nil {
			return nil, err
		}
		v = append(v, i)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeFloatPtrs ...
func DecodeFloatPtrs[T float32 | float64](d *jsontext.Decoder) ([]*T, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []*T
	for {
		if d.PeekKind() == ']' {
			break
		}
		i, err := DecodeFloatPtr[T](d)
		if err != nil {
			return nil, err
		}
		v = append(v, i)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeStringPtrs ...
func DecodeStringPtrs(d *jsontext.Decoder) ([]*string, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []*string
	for {
		if d.PeekKind() == ']' {
			break
		}
		s, err := DecodeStringPtr(d)
		if err != nil {
			return nil, err
		}
		v = append(v, s)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}
