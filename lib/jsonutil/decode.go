package jsonutil

import (
	"encoding/base64"
	"encoding/json/jsontext"
	"errors"
	"strconv"

	"github.com/lvan100/golib/errutil"
)

var ErrNull = errors.New("null")

// Object ...
type Object interface {
	DecodeJSON(d *jsontext.Decoder) error
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

//////////////////////////////////// parse ////////////////////////////////////

// ParseBool ...
func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// ParseInt ...
func ParseInt[T int | int8 | int16 | int32 | int64](s string) (T, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return T(v), nil
}

// ParseUint ...
func ParseUint[T uint | uint8 | uint16 | uint32 | uint64](s string) (T, error) {
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return T(u), nil
}

// ParseFloat ...
func ParseFloat[T float32 | float64](s string) (T, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return T(f), nil
}

// ParseString ...
func ParseString(s string) (string, error) {
	return s, nil
}

// ParseBytes ...
func ParseBytes(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

//////////////////////////////////// value ////////////////////////////////////

// DecodeObject ...
func DecodeObject[T Object](f func() T) func(d *jsontext.Decoder) (T, error) {
	return func(d *jsontext.Decoder) (T, error) {
		var zero T
		switch d.PeekKind() {
		case 'n':
			return zero, nil
		case '{':
			v := f()
			if err := v.DecodeJSON(d); err != nil {
				return zero, err
			}
			return v, nil
		default:
			return zero, errutil.Explain(nil, "invalid JSON: expected object")
		}
	}
}

// DecodeValue ...
func DecodeValue[T any](parseFn func(string) (T, error)) func(d *jsontext.Decoder) (T, error) {
	return func(d *jsontext.Decoder) (T, error) {
		var zero T
		value, err := d.ReadValue()
		if err != nil {
			return zero, err
		}
		switch value.Kind() {
		case 'n':
			return zero, ErrNull
		case 'f', 't', '0':
			return parseFn(value.String())
		case '"':
			var s string
			s, err = strconv.Unquote(value.String())
			if err != nil {
				return zero, err
			}
			return parseFn(s)
		default:
			return zero, errutil.Explain(err, "invalid JSON: expected value")
		}
	}
}

// DecodeValuePtr ...
func DecodeValuePtr[T any](parseFn func(string) (T, error)) func(d *jsontext.Decoder) (*T, error) {
	return func(d *jsontext.Decoder) (*T, error) {
		v, err := DecodeValue(parseFn)(d)
		if err != nil {
			if errors.Is(err, ErrNull) {
				return nil, nil
			}
			return nil, err
		}
		return &v, nil
	}
}

var DecodeBool = DecodeValue(ParseBool)
var DecodeBoolPtr = DecodeValuePtr(ParseBool)

var DecodeInt = DecodeValue(ParseInt[int])
var DecodeInt8 = DecodeValue(ParseInt[int8])
var DecodeInt16 = DecodeValue(ParseInt[int16])
var DecodeInt32 = DecodeValue(ParseInt[int32])
var DecodeInt64 = DecodeValue(ParseInt[int64])

var DecodeIntPtr = DecodeValuePtr(ParseInt[int])
var DecodeInt8Ptr = DecodeValuePtr(ParseInt[int8])
var DecodeInt16Ptr = DecodeValuePtr(ParseInt[int16])
var DecodeInt32Ptr = DecodeValuePtr(ParseInt[int32])
var DecodeInt64Ptr = DecodeValuePtr(ParseInt[int64])

var DecodeUint = DecodeValue(ParseUint[uint])
var DecodeUint8 = DecodeValue(ParseUint[uint8])
var DecodeUint16 = DecodeValue(ParseUint[uint16])
var DecodeUint32 = DecodeValue(ParseUint[uint32])
var DecodeUint64 = DecodeValue(ParseUint[uint64])

var DecodeUintPtr = DecodeValuePtr(ParseUint[uint])
var DecodeUint8Ptr = DecodeValuePtr(ParseUint[uint8])
var DecodeUint16Ptr = DecodeValuePtr(ParseUint[uint16])
var DecodeUint32Ptr = DecodeValuePtr(ParseUint[uint32])
var DecodeUint64Ptr = DecodeValuePtr(ParseUint[uint64])

var DecodeFloat32 = DecodeValue[float32]
var DecodeFloat64 = DecodeValue[float64]

var DecodeFloat32Ptr = DecodeValuePtr[float32]
var DecodeFloat64Ptr = DecodeValuePtr[float64]

var DecodeString = DecodeValue(ParseString)
var DecodeStringPtr = DecodeValuePtr(ParseString)

var DecodeBytes = DecodeValue(ParseBytes)

//////////////////////////////////// array ////////////////////////////////////

// DecodeArray ...
func DecodeArray[T any](f func(d *jsontext.Decoder) (T, error)) func(d *jsontext.Decoder) ([]T, error) {
	return func(d *jsontext.Decoder) ([]T, error) {
		switch d.PeekKind() {
		case 'n':
			return nil, nil
		case '[':
			var v []T
			if err := DecodeArrayBegin(d); err != nil {
				return nil, err
			}
			for {
				if d.PeekKind() == ']' {
					break
				}
				i, err := f(d)
				if err != nil {
					return nil, err
				}
				v = append(v, i)
			}
			if err := DecodeArrayEnd(d); err != nil {
				return nil, err
			}
			return v, nil
		default:
			return nil, errutil.Explain(nil, "invalid JSON: expected array")
		}
	}
}

// DecodeObjects ...
func DecodeObjects[T Object](f func() T) func(d *jsontext.Decoder) ([]T, error) {
	return DecodeArray(DecodeObject(f))
}

var DecodeBools = DecodeArray(DecodeValue(ParseBool))
var DecodeBoolPtrs = DecodeArray(DecodeValuePtr(ParseBool))

var DecodeInts = DecodeArray(DecodeValue(ParseInt[int]))
var DecodeInt8s = DecodeArray(DecodeValue(ParseInt[int8]))
var DecodeInt16s = DecodeArray(DecodeValue(ParseInt[int16]))
var DecodeInt32s = DecodeArray(DecodeValue(ParseInt[int32]))
var DecodeInt64s = DecodeArray(DecodeValue(ParseInt[int64]))

var DecodeIntPtrs = DecodeArray(DecodeValuePtr(ParseInt[int]))
var DecodeInt8Ptrs = DecodeArray(DecodeValuePtr(ParseInt[int8]))
var DecodeInt16Ptrs = DecodeArray(DecodeValuePtr(ParseInt[int16]))
var DecodeInt32Ptrs = DecodeArray(DecodeValuePtr(ParseInt[int32]))
var DecodeInt64Ptrs = DecodeArray(DecodeValuePtr(ParseInt[int64]))

var DecodeUints = DecodeArray(DecodeValue(ParseUint[uint]))
var DecodeUint8s = DecodeArray(DecodeValue(ParseUint[uint8]))
var DecodeUint16s = DecodeArray(DecodeValue(ParseUint[uint16]))
var DecodeUint32s = DecodeArray(DecodeValue(ParseUint[uint32]))
var DecodeUint64s = DecodeArray(DecodeValue(ParseUint[uint64]))

var DecodeUintPtrs = DecodeArray(DecodeValuePtr(ParseUint[uint]))
var DecodeUint8Ptrs = DecodeArray(DecodeValuePtr(ParseUint[uint8]))
var DecodeUint16Ptrs = DecodeArray(DecodeValuePtr(ParseUint[uint16]))
var DecodeUint32Ptrs = DecodeArray(DecodeValuePtr(ParseUint[uint32]))
var DecodeUint64Ptrs = DecodeArray(DecodeValuePtr(ParseUint[uint64]))

var DecodeFloat32s = DecodeArray(DecodeValue(ParseFloat[float32]))
var DecodeFloat64s = DecodeArray(DecodeValue(ParseFloat[float64]))

var DecodeFloat32Ptrs = DecodeArray(DecodeValuePtr(ParseFloat[float32]))
var DecodeFloat64Ptrs = DecodeArray(DecodeValuePtr(ParseFloat[float64]))

var DecodeStrings = DecodeArray(DecodeValue(ParseString))
var DecodeStringPtrs = DecodeArray(DecodeValuePtr(ParseString))

////////////////////////////////// map ////////////////////////////////////////

// DecodeMap ...
func DecodeMap[K comparable, V any](d *jsontext.Decoder,
	parseKeyFn func(d *jsontext.Decoder) (K, error),
	parseValueFn func(d *jsontext.Decoder) (V, error)) (map[K]V, error) {
	switch d.PeekKind() {
	case 'n':
		return nil, nil
	case '{':
		m := make(map[K]V)
		if err := DecodeObjectBegin(d); err != nil {
			return nil, err
		}
		for {
			if d.PeekKind() == '}' {
				break
			}
			key, err := parseKeyFn(d)
			if err != nil {
				return nil, err
			}
			val, err := parseValueFn(d)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
		if err := DecodeObjectEnd(d); err != nil {
			return nil, err
		}
		return m, nil
	default:
		return nil, errutil.Explain(nil, "invalid JSON: expected object or map")
	}
}
