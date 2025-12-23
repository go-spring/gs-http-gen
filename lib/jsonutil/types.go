package jsonutil

import (
	"encoding/base64"
	"math"
	"strconv"

	"github.com/lvan100/golib/errutil"
)

// ParseBool ...
func ParseBool(s string, k Kind) (bool, error) {
	if k != 'f' && k != 't' {
		return false, errutil.Explain(nil, "invalid JSON: expected boolean")
	}
	return strconv.ParseBool(s)
}

// DecodeBool ...
func DecodeBool(d Decoder) (bool, error) {
	return DecodeValue(ParseBool)(d)
}

// DecodeBoolPtr ...
func DecodeBoolPtr(d Decoder) (*bool, error) {
	return DecodeValuePtr(ParseBool)(d)
}

// OverflowInt ...
func OverflowInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](v int64) bool {
	var z T
	switch any(z).(type) {
	case int:
		return v > math.MaxInt || v < math.MinInt
	case int8:
		return v > math.MaxInt8 || v < math.MinInt8
	case int16:
		return v > math.MaxInt16 || v < math.MinInt16
	case int32:
		return v > math.MaxInt32 || v < math.MinInt32
	case int64:
		return v > math.MaxInt64 || v < math.MinInt64
	}
	return false
}

// ParseInt ...
func ParseInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](s string, k Kind) (T, error) {
	if k != '0' {
		return 0, errutil.Explain(nil, "invalid JSON: expected number")
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if OverflowInt[T](v) {
		return 0, errutil.Explain(nil, "invalid JSON: number out of range")
	}
	return T(v), nil
}

// DecodeInt ...
func DecodeInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](d Decoder) (T, error) {
	return DecodeValue(ParseInt[T])(d)
}

// DecodeIntPtr ...
func DecodeIntPtr[T ~int | ~int8 | ~int16 | ~int32 | ~int64](d Decoder) (*T, error) {
	return DecodeValuePtr(ParseInt[T])(d)
}

// ParseIntKey ...
func ParseIntKey[T ~int | ~int8 | ~int16 | ~int32 | ~int64](s string, k Kind) (T, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if OverflowInt[T](v) {
		return 0, errutil.Explain(nil, "invalid JSON: number out of range")
	}
	return T(v), nil
}

// DecodeIntKey ...
func DecodeIntKey[T ~int | ~int8 | ~int16 | ~int32 | ~int64](d Decoder) (T, error) {
	return DecodeValue(ParseIntKey[T])(d)
}

// OverflowUint ...
func OverflowUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](v uint64) bool {
	var z T
	switch any(z).(type) {
	case uint:
		return v > math.MaxUint
	case uint8:
		return v > math.MaxUint8
	case uint16:
		return v > math.MaxUint16
	case uint32:
		return v > math.MaxUint32
	}
	return false
}

// ParseUint ...
func ParseUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](s string, k Kind) (T, error) {
	if k != '0' {
		return 0, errutil.Explain(nil, "invalid JSON: expected number")
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if OverflowUint[T](v) {
		return 0, errutil.Explain(nil, "invalid JSON: number out of range")
	}
	return T(v), nil
}

// DecodeUint ...
func DecodeUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](d Decoder) (T, error) {
	return DecodeValue(ParseUint[T])(d)
}

// DecodeUintPtr ...
func DecodeUintPtr[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](d Decoder) (*T, error) {
	return DecodeValuePtr(ParseUint[T])(d)
}

// ParseUintKey ...
func ParseUintKey[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](s string, k Kind) (T, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if OverflowUint[T](v) {
		return 0, errutil.Explain(nil, "invalid JSON: number out of range")
	}
	return T(v), nil
}

// DecodeUintKey ...
func DecodeUintKey[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](d Decoder) (T, error) {
	return DecodeValue(ParseUintKey[T])(d)
}

// OverflowFloat ...
func OverflowFloat[T ~float32 | ~float64](v float64) bool {
	var z T
	switch any(z).(type) {
	case float32:
		return v > math.MaxFloat32 || v < -math.MaxFloat32
	}
	return false
}

// ParseFloat ...
func ParseFloat[T ~float32 | ~float64](s string, k Kind) (T, error) {
	if k != '0' {
		return 0, errutil.Explain(nil, "invalid JSON: expected number")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	if OverflowFloat[T](f) {
		return 0, errutil.Explain(nil, "invalid JSON: number out of range")
	}
	return T(f), nil
}

// DecodeFloat ...
func DecodeFloat[T ~float32 | ~float64](d Decoder) (T, error) {
	return DecodeValue(ParseFloat[T])(d)
}

// DecodeFloatPtr ...
func DecodeFloatPtr[T ~float32 | ~float64](d Decoder) (*T, error) {
	return DecodeValuePtr(ParseFloat[T])(d)
}

// ParseString ...
func ParseString(s string, k Kind) (string, error) {
	if k != '"' {
		return "", errutil.Explain(nil, "invalid JSON: expected string")
	}
	return s, nil
}

// DecodeString ...
func DecodeString(d Decoder) (string, error) {
	return DecodeValue(ParseString)(d)
}

// DecodeStringPtr ...
func DecodeStringPtr(d Decoder) (*string, error) {
	return DecodeValuePtr(ParseString)(d)
}

// ParseBytes ...
func ParseBytes(s string, k Kind) ([]byte, error) {
	if k != '"' {
		return nil, errutil.Explain(nil, "invalid JSON: expected string")
	}
	return base64.StdEncoding.DecodeString(s)
}

// DecodeBytes ...
func DecodeBytes(d Decoder) ([]byte, error) {
	return DecodeValue(ParseBytes)(d)
}
