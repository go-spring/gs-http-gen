package jsonutil

import (
	"errors"

	"github.com/lvan100/golib/errutil"
)

// HashKey returns a hash value for the given string.
// 业界推荐算法，对于短字符串，碰撞的概率很低很低.
func HashKey(s string) uint64 {
	const (
		offset = 14695981039346656037
		prime  = 1099511628211
	)
	h := uint64(offset)
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime
	}
	return h
}

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
	Unmarshal(b []byte, i any) error
	PeekKind() Kind
	ReadToken() (token string, _ Kind, _ error)
	ReadValue() (value []byte, _ error)
	SkipValue() error
}

var ErrNull = errors.New("json: null")

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

// DecodeAny ...
func DecodeAny[T any](d Decoder) (T, error) {
	var v T
	b, err := d.ReadValue()
	if err != nil {
		return v, err
	}
	if err = d.Unmarshal(b, &v); err != nil {
		return v, err
	}
	return v, nil
}

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
	newFn func() T,
) func(d Decoder) (T, error) {
	return func(d Decoder) (T, error) {
		var zero T
		switch d.PeekKind() {
		case 'n':
			_, _, _ = d.ReadToken()
			return zero, nil
		case '{':
			v := newFn()
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
	parseFn func(d Decoder) (T, error),
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
				i, err := parseFn(d)
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
	parseValFn func(d Decoder) (V, error),
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
				val, err := parseValFn(d)
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
			return nil, errutil.Explain(nil, "invalid JSON: expected map")
		}
	}
}
