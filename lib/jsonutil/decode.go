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

// DecodeValue ...
func DecodeValue[T any](d *jsontext.Decoder, parseFn func(string) (T, error)) (T, error) {
	var v T
	value, err := d.ReadValue()
	if err != nil {
		return v, err
	}
	switch value.Kind() {
	case 'n':
		return v, ErrNull
	case 't', 'f', '0':
		return parseFn(value.String())
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

// DecodeBoolV2 ...
func DecodeBoolV2(d *jsontext.Decoder) (bool, error) {
	return DecodeValue(d, func(s string) (bool, error) {
		return strconv.ParseBool(s)
	})
}

// DecodeBool ...
func DecodeBool(d *jsontext.Decoder) (bool, error) {
	return DecodeBoolV2(d)
	//value, err := d.ReadValue()
	//if err != nil {
	//	return false, err
	//}
	//switch value.Kind() {
	//case 't':
	//	return true, nil
	//case 'f':
	//	return false, nil
	//case 'n':
	//	return false, ErrNull
	//case '"':
	//	var s string
	//	s, err = strconv.Unquote(value.String())
	//	if err != nil {
	//		return false, err
	//	}
	//	return strconv.ParseBool(s)
	//default:
	//	return false, errutil.Explain(err, "invalid JSON: expected boolean")
	//}
}

// DecodeIntV2 ...
func DecodeIntV2[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) (T, error) {
	return DecodeValue(d, func(s string) (T, error) {
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
	return DecodeValue(d, func(s string) (T, error) {
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
	return DecodeValue(d, func(s string) (T, error) {
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

// DecodeStringV2 ...
func DecodeStringV2(d *jsontext.Decoder) (string, error) {
	return DecodeValue(d, func(s string) (string, error) {
		return s, nil
	})
}

// DecodeString decodes a string value from the given JSON decoder.
func DecodeString(d *jsontext.Decoder) (string, error) {
	return DecodeStringV2(d)
	//value, err := d.ReadValue()
	//if err != nil {
	//	return "", err
	//}
	//switch value.Kind() {
	//case 'n':
	//	return "", ErrNull
	//case '"':
	//	return strconv.Unquote(value.String())
	//default:
	//	return "", errutil.Explain(err, "invalid JSON: expected string")
	//}
}

// DecodeBytesV2 ...
func DecodeBytesV2(d *jsontext.Decoder) ([]byte, error) {
	return DecodeValue(d, func(s string) ([]byte, error) {
		return base64.StdEncoding.DecodeString(s)
	})
}

// DecodeBytes decodes a byte array value from the given JSON decoder.
func DecodeBytes(d *jsontext.Decoder) ([]byte, error) {
	return DecodeBytesV2(d)
	//value, err := d.ReadValue()
	//if err != nil {
	//	return nil, err
	//}
	//switch value.Kind() {
	//case 'n':
	//	return nil, nil
	//case '"':
	//	s, err := strconv.Unquote(value.String())
	//	if err != nil {
	//		return nil, err
	//	}
	//	return base64.StdEncoding.DecodeString(s)
	//default:
	//	return nil, errutil.Explain(err, "invalid JSON: expected []byte")
	//}
}

// DecodePtr ...
func DecodePtr[T any](d *jsontext.Decoder, parseFn func(string) (T, error)) (*T, error) {
	value, err := d.ReadValue()
	if err != nil {
		return nil, err
	}
	switch value.Kind() {
	case 'n':
		return nil, nil
	case 't', 'f', '0':
		i, err := parseFn(value.String())
		if err != nil {
			return nil, err
		}
		return &i, nil
	case '"':
		s, err := strconv.Unquote(value.String())
		if err != nil {
			return nil, err
		}
		i, err := parseFn(s)
		if err != nil {
			return nil, err
		}
		return &i, nil
	default:
		return nil, errutil.Explain(err, "invalid JSON: expected integer")
	}
}

// DecodeBoolPtrV2 ...
func DecodeBoolPtrV2(d *jsontext.Decoder) (*bool, error) {
	return DecodePtr(d, func(s string) (bool, error) {
		return strconv.ParseBool(s)
	})
}

// DecodeBoolPtr ...
func DecodeBoolPtr(d *jsontext.Decoder) (*bool, error) {
	return DecodeBoolPtrV2(d)
	//b, err := DecodeBool(d)
	//if err != nil {
	//	if errors.Is(err, ErrNull) {
	//		return nil, nil
	//	}
	//	return nil, err
	//}
	//return &b, nil
}

// DecodeIntPtrV2 ...
func DecodeIntPtrV2[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) (*T, error) {
	return DecodePtr(d, func(s string) (T, error) {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
}

// DecodeIntPtr ...
func DecodeIntPtr[T int | int8 | int16 | int32 | int64](d *jsontext.Decoder) (*T, error) {
	return DecodeIntPtrV2[T](d)
	//i, err := DecodeInt[T](d)
	//if err != nil {
	//	if errors.Is(err, ErrNull) {
	//		return nil, nil
	//	}
	//	return nil, err
	//}
	//return &i, nil
}

// DecodeUintPtrV2 ...
func DecodeUintPtrV2[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) (*T, error) {
	return DecodePtr(d, func(s string) (T, error) {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
}

// DecodeUintPtr ...
func DecodeUintPtr[T uint | uint8 | uint16 | uint32 | uint64](d *jsontext.Decoder) (*T, error) {
	return DecodeUintPtrV2[T](d)
	//u, err := DecodeUint[T](d)
	//if err != nil {
	//	if errors.Is(err, ErrNull) {
	//		return nil, nil
	//	}
	//	return nil, err
	//}
	//return &u, nil
}

// DecodeFloatPtrV2 ...
func DecodeFloatPtrV2[T float32 | float64](d *jsontext.Decoder) (*T, error) {
	return DecodePtr(d, func(s string) (T, error) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
}

// DecodeFloatPtr ...
func DecodeFloatPtr[T float32 | float64](d *jsontext.Decoder) (*T, error) {
	return DecodeFloatPtrV2[T](d)
	//f, err := DecodeFloat[T](d)
	//if err != nil {
	//	if errors.Is(err, ErrNull) {
	//		return nil, nil
	//	}
	//	return nil, err
	//}
	//return &f, nil
}

// DecodeStringPtrV2 ...
func DecodeStringPtrV2(d *jsontext.Decoder) (*string, error) {
	return DecodePtr(d, func(s string) (string, error) {
		return s, nil
	})
}

// DecodeStringPtr ...
func DecodeStringPtr(d *jsontext.Decoder) (*string, error) {
	return DecodeStringPtrV2(d)
	//s, err := DecodeString(d)
	//if err != nil {
	//	if errors.Is(err, ErrNull) {
	//		return nil, nil
	//	}
	//	return nil, err
	//}
	//return &s, nil
}

// DecodeArray ...
func DecodeArray[T any](d *jsontext.Decoder, parseFn func(string) (T, error)) ([]T, error) {
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
			return nil, err
		}
		var i T
		switch value.Kind() {
		case 'n':
			return nil, errutil.Explain(nil, "null value is not allowed")
		case 't', 'f', '0':
			i, err = parseFn(value.String())
			if err != nil {
				return nil, err
			}
		case '"':
			s, err := strconv.Unquote(value.String())
			if err != nil {
				return nil, err
			}
			i, err = parseFn(s)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errutil.Explain(err, "invalid JSON: expected integer")
		}
		v = append(v, i)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeBoolsV2 ...
func DecodeBoolsV2(d *jsontext.Decoder) ([]bool, error) {
	return DecodeArray[bool](d, strconv.ParseBool)
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
	return DecodeArray(d, func(s string) (T, error) {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
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
	return DecodeArray(d, func(s string) (T, error) {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
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
	return DecodeArray(d, func(s string) (T, error) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
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
	return DecodeArray(d, func(s string) (string, error) {
		return s, nil
	})
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

// DecodePtrArray ...
func DecodePtrArray[T any](d *jsontext.Decoder, parseFn func(string) (T, error)) ([]*T, error) {
	if err := DecodeArrayBegin(d); err != nil {
		return nil, err
	}
	var v []*T
	for {
		if d.PeekKind() == ']' {
			break
		}
		value, err := d.ReadValue()
		if err != nil {
			return nil, err
		}
		var p *T
		switch value.Kind() {
		case 'n':
		case 't', 'f', '0':
			i, err := parseFn(value.String())
			if err != nil {
				return nil, err
			}
			p = &i
		case '"':
			s, err := strconv.Unquote(value.String())
			if err != nil {
				return nil, err
			}
			i, err := parseFn(s)
			if err != nil {
				return nil, err
			}
			p = &i
		default:
			return nil, errutil.Explain(err, "invalid JSON: expected integer")
		}
		v = append(v, p)
	}
	if err := DecodeArrayEnd(d); err != nil {
		return nil, err
	}
	return v, nil
}

// DecodeBoolPtrsV2 ...
func DecodeBoolPtrsV2(d *jsontext.Decoder) ([]*bool, error) {
	return DecodePtrArray(d, func(s string) (bool, error) {
		return strconv.ParseBool(s)
	})
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
	return DecodePtrArray(d, func(s string) (T, error) {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
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
	return DecodePtrArray(d, func(s string) (T, error) {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
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
	return DecodePtrArray(d, func(s string) (T, error) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return T(v), nil
	})
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
	return DecodePtrArray(d, func(s string) (string, error) {
		return s, nil
	})
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

//// DecodeObjectArray ...
//func DecodeObjectArray[T Object](d *jsontext.Decoder) ([]*T, error) {
//	if err := DecodeArrayBegin(d); err != nil {
//		return nil, err
//	}
//	var v []*T
//	for {
//		if d.PeekKind() == ']' {
//			break
//		}
//		var i T
//		if err := i.DecodeJSON(d); err != nil {
//			return nil, err
//		}
//		v = append(v, &i)
//	}
//	if err := DecodeArrayEnd(d); err != nil {
//		return nil, err
//	}
//	return v, nil
//}
