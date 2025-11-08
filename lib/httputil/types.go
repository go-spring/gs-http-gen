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

package httputil

import (
	"encoding/json"
	"strconv"

	"github.com/spf13/cast"
)

// Bool represents a nullable boolean value.
type Bool struct {
	Valid bool
	Value bool
}

// Set sets the value of the Bool.
func (b *Bool) Set(v bool) {
	b.Valid = true
	b.Value = v
}

// MarshalJSON implements the json.Marshaler interface.
func (b *Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatBool(b.Value)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *Bool) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &b.Value); err != nil {
		return err
	}
	b.Valid = true
	return nil
}

// AnyBool represents a nullable boolean value that can be of any type.
type AnyBool struct {
	Bool
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *AnyBool) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	b.Value, err = cast.ToBoolE(v)
	if err != nil {
		return err
	}
	b.Valid = true
	return nil
}

// Int represents a nullable integer value.
type Int struct {
	Valid bool
	Value int64
}

// Set sets the value of the Int.
func (i *Int) Set(v int64) {
	i.Valid = true
	i.Value = v
}

// MarshalJSON implements the json.Marshaler interface.
func (i *Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(i.Value, 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (i *Int) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &i.Value); err != nil {
		return err
	}
	i.Valid = true
	return nil
}

// AnyInt represents a nullable integer value that can be of any type.
type AnyInt struct {
	Int
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (i *AnyInt) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	i.Value, err = cast.ToInt64E(v)
	if err != nil {
		return err
	}
	i.Valid = true
	return nil
}

// Uint represents a nullable unsigned integer value.
type Uint struct {
	Valid bool
	Value uint64
}

// Set sets the value of the Uint.
func (u *Uint) Set(v uint64) {
	u.Valid = true
	u.Value = v
}

// MarshalJSON implements the json.Marshaler interface.
func (u *Uint) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(u.Value, 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (u *Uint) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &u.Value); err != nil {
		return err
	}
	u.Valid = true
	return nil
}

// AnyUint represents a nullable unsigned integer value that can be of any type.
type AnyUint struct {
	Uint
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (u *AnyUint) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	u.Value, err = cast.ToUint64E(v)
	if err != nil {
		return err
	}
	u.Valid = true
	return nil
}

// Float represents a nullable floating-point value.
type Float struct {
	Valid bool
	Value float64
}

// Set sets the value of the Float.
func (f *Float) Set(v float64) {
	f.Valid = true
	f.Value = v
}

// MarshalJSON implements the json.Marshaler interface.
func (f *Float) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatFloat(f.Value, 'f', -1, 64)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *Float) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &f.Value); err != nil {
		return err
	}
	f.Valid = true
	return nil
}

// AnyFloat represents a nullable floating-point value that can be of any type.
type AnyFloat struct {
	Float
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *AnyFloat) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	f.Value, err = cast.ToFloat64E(v)
	if err != nil {
		return err
	}
	f.Valid = true
	return nil
}

// String represents a nullable string value.
type String struct {
	Valid bool
	Value string
}

// Set sets the value of the String.
func (s *String) Set(v string) {
	s.Valid = true
	s.Value = v
}

// MarshalJSON implements the json.Marshaler interface.
func (s *String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.Quote(s.Value)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *String) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.Value); err != nil {
		return err
	}
	s.Valid = true
	return nil
}

// AnyString represents a nullable string value that can be of any type.
type AnyString struct {
	String
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *AnyString) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	s.Value, err = cast.ToStringE(v)
	if err != nil {
		return err
	}
	s.Valid = true
	return nil
}
