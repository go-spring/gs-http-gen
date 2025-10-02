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

package tidl

// MetaInfo represents metadata about the parsed document.
// It contains information about the document's metadata,
// such as its name, version, and configuration.
type MetaInfo struct {
	Name    string         `json:"name"`
	Version string         `json:"version"`
	Config  map[string]any `json:"config"`
}

// Position represents the start and stop line numbers of a parsed element.
// This allows tracing back to the original source code location.
type Position struct {
	Start int
	Stop  int
}

// Document represents the root node of the parsed file.
// It contains all top-level definitions such as constants, enums, types, and RPCs.
// Additionally, it stores any global comments that are not attached to specific nodes.
type Document struct {
	Comments []Comment
	Consts   []Const
	Enums    []Enum
	Types    []Type
	RPCs     []RPC
}

// Comment represents a single comment block or line.
// Single == true means it was parsed from a single-line comment (e.g. //).
// Single == false means it was parsed from a multi-line block comment (e.g. /* ... */).
type Comment struct {
	Text     string
	Single   bool
	Position Position
}

// Comments groups the two major comment placements:
//   - Top: comments located above a declaration.
//   - Right: comments located at the end of a declaration's line.
type Comments struct {
	Top   []Comment
	Right *Comment
}

// Const represents a constant definition in the parsed document.
type Const struct {
	Type     string   // Data type of the constant
	Name     string   // Constant name
	Value    string   // Literal value
	Position Position // Location in source code
	Comments Comments // Associated comments
}

// Enum represents an enum type definition.
type Enum struct {
	Name     string      // Enum name
	Fields   []EnumField // Enum fields
	Position Position    // Location in source code
	Comments Comments    // Associated comments
}

// EnumField represents a single field inside an enum definition.
type EnumField struct {
	Name     string   // Field name
	Value    int64    // Enum constant value
	Position Position // Location in source code
	Comments Comments // Associated comments
}

// Type represents a custom user-defined type (struct-like).
type Type struct {
	Name        string         // Type name
	OneOf       bool           // oneof
	Redefined   *RedefinedType // type a b
	GenericName *string        // Optional generic type parameter (if present)
	Fields      []TypeField    // Type fields
	Position    Position       // Location in source code
	Comments    Comments       // Associated comments
}

type RedefinedType struct {
	Name        string
	GenericType TypeDefinition
}

func (t RedefinedType) Text() string {
	return t.Name + "<" + t.GenericType.Text() + ">"
}

// TypeField represents a single field inside a user-defined type.
type TypeField struct {
	FieldType   TypeDefinition // The type of the field (primitive, user, container, etc.)
	Name        string         // Field name (maybe empty for embedded types)
	Default     *string        // Optional default value (if specified)
	Annotations []Annotation   // Field annotations (key-value metadata)
	Position    Position       // Location in source code
	Comments    Comments       // Associated comments
}

// TypeDefinition is the interface implemented by all type representations.
// The Text() method provides a human-readable representation of the type.
type TypeDefinition interface {
	Text() string
}

// EmbedType represents an embedded type field (similar to composition in Go).
type EmbedType struct {
	Name     string // Embedded type name
	Optional bool   // Whether the type is optional (nullable)
}

func (t EmbedType) Text() string {
	if t.Optional {
		return t.Name + "?"
	}
	return t.Name
}

// AnyType represents the special "any" type.
type AnyType struct{}

func (t AnyType) Text() string {
	return "any"
}

// MarshalText implements encoding.TextMarshaler for AnyType.
func (t AnyType) MarshalText() (text []byte, err error) {
	return []byte(t.Text()), nil
}

// BaseType represents a primitive type (e.g., int, string, bool).
type BaseType struct {
	Name     string // Type name
	Optional bool   // Whether the type is optional (nullable)
}

func (t BaseType) Text() string {
	if t.Optional {
		return t.Name + "?"
	}
	return t.Name
}

// UserType represents a reference to a user-defined type.
type UserType struct {
	Name     string // Type name
	Optional bool   // Whether the type is optional (nullable)
}

func (t UserType) Text() string {
	if t.Optional {
		return t.Name + "?"
	}
	return t.Name
}

// BinaryType represents the "binary" type (raw bytes).
type BinaryType struct{}

func (t BinaryType) Text() string {
	return "binary"
}

// MarshalText implements encoding.TextMarshaler for BinaryType.
func (t BinaryType) MarshalText() (text []byte, err error) {
	return []byte(t.Text()), nil
}

// MapType represents a key-value container type (map<K,V>).
type MapType struct {
	Key   string         // Key type
	Value TypeDefinition // Value type
}

func (t MapType) Text() string {
	return "map<" + t.Key + ", " + t.Value.Text() + ">"
}

// ListType represents a list container type (list<T>).
type ListType struct {
	Item TypeDefinition // Element type
}

func (t ListType) Text() string {
	return "list<" + t.Item.Text() + ">"
}

// Annotation represents metadata attached to types, fields, or RPCs.
type Annotation struct {
	Key      string   // Annotation key
	Value    *string  // Optional value
	Position Position // Location in source code
	Comments Comments // Associated comments
}

// RPC represents a remote procedure call definition.
type RPC struct {
	Name        string       // RPC name
	Request     string       // Request type (raw string)
	Response    RespType     // Response type
	Annotations []Annotation // RPC annotations (metadata)
	Position    Position     // Location in source code
	Comments    Comments     // Associated comments
}

// RespType represents the response type of an RPC.
// It supports both streaming and generic forms.
type RespType struct {
	Stream   bool   // Whether the response is a stream
	TypeName string // Base type name
	UserType *UserType
}

// Text returns a human-readable representation of the response type.
func (t RespType) Text() string {
	if t.Stream {
		return "stream<" + t.UserType.Text() + ">"
	}
	if t.UserType != nil {
		return t.TypeName + "<" + t.UserType.Text() + ">"
	}
	return t.TypeName
}
