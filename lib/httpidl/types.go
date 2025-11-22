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

package httpidl

import (
	"github.com/go-spring/gs-http-gen/lib/pathidl"
	"github.com/go-spring/gs-http-gen/lib/validate"
)

// MetaInfo represents metadata about the parsed document.
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

// Comment represents a single comment block or line.
// Single == true means it was parsed from a single-line comment (e.g. //).
// Single == false means it was parsed from a multi-line block comment (e.g. /* ... */).
type Comment struct {
	Text     []string
	Single   bool
	Position Position
}

// Comments groups the two major comment placements:
//   - Top: comments located above a declaration.
//   - Right: comments located at the end of a declaration's line.
type Comments struct {
	Above []Comment
	Right *Comment
}

// Exists returns true if there are any comments associated with this node.
func (c Comments) Exists() bool {
	return len(c.Above) > 0 || c.Right != nil
}

// Document represents the root node of the parsed file.
// It contains all top-level definitions such as constants, enums, types, RPCs, and SSEs.
// Additionally, it stores any global comments that are not attached to specific nodes.
type Document struct {
	Comments []Comment
	Consts   []Const
	Enums    []Enum
	Types    []Type
	RPCs     []RPC
	SSEs     []SSE

	EnumTypes map[string]int // Name -> index
	TypeTypes map[string]int // Name -> index
	UserTypes map[string]struct{}
}

// Annotation represents metadata attached to types, fields, RPCs, or SSEs.
type Annotation struct {
	Key      string   // Annotation key
	Value    *string  // Optional annotation value
	Position Position // Location in source code
	Comments Comments // Associated comments
}

// Const represents a constant definition in the parsed document.
type Const struct {
	Type     BaseType // Data type of the constant
	Name     string   // Name of the constant
	Value    string   // Literal value
	Position Position // Location in source code
	Comments Comments // Associated comments
}

// Enum represents an enum type definition.
type Enum struct {
	Name     string      // Name of the enum
	OneOf    bool        // Indicates whether this enum is used in oneof
	Fields   []EnumField // List of fields
	Position Position    // Location in source code
	Comments Comments    // Associated comments
}

// EnumField represents a single field inside an enum definition.
type EnumField struct {
	Name     string   // Name of the enum field
	Value    int64    // Integer value assigned to the enum field
	Position Position // Location in source code
	Comments Comments // Associated comments
}

// TypeDefinition is the interface implemented by all type representations.
// The Text() method returns a human-readable representation of the type.
type TypeDefinition interface {
	Text() string
}

// JSONTag represents the JSON tag of a field.
type JSONTag struct {
	Name      string
	OmitEmpty bool
}

// FormTag represents the form tag of a field.
type FormTag struct {
	Name string
}

// Binding represents a field binding from path, or query
type Binding struct {
	From string // Source: path/query
	Name string // Field name in the source
}

// Type represents a custom user-defined type (struct-like).
type Type struct {
	Name         string      // Name of the type
	OneOf        bool        // Indicates whether this type is a oneof
	InstType     *InstType   // Represents a type alias (e.g., type A B<T>)
	GenericParam *string     // Optional generic type parameter (if present)
	Fields       []TypeField // Type fields
	Position     Position    // Location in source code
	Comments     Comments    // Associated comments

	Embedded  bool // Embedded
	Validate  bool // Validate
	Request   bool // Request
	OnRequest bool // OnRequest
	OnForm    bool // OnForm
}

// TypeField represents a single field inside a user-defined type.
type TypeField struct {
	Name        string         // Name of the field
	Type        TypeDefinition // Type of the field
	Annotations []Annotation   // Additional metadata (key-value pairs)
	Position    Position       // Location in source code
	Comments    Comments       // Associated comments

	JSONTag        JSONTag       // JSON tag
	FormTag        FormTag       // Form tag
	Binding        *Binding      // Field binding
	Required       bool          // Required
	ValidateExpr   validate.Expr // Validate expression
	ValidateNested bool          // Nested validate
	EnumAsString   bool          // Enum as string
}

// InstType represents a type alias with optional generic arguments.
type InstType struct {
	Name        string         // Name of the aliased type
	GenericType TypeDefinition // The generic type parameter
}

func (t InstType) Text() string {
	return t.Name + "<" + t.GenericType.Text() + ">"
}

// EmbedType represents an embedded type field (similar to composition in Go).
type EmbedType struct {
	Name string // Name of the embedded type
}

func (t EmbedType) Text() string {
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
	Name string // Name of the primitive type
}

func (t BaseType) Text() string {
	return t.Name
}

// UserType represents a reference to a user-defined type.
type UserType struct {
	Name string // Name of the referenced type
}

func (t UserType) Text() string {
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

// RPC represents a remote procedure call definition.
type RPC struct {
	Name        string       // Name of the RPC
	Request     string       // Request type
	Response    string       // Response type
	Annotations []Annotation // Metadata attached to the RPC
	Position    Position     // Location in source code
	Comments    Comments     // Associated comments

	Path        string // HTTP path
	Method      string // HTTP method (GET, POST, etc.)
	ContentType string // HTTP Content-Type

	ConnTimeout  int // Connection timeout, ms
	ReadTimeout  int // Read timeout, ms
	WriteTimeout int // Write timeout, ms

	PathSegments []pathidl.Segment
	PathParams   map[string]string // path => field name of request type
}

// SSE represents a server-sent event definition.
type SSE struct {
	Name        string       // Name of the SSE
	Request     string       // Request type
	Response    string       // Response type
	Annotations []Annotation // Metadata attached to the SSE
	Position    Position     // Location in source code
	Comments    Comments     // Associated comments

	Path        string // HTTP path
	Method      string // HTTP method (GET, POST, etc.)
	ContentType string // HTTP Content-Type

	ConnTimeout  int // Connection timeout, ms
	ReadTimeout  int // Read timeout, ms
	WriteTimeout int // Write timeout, ms

	PathSegments []pathidl.Segment
	PathParams   map[string]string // path => field name of request type
}
