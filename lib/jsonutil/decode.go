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
	key, err := DecodeString()(d)
	if err != nil {
		if errors.Is(err, ErrNull) {
			return "", errutil.Explain(err, "invalid JSON: expected key")
		}
		return "", err
	}
	return key, nil
}

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

// DecodeBool ...
func DecodeBool() func(d *jsontext.Decoder) (bool, error) {
	return DecodeValue(ParseBool)
}

// DecodeBoolPtr ...
func DecodeBoolPtr() func(d *jsontext.Decoder) (*bool, error) {
	return DecodeValuePtr(ParseBool)
}

// DecodeInt ...
func DecodeInt[T int | int8 | int16 | int32 | int64]() func(d *jsontext.Decoder) (T, error) {
	return DecodeValue(ParseInt[T])
}

// DecodeIntPtr ...
func DecodeIntPtr[T int | int8 | int16 | int32 | int64]() func(d *jsontext.Decoder) (*T, error) {
	return DecodeValuePtr(ParseInt[T])
}

// DecodeUint ...
func DecodeUint[T uint | uint8 | uint16 | uint32 | uint64]() func(d *jsontext.Decoder) (T, error) {
	return DecodeValue(ParseUint[T])
}

// DecodeUintPtr ...
func DecodeUintPtr[T uint | uint8 | uint16 | uint32 | uint64]() func(d *jsontext.Decoder) (*T, error) {
	return DecodeValuePtr(ParseUint[T])
}

// DecodeFloat ...
func DecodeFloat[T float32 | float64]() func(d *jsontext.Decoder) (T, error) {
	return DecodeValue(ParseFloat[T])
}

// DecodeFloatPtr ...
func DecodeFloatPtr[T float32 | float64]() func(d *jsontext.Decoder) (*T, error) {
	return DecodeValuePtr(ParseFloat[T])
}

// DecodeString ...
func DecodeString() func(d *jsontext.Decoder) (string, error) {
	return DecodeValue(ParseString)
}

// DecodeStringPtr ...
func DecodeStringPtr() func(d *jsontext.Decoder) (*string, error) {
	return DecodeValuePtr(ParseString)
}

// DecodeBytes ...
func DecodeBytes() func(d *jsontext.Decoder) ([]byte, error) {
	return DecodeValue(ParseBytes)
}

// DecodeArray ...
func DecodeArray[T any](d *jsontext.Decoder, f func(d *jsontext.Decoder) (T, error)) ([]T, error) {
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

// DecodeObjects ...
func DecodeObjects[T Object](d *jsontext.Decoder, ctorFn func() T) ([]T, error) {
	return DecodeArray(d, DecodeObject(ctorFn))
}

// DecodeValues ...
func DecodeValues[T any](d *jsontext.Decoder, parseFn func(string) (T, error)) ([]T, error) {
	return DecodeArray(d, DecodeValue(parseFn))
}

// DecodeBoolsV2 ...
func DecodeBoolsV2(d *jsontext.Decoder) ([]bool, error) {
	return DecodeValues[bool](d, ParseBool)
}

// DecodeBools ...
func DecodeBools(d *jsontext.Decoder) ([]bool, error) {
	return DecodeBoolsV2(d)
	//if err := DecodeArrayBegin(d); err != nil {
	//	return nil, err
	//}
	//var v []bool
	//for {
	//	if d.PeekKind() == ']' {
	//		break
	//	}
	//	i, err := DecodeBool(d)
	//	if err != nil {
	//		if errors.Is(err, ErrNull) {
	//			return nil, errutil.Explain(nil, "null value is not allowed")
	//		}
	//		return nil, err
	//	}
	//	v = append(v, i)
	//}
	//if err := DecodeArrayEnd(d); err != nil {
	//	return nil, err
	//}
	//return v, nil
}

// DecodeIntsV2 ...
func DecodeIntsV2[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) ([]T, error) {
	return DecodeValues(d, ParseInt[T])
}

// DecodeInts ...
func DecodeInts[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) ([]T, error) {
	return DecodeIntsV2[T](d)
	//if err := DecodeArrayBegin(d); err != nil {
	//	return nil, err
	//}
	//var v []T
	//for {
	//	if d.PeekKind() == ']' {
	//		break
	//	}
	//	i, err := DecodeInt[T](d)
	//	if err != nil {
	//		if errors.Is(err, ErrNull) {
	//			return nil, errutil.Explain(nil, "null value is not allowed")
	//		}
	//		return nil, err
	//	}
	//	v = append(v, i)
	//}
	//if err := DecodeArrayEnd(d); err != nil {
	//	return nil, err
	//}
	//return v, nil
}

// DecodeUintsV2 ...
func DecodeUintsV2[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) ([]T, error) {
	return DecodeValues(d, ParseUint[T])
}

// DecodeUints ...
func DecodeUints[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) ([]T, error) {
	return DecodeUintsV2[T](d)
	//if err := DecodeArrayBegin(d); err != nil {
	//	return nil, err
	//}
	//var v []T
	//for {
	//	if d.PeekKind() == ']' {
	//		break
	//	}
	//	i, err := DecodeUint[T](d)
	//	if err != nil {
	//		if errors.Is(err, ErrNull) {
	//			return nil, errutil.Explain(nil, "null value is not allowed")
	//		}
	//		return nil, err
	//	}
	//	v = append(v, i)
	//}
	//if err := DecodeArrayEnd(d); err != nil {
	//	return nil, err
	//}
	//return v, nil
}

// DecodeFloatsV2 ...
func DecodeFloatsV2[T float32 | float64](d *jsontext.Decoder) ([]T, error) {
	return DecodeValues(d, ParseFloat[T])
}

// DecodeFloats ...
func DecodeFloats[T float32 | float64](d *jsontext.Decoder) ([]T, error) {
	return DecodeFloatsV2[T](d)
	//if err := DecodeArrayBegin(d); err != nil {
	//	return nil, err
	//}
	//var v []T
	//for {
	//	if d.PeekKind() == ']' {
	//		break
	//	}
	//	i, err := DecodeFloat[T](d)
	//	if err != nil {
	//		if errors.Is(err, ErrNull) {
	//			return nil, errutil.Explain(nil, "null value is not allowed")
	//		}
	//	}
	//	v = append(v, i)
	//}
	//if err := DecodeArrayEnd(d); err != nil {
	//	return nil, err
	//}
	//return v, nil
}

// DecodeStringsV2 ...
func DecodeStringsV2(d *jsontext.Decoder) ([]string, error) {
	return DecodeValues(d, ParseString)
}

// DecodeStrings ...
func DecodeStrings(d *jsontext.Decoder) ([]string, error) {
	return DecodeStringsV2(d)
	//if err := DecodeArrayBegin(d); err != nil {
	//	return nil, err
	//}
	//var v []string
	//for {
	//	if d.PeekKind() == ']' {
	//		break
	//	}
	//	s, err := DecodeString(d)
	//	if err != nil {
	//		if errors.Is(err, ErrNull) {
	//			return nil, errutil.Explain(nil, "null value is not allowed")
	//		}
	//		return nil, err
	//	}
	//	v = append(v, s)
	//}
	//if err := DecodeArrayEnd(d); err != nil {
	//	return nil, err
	//}
	//return v, nil
}

// DecodeValuePtrs ...
func DecodeValuePtrs[T any](d *jsontext.Decoder, parseFn func(string) (T, error)) ([]*T, error) {
	return DecodeArray(d, DecodeValuePtr(parseFn))
}

// DecodeBoolPtrsV2 ...
func DecodeBoolPtrsV2(d *jsontext.Decoder) ([]*bool, error) {
	return DecodeValuePtrs(d, ParseBool)
}

// DecodeBoolPtrs ...
func DecodeBoolPtrs(d *jsontext.Decoder) ([]*bool, error) {
	return DecodeBoolPtrsV2(d)
	//if err := DecodeArrayBegin(d); err != nil {
	//	return nil, err
	//}
	//var v []*bool
	//for {
	//	if d.PeekKind() == ']' {
	//		break
	//	}
	//	i, err := DecodeBoolPtr(d)
	//	if err != nil {
	//		return nil, err
	//	}
	//	v = append(v, i)
	//}
	//if err := DecodeArrayEnd(d); err != nil {
	//	return nil, err
	//}
	//return v, nil
}

// DecodeIntPtrsV2 ...
func DecodeIntPtrsV2[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) ([]*T, error) {
	return DecodeValuePtrs(d, ParseInt[T])
}

// DecodeIntPtrs ...
func DecodeIntPtrs[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) ([]*T, error) {
	return DecodeIntPtrsV2[T](d)
	//if err := DecodeArrayBegin(d); err != nil {
	//	return nil, err
	//}
	//var v []*T
	//for {
	//	if d.PeekKind() == ']' {
	//		break
	//	}
	//	i, err := DecodeIntPtr[T](d)
	//	if err != nil {
	//		return nil, err
	//	}
	//	v = append(v, i)
	//}
	//if err := DecodeArrayEnd(d); err != nil {
	//	return nil, err
	//}
	//return v, nil
}

// DecodeUintPtrsV2 ...
func DecodeUintPtrsV2[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) ([]*T, error) {
	return DecodeValuePtrs(d, ParseUint[T])
}

// DecodeUintPtrs ...
func DecodeUintPtrs[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) ([]*T, error) {
	return DecodeUintPtrsV2[T](d)
	//if err := DecodeArrayBegin(d); err != nil {
	//	return nil, err
	//}
	//var v []*T
	//for {
	//	if d.PeekKind() == ']' {
	//		break
	//	}
	//	i, err := DecodeUintPtr[T](d)
	//	if err != nil {
	//		return nil, err
	//	}
	//	v = append(v, i)
	//}
	//if err := DecodeArrayEnd(d); err != nil {
	//	return nil, err
	//}
	//return v, nil
}

func DecodeFloatPtrsV2[T float32 | float64](d *jsontext.Decoder) ([]*T, error) {
	return DecodeValuePtrs(d, ParseFloat[T])
}

// DecodeFloatPtrs ...
func DecodeFloatPtrs[T float32 | float64](d *jsontext.Decoder) ([]*T, error) {
	return DecodeFloatPtrsV2[T](d)
	//if err := DecodeArrayBegin(d); err != nil {
	//	return nil, err
	//}
	//var v []*T
	//for {
	//	if d.PeekKind() == ']' {
	//		break
	//	}
	//	i, err := DecodeFloatPtr[T](d)
	//	if err != nil {
	//		return nil, err
	//	}
	//	v = append(v, i)
	//}
	//if err := DecodeArrayEnd(d); err != nil {
	//	return nil, err
	//}
	//return v, nil
}

// DecodeStringPtrsV2 ...
func DecodeStringPtrsV2(d *jsontext.Decoder) ([]*string, error) {
	return DecodeValuePtrs(d, ParseString)
}

// DecodeStringPtrs ...
func DecodeStringPtrs(d *jsontext.Decoder) ([]*string, error) {
	return DecodeStringPtrsV2(d)
	//if err := DecodeArrayBegin(d); err != nil {
	//	return nil, err
	//}
	//var v []*string
	//for {
	//	if d.PeekKind() == ']' {
	//		break
	//	}
	//	s, err := DecodeStringPtr(d)
	//	if err != nil {
	//		return nil, err
	//	}
	//	v = append(v, s)
	//}
	//if err := DecodeArrayEnd(d); err != nil {
	//	return nil, err
	//}
	//return v, nil
}

func DecodeMap[K comparable, V any](d *jsontext.Decoder,
	parseKeyFn func(d *jsontext.Decoder) (K, error),
	parseValueFn func(d *jsontext.Decoder) (V, error)) (any, error) {
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
