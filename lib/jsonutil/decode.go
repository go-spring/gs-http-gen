package jsonutil

import (
	"encoding/base64"
	"encoding/json/jsontext"
	"errors"
	"math"
	"strconv"

	"github.com/lvan100/golib/errutil"
)

// Kind represents each possible JSON token kind with a single byte,
// which is conveniently the first byte of that kind's grammar
// with the restriction that numbers always be represented with '0':
//
//   - 'n': null
//   - 'f': false
//   - 't': true
//   - '"': string
//   - '0': number
//   - '{': object begin
//   - '}': object end
//   - '[': array begin
//   - ']': array end
//
// An invalid kind is usually represented using 0,
// but may be non-zero due to invalid JSON data.
type Kind byte

const InvalidKind Kind = 0

// Decoder ...
type Decoder interface {
	PeekKind() Kind
	ReadToken() (string, Kind, error)
	SkipValue() error
}

// JSONv2Decoder ...
type JSONv2Decoder struct {
	*jsontext.Decoder
}

// NewJSONv2Decoder ...
func NewJSONv2Decoder(d *jsontext.Decoder) Decoder {
	return &JSONv2Decoder{d}
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

// SkipValue ...
func (d *JSONv2Decoder) SkipValue() error {
	return d.Decoder.SkipValue()
}

var ErrNull = errors.New("null")

// Object ...
type Object interface {
	DecodeJSON(d Decoder) error
}

// DecodeObjectBegin ...
func DecodeObjectBegin(d Decoder) error {
	_, tokenKind, err := d.ReadToken()
	if err != nil {
		return err
	}
	if tokenKind != '{' {
		return errutil.Explain(err, "invalid JSON: expected object")
	}
	return nil
}

// DecodeObjectEnd ...
func DecodeObjectEnd(d Decoder) error {
	_, tokenKind, err := d.ReadToken()
	if err != nil {
		return err
	}
	if tokenKind != '}' {
		return errutil.Explain(err, "invalid JSON: expected end of object")
	}
	return nil
}

// DecodeArrayBegin ...
func DecodeArrayBegin(d Decoder) error {
	_, tokenKind, err := d.ReadToken()
	if err != nil {
		return err
	}
	if tokenKind != '[' {
		return errutil.Explain(err, "invalid JSON: expected array")
	}
	return nil
}

// DecodeArrayEnd ...
func DecodeArrayEnd(d Decoder) error {
	_, tokenKind, err := d.ReadToken()
	if err != nil {
		return err
	}
	if tokenKind != ']' {
		return errutil.Explain(err, "invalid JSON: expected end of array")
	}
	return nil
}

// DecodeKey ...
func DecodeKey(d Decoder) (string, error) {
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
func ParseBool(s string, k Kind) (bool, error) {
	if k != 'f' && k != 't' {
		return false, errutil.Explain(nil, "invalid JSON: expected boolean")
	}
	return strconv.ParseBool(s)
}

// OverflowInt ...
func OverflowInt[T int | int8 | int16 | int32 | int64](v int64) bool {
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
func ParseInt[T int | int8 | int16 | int32 | int64](s string, k Kind) (T, error) {
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

// OverflowUint ...
func OverflowUint[T uint | uint8 | uint16 | uint32 | uint64](v uint64) bool {
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
func ParseUint[T uint | uint8 | uint16 | uint32 | uint64](s string, k Kind) (T, error) {
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

// OverflowFloat ...
func OverflowFloat[T float32 | float64](v float64) bool {
	var z T
	switch any(z).(type) {
	case float32:
		return v > math.MaxFloat32 || v < -math.MaxFloat32
	}
	return false
}

// ParseFloat ...
func ParseFloat[T float32 | float64](s string, k Kind) (T, error) {
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

// ParseString ...
func ParseString(s string, k Kind) (string, error) {
	if k != '"' {
		return "", errutil.Explain(nil, "invalid JSON: expected string")
	}
	return s, nil
}

// ParseBytes ...
func ParseBytes(s string, k Kind) ([]byte, error) {
	if k != '"' {
		return nil, errutil.Explain(nil, "invalid JSON: expected string")
	}
	return base64.StdEncoding.DecodeString(s)
}

//////////////////////////////////// decode ////////////////////////////////////

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

// DecodeValue ...
func DecodeValue[T any](
	parseFn func(string, Kind) (T, error),
) func(d Decoder) (T, error) {
	return func(d Decoder) (T, error) {
		var zero T
		token, tokenKind, err := d.ReadToken()
		if err != nil {
			return zero, err
		}
		switch tokenKind {
		case 'n':
			return zero, ErrNull
		case 'f', 't', '0', '"':
			return parseFn(token, tokenKind)
		default:
			return zero, errutil.Explain(err, "invalid JSON: expected value")
		}
	}
}

// DecodeValuePtr ...
func DecodeValuePtr[T any](
	parseFn func(string, Kind) (T, error),
) func(d Decoder) (*T, error) {
	return func(d Decoder) (*T, error) {
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

// DecodeObject ...
func DecodeObject[T Object](
	f func() T,
) func(d Decoder) (T, error) {
	return func(d Decoder) (T, error) {
		var zero T
		switch d.PeekKind() {
		case 'n':
			_, _, _ = d.ReadToken()
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

// DecodeArray ...
func DecodeArray[T any](
	f func(d Decoder) (T, error),
) func(d Decoder) ([]T, error) {
	return func(d Decoder) ([]T, error) {
		switch d.PeekKind() {
		case 'n':
			_, _, _ = d.ReadToken()
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

// DecodeMap ...
func DecodeMap[K comparable, V any](
	parseKeyFn func(d Decoder) (K, error),
	parseValueFn func(d Decoder) (V, error),
) func(d Decoder) (map[K]V, error) {
	return func(d Decoder) (map[K]V, error) {
		switch d.PeekKind() {
		case 'n':
			_, _, _ = d.ReadToken()
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
}
