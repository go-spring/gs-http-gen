/*
 * Copyright 2025 The Go-Spring Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jsonutil

import (
	"encoding/base64"
	"math"
	"strconv"

	"github.com/lvan100/golib/errutil"
)

// ParseBool parses a JSON boolean token into a Go bool.
// The input Kind must be 't' or 'f', otherwise an error is returned.
func ParseBool(_ string, k Kind) (bool, error) {
	if k != 'f' && k != 't' {
		return false, errutil.Explain(nil, "invalid JSON: expected boolean")
	}
	return k == 't', nil
}

// DecodeBool reads the next JSON value from the decoder and parses it as bool.
func DecodeBool(d Decoder) (bool, error) {
	return DecodeValue(ParseBool)(d)
}

// DecodeBoolPtr reads the next JSON value and parses it as *bool.
// Returns nil if the JSON token is null.
func DecodeBoolPtr(d Decoder) (*bool, error) {
	return DecodeValuePtr(ParseBool)(d)
}

// OverflowInt checks whether an int64 value exceeds the bounds of the target integer type T.
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

// ParseInt parses a JSON number token into an integer type T.
// Returns an error if the token is not a number or if the value overflows.
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

// DecodeInt reads the next JSON value and parses it into an integer type T.
func DecodeInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](d Decoder) (T, error) {
	return DecodeValue(ParseInt[T])(d)
}

// DecodeIntPtr reads the next JSON value and parses it into a pointer to integer type T.
// Returns nil if the JSON token is null.
func DecodeIntPtr[T ~int | ~int8 | ~int16 | ~int32 | ~int64](d Decoder) (*T, error) {
	return DecodeValuePtr(ParseInt[T])(d)
}

// ParseIntKey parses a JSON object key as an integer type T.
// Returns an error if parsing fails or the value overflows.
func ParseIntKey[T ~int | ~int8 | ~int16 | ~int32 | ~int64](s string, _ Kind) (T, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if OverflowInt[T](v) {
		return 0, errutil.Explain(nil, "invalid JSON: number out of range")
	}
	return T(v), nil
}

// DecodeIntKey reads a JSON object key and parses it as an integer type T.
func DecodeIntKey[T ~int | ~int8 | ~int16 | ~int32 | ~int64](d Decoder) (T, error) {
	return DecodeValue(ParseIntKey[T])(d)
}

// OverflowUint checks whether a uint64 value exceeds the bounds of the target unsigned type T.
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

// ParseUint parses a JSON number token into an unsigned integer type T.
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

// DecodeUint reads the next JSON value and parses it as an unsigned integer type T.
func DecodeUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](d Decoder) (T, error) {
	return DecodeValue(ParseUint[T])(d)
}

// DecodeUintPtr reads the next JSON value and parses it into a pointer to unsigned type T.
func DecodeUintPtr[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](d Decoder) (*T, error) {
	return DecodeValuePtr(ParseUint[T])(d)
}

// ParseUintKey parses a JSON object key as an unsigned integer type T.
func ParseUintKey[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](s string, _ Kind) (T, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if OverflowUint[T](v) {
		return 0, errutil.Explain(nil, "invalid JSON: number out of range")
	}
	return T(v), nil
}

// DecodeUintKey reads a JSON object key and parses it as an unsigned integer type T.
func DecodeUintKey[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](d Decoder) (T, error) {
	return DecodeValue(ParseUintKey[T])(d)
}

// OverflowFloat checks whether a float64 value exceeds the bounds of the target float type T.
func OverflowFloat[T ~float32 | ~float64](v float64) bool {
	var z T
	switch any(z).(type) {
	case float32:
		return v > math.MaxFloat32 || v < -math.MaxFloat32
	}
	return false
}

// ParseFloat parses a JSON number token into a float type T.
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

// DecodeFloat reads the next JSON value and parses it as a float type T.
func DecodeFloat[T ~float32 | ~float64](d Decoder) (T, error) {
	return DecodeValue(ParseFloat[T])(d)
}

// DecodeFloatPtr reads the next JSON value and parses it into a pointer to float type T.
func DecodeFloatPtr[T ~float32 | ~float64](d Decoder) (*T, error) {
	return DecodeValuePtr(ParseFloat[T])(d)
}

// ParseString parses a JSON string token into a Go string.
func ParseString(s string, k Kind) (string, error) {
	if k != '"' {
		return "", errutil.Explain(nil, "invalid JSON: expected string")
	}
	return s, nil
}

// DecodeString reads the next JSON value and parses it as a string.
func DecodeString(d Decoder) (string, error) {
	return DecodeValue(ParseString)(d)
}

// DecodeStringPtr reads the next JSON value and parses it as a pointer to string.
func DecodeStringPtr(d Decoder) (*string, error) {
	return DecodeValuePtr(ParseString)(d)
}

// ParseBytes parses a JSON string token as base64-encoded bytes.
func ParseBytes(s string, k Kind) ([]byte, error) {
	if k != '"' {
		return nil, errutil.Explain(nil, "invalid JSON: expected string")
	}
	return base64.StdEncoding.DecodeString(s)
}

// DecodeBytes reads the next JSON value and parses it as base64-decoded bytes.
func DecodeBytes(d Decoder) ([]byte, error) {
	return DecodeValue(ParseBytes)(d)
}
