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
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// JSONv2Decoder wraps jsontext.Decoder to implement the Decoder interface.
// It provides streaming JSON decoding with convenience methods for reading tokens and values.
type JSONv2Decoder struct {
	*jsontext.Decoder
}

// NewJSONv2Decoder creates a new JSONv2Decoder that implements the Decoder interface.
// It wraps an existing jsontext.Decoder.
func NewJSONv2Decoder(d *jsontext.Decoder) Decoder {
	return &JSONv2Decoder{d}
}

// Unmarshal decodes JSON bytes into a Go value using the standard jsonv2.Unmarshal.
// This satisfies the Decoder interface.
func (d *JSONv2Decoder) Unmarshal(b []byte, i any) error {
	return json.Unmarshal(b, i)
}

// toKind converts jsontext.Kind to the local Kind type used by jsonutil.
// Returns InvalidKind if the input kind does not match any known JSON token.
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

// PeekKind returns the Kind of the next JSON token without consuming it.
func (d *JSONv2Decoder) PeekKind() Kind {
	return toKind(d.Decoder.PeekKind())
}

// ReadToken reads the next JSON token and returns its string representation,
// kind, and any error encountered. The decoder state advances past the returned token.
func (d *JSONv2Decoder) ReadToken() (string, Kind, error) {
	token, err := d.Decoder.ReadToken()
	if err != nil {
		return "", 0, err
	}
	return token.String(), toKind(token.Kind()), nil
}

// ReadValue reads the next JSON value (scalar, object, or array) as raw bytes.
// The decoder state advances past the returned value.
func (d *JSONv2Decoder) ReadValue() ([]byte, error) {
	return d.Decoder.ReadValue()
}

// SkipValue skips over the next JSON value (scalar, object, or array).
// The decoder state advances past the skipped value.
func (d *JSONv2Decoder) SkipValue() error {
	return d.Decoder.SkipValue()
}
