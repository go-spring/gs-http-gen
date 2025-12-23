package jsonutil

import (
	"github.com/lvan100/golib/errutil"
)

// HashKey returns a 64-bit hash value for the given string using the FNV-1a algorithm.
// Suitable for short strings with low collision probability.
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

// Decoder defines a streaming JSON decoder interface.
type Decoder interface {
	// Unmarshal decodes JSON bytes into an arbitrary Go value.
	Unmarshal(b []byte, i any) error
	// PeekKind returns the Kind of the next token without consuming it.
	PeekKind() Kind
	// ReadToken reads the next token and returns its string value, kind, and error.
	ReadToken() (token string, _ Kind, _ error)
	// ReadValue reads the next value, which may be a complete JSON
	// node (object, array, or scalar), as bytes.
	ReadValue() (value []byte, _ error)
	// SkipValue skips the next value (maybe a complete JSON node).
	SkipValue() error
}

// Object represents a JSON-mappable object that supports streaming decoding.
type Object interface {
	// DecodeJSON reads JSON data from the Decoder and populates the object.
	DecodeJSON(d Decoder) error
}

// DecodeObjectBegin consumes the opening '{' token of a JSON object.
// Returns an error if the next token is not '{'.
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

// DecodeObjectEnd consumes the closing '}' token of a JSON object.
// Returns an error if the next token is not '}'.
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

// DecodeArrayBegin consumes the opening '[' token of a JSON array.
// Returns an error if the next token is not '['.
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

// DecodeArrayEnd consumes the closing ']' token of a JSON array.
// Returns an error if the next token is not ']'.
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

// DecodeAny decodes the next JSON value (scalar, object, or array)
// into a Go value using Decoder.Unmarshal.
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

// DecodeValue parses a scalar JSON value (number, boolean, or string) using parseFn.
// Returns an error if the next token is null or invalid.
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
			return zero, errutil.Explain(nil, "invalid JSON: unexpected null")
		case 'f', 't', '0', '"':
			v, err := parseFn(token, tokenKind)
			if err != nil {
				return zero, err
			}
			return v, nil
		default:
			return zero, errutil.Explain(err, "invalid JSON: expected value")
		}
	}
}

// DecodeValuePtr parses a scalar JSON value into a pointer type.
// Returns nil if the next token is null.
func DecodeValuePtr[T any](
	parseFn func(string, Kind) (T, error),
) func(d Decoder) (*T, error) {
	return func(d Decoder) (*T, error) {
		token, tokenKind, err := d.ReadToken()
		if err != nil {
			return nil, err
		}
		switch tokenKind {
		case 'n':
			return nil, nil
		case 'f', 't', '0', '"':
			v, err := parseFn(token, tokenKind)
			if err != nil {
				return nil, err
			}
			return &v, nil
		default:
			return nil, errutil.Explain(err, "invalid JSON: expected value")
		}
	}
}

// DecodeObject decodes a JSON object into a struct that implements the Object interface.
// Returns the zero value if the next token is null.
// Internally calls DecodeJSON on the object to populate its fields.
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

// DecodeArray decodes a JSON array of arbitrary type.
// parseFn is used to parse each element of the array.
// Returns nil if the next token is null.
func DecodeArray[T any](
	parseFn func(d Decoder) (T, error),
) func(d Decoder) ([]T, error) {
	return func(d Decoder) ([]T, error) {
		switch d.PeekKind() {
		case 'n':
			_, _, _ = d.ReadToken()
			return nil, nil
		case '[':
			v := make([]T, 0)
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

// DecodeMap decodes a JSON object into a Go map.
// parseKeyFn and parseValFn are used to parse each key and value.
// Returns nil if the next token is null.
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
			return nil, errutil.Explain(nil, "invalid JSON: expected object")
		}
	}
}
