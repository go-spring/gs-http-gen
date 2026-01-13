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

package golang

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-spring/gs-http-gen/lib/httpidl"
	"github.com/go-spring/gs-http-gen/lib/pathidl"
	"github.com/go-spring/stdlib/errutil"
)

// TypeKind represents kind of a Go field type
type TypeKind int

const (
	TypeKindBool = TypeKind(iota)
	TypeKindBoolPtr
	TypeKindInt
	TypeKindIntPtr
	TypeKindUint
	TypeKindUintPtr
	TypeKindFloat
	TypeKindFloatPtr
	TypeKindString
	TypeKindStringPtr
	TypeKindBytes
	TypeKindAny
	TypeKindEnum
	TypeKindEnumPtr
	TypeKindEnumAsString
	TypeKindEnumAsStringPtr
	TypeKindStructPtr
	TypeKindList
	TypeKindMap
)

// IsNotNullable returns true if the field is not nullable
func IsNotNullable(x TypeKind) bool {
	switch x {
	case TypeKindBool, TypeKindInt, TypeKindUint,
		TypeKindFloat, TypeKindString, TypeKindEnum,
		TypeKindEnumAsString:
		return true
	default:
		return false
	}
}

// IsPointer returns true if the field is a pointer
func IsPointer(x TypeKind) bool {
	switch x {
	case TypeKindBoolPtr, TypeKindIntPtr, TypeKindUintPtr,
		TypeKindFloatPtr, TypeKindStringPtr, TypeKindEnumPtr,
		TypeKindEnumAsStringPtr, TypeKindStructPtr:
		return true
	default:
		return false
	}
}

// Const represents a Go constant
type Const struct {
	httpidl.Const
	Type string
}

// Enum represents a Go enum
type Enum struct {
	httpidl.Enum
}

// Type represents a Go struct
type Type struct {
	httpidl.Type
	Fields []TypeField
}

// TypeField represents a field in a Go struct
type TypeField struct {
	httpidl.TypeField
	Name     string
	Type     string
	TypeKind []TypeKind
}

// IsPointer returns true if the field is a pointer
func (x *TypeField) IsPointer() bool {
	return IsPointer(x.TypeKind[0])
}

// FieldTag returns the field tag
func (x *TypeField) FieldTag() string {
	var tags []string

	// JSON tag
	{
		var sb strings.Builder
		sb.WriteString(`json:"`)
		sb.WriteString(x.JSONTag.Name)
		if x.JSONTag.OmitEmpty {
			sb.WriteString(",omitempty")
		}
		sb.WriteString(`"`)
		tags = append(tags, sb.String())
	}

	// Form tag
	if x.Binding == nil {
		s := fmt.Sprintf(`form:"%s"`, x.FormTag.Name)
		tags = append(tags, s)
	} else {
		s := fmt.Sprintf(`%s:"%s"`, x.Binding.Source, x.Binding.Field)
		tags = append(tags, s)
	}

	// Validate tag
	if x.Required {
		tags = append(tags, `validate:"required"`)
	}

	// Default tag
	if x.CompatDefault != nil {
		tags = append(tags, fmt.Sprintf(`default:"%s"`, *x.CompatDefault))
	}

	return "`" + strings.Join(tags, " ") + "`"
}

// RPC represents a single remote procedure call with HTTP metadata.
type RPC struct {
	httpidl.RPC
	Response     string     // Response type
	RespTypeKind []TypeKind // Response type kind

	FormatPath   string            // Formatted HTTP path
	PathParams   map[string]string // HTTP path parameters
	PathSegments []pathidl.Segment // HTTP path segments
}

type GoSpec struct {
	Meta   *httpidl.MetaInfo
	Files  map[string]httpidl.Document
	Funcs  map[string]httpidl.ValidateFunc
	Consts map[string][]Const
	Enums  map[string][]Enum
	Types  map[string][]Type
	RPCs   []RPC
}

// Convert converts an IDL project to Go code.
func Convert(dir string) (GoSpec, error) {
	project, err := httpidl.ParseDir(dir)
	if err != nil {
		return GoSpec{}, err
	}

	spec := GoSpec{
		Meta:   project.Meta,
		Files:  project.Files,
		Funcs:  project.Funcs,
		Consts: make(map[string][]Const),
		Enums:  make(map[string][]Enum),
		Types:  make(map[string][]Type),
	}

	// Collect all RPC definitions
	for _, doc := range project.Files {
		for _, r := range doc.RPCs {

			var response string
			if a, ok := httpidl.FindAnnotation(r.Annotations, "resp.go.type"); ok {
				if a.Value == nil {
					return GoSpec{}, errutil.Explain(nil, `annotation "resp.go.type" must have a value`)
				}
				s := strings.Trim(strings.TrimSpace(*a.Value), "\"")
				if s == "" {
					return GoSpec{}, errutil.Explain(nil, `annotation "resp.go.type" must not be empty`)
				}
				response = s
			} else {
				switch typ := r.Response.(type) {
				case httpidl.BytesType:
					response = "[]byte"
				case httpidl.BaseType:
					if response, err = goBaseType(typ.Name); err != nil {
						return GoSpec{}, err
					}
				case httpidl.UserType:
					if _, ok = httpidl.FindEnum(spec.Files, typ.Name); ok {
						response = typ.Name
					} else {
						response = "*" + typ.Name
					}
				default:
					if response, err = goTypeDef(spec, typ); err != nil {
						return GoSpec{}, err
					}
				}
			}

			typeKind, err := getTypeKind(spec, response)
			if err != nil {
				return GoSpec{}, err
			}

			rpc := RPC{
				RPC:          r,
				Response:     response,
				RespTypeKind: typeKind,
				FormatPath:   r.Path, // 假设是普通路径
				PathParams:   r.PathParams,
				PathSegments: r.PathSegments,
			}
			spec.RPCs = append(spec.RPCs, rpc)
		}
	}

	sort.Slice(spec.RPCs, func(i, j int) bool {
		return spec.RPCs[i].Name < spec.RPCs[j].Name
	})

	for fileName, doc := range project.Files {
		consts, err := convertConsts(spec, doc)
		if err != nil {
			return GoSpec{}, errutil.Explain(nil, "convert consts error: %w", err)
		}
		enums, err := convertEnums(spec, doc)
		if err != nil {
			return GoSpec{}, errutil.Explain(nil, "convert enums error: %w", err)
		}
		types, err := convertTypes(spec, doc)
		if err != nil {
			return GoSpec{}, errutil.Explain(nil, "convert types error: %w", err)
		}
		spec.Consts[fileName] = consts
		spec.Enums[fileName] = enums
		spec.Types[fileName] = types
	}

	for i, rpc := range spec.RPCs {
		for k, s := range rpc.PathParams {
			rpc.PathParams[k] = httpidl.ToPascal(s)
		}
		var formatPath strings.Builder
		for _, seg := range rpc.PathSegments {
			formatPath.WriteString("/")
			if seg.Type == pathidl.Static {
				formatPath.WriteString(seg.Value)
				continue
			}
			formatPath.WriteString("%v")
		}
		rpc.FormatPath = formatPath.String()
		spec.RPCs[i] = rpc
	}

	return spec, nil
}

// convertConsts converts IDL constants to Go constants
func convertConsts(spec GoSpec, doc httpidl.Document) ([]Const, error) {
	var ret []Const
	for _, c := range doc.Consts {
		typeName, err := goBaseType(c.Type.Name)
		if err != nil {
			return nil, err
		}
		ret = append(ret, Const{
			Const: c,
			Type:  typeName,
		})
	}
	return ret, nil
}

// convertEnums converts IDL enums to Go enums
func convertEnums(spec GoSpec, doc httpidl.Document) ([]Enum, error) {
	var ret []Enum
	for _, e := range doc.Enums {
		ret = append(ret, Enum{e})
	}
	return ret, nil
}

// convertTypes converts IDL struct types to Go struct types
func convertTypes(spec GoSpec, doc httpidl.Document) ([]Type, error) {
	var ret []Type
	for _, t := range doc.Types {
		// Skip generic types (they need instantiation)
		if t.GenericParam != nil {
			continue
		}
		typ, err := convertType(spec, t)
		if err != nil {
			return nil, err
		}
		ret = append(ret, typ)
	}
	return ret, nil
}

// convertType converts an IDL struct type to a Go struct type
func convertType(spec GoSpec, t httpidl.Type) (Type, error) {
	r := Type{Type: t}
	for _, f := range t.Fields {
		fieldName := httpidl.ToPascal(f.Name)

		// Get the type name
		typeName, err := goType(spec, f)
		if err != nil {
			return Type{}, errutil.Explain(nil, "get type name for field %s in type %s error: %w", f.Name, r.Name, err)
		}

		// Determine the category of the field (base, enum, struct, list, map)
		typeKind, err := getTypeKind(spec, typeName)
		if err != nil {
			return Type{}, errutil.Explain(nil, "get type kind for field %s in type %s error: %w", f.Name, r.Name, err)
		}
		if f.Required && IsPointer(typeKind[0]) {
			return Type{}, errutil.Explain(nil, "field %s in type %s is required but has pointer type", f.Name, r.Name)
		}
		if !f.Required && IsNotNullable(typeKind[0]) {
			return Type{}, errutil.Explain(nil, "field %s in type %s is not required but has nullable type", f.Name, r.Name)
		}

		// Add the field to the struct
		field := TypeField{
			TypeField: f,
			Name:      fieldName,
			Type:      typeName,
			TypeKind:  typeKind,
		}
		r.Fields = append(r.Fields, field)
	}
	return r, nil
}

// goBaseType returns the Go type name for a given IDL base type.
func goBaseType(typeName string) (string, error) {
	switch typeName {
	case "string":
		return "string", nil
	case "int":
		return "int64", nil
	case "uint":
		return "uint64", nil
	case "float":
		return "float64", nil
	case "bool":
		return "bool", nil
	default:
		return "", errutil.Explain(nil, "unknown base type: %s", typeName)
	}
}

// goType returns the Go type name for a given IDL type
func goType(spec GoSpec, f httpidl.TypeField) (string, error) {
	if a, ok := httpidl.FindAnnotation(f.Annotations, "go.type"); ok {
		if a.Value == nil {
			return "", errutil.Explain(nil, `annotation "go.type" must have a value`)
		}
		s := strings.Trim(strings.TrimSpace(*a.Value), "\"")
		if s == "" {
			return "", errutil.Explain(nil, `annotation "go.type" must not be empty`)
		}
		return s, nil
	}

	switch typ := f.Type.(type) {
	case httpidl.BytesType:
		return "[]byte", nil
	case httpidl.BaseType:
		s, err := goBaseType(typ.Name)
		if err != nil {
			return "", err
		}
		if f.Required {
			return s, nil
		}
		return "*" + s, nil
	case httpidl.UserType:
		typeName := typ.Name
		_, isEnumType := httpidl.FindEnum(spec.Files, typeName)
		if f.EnumAsString {
			typeName += "AsString"
		}
		if isEnumType && f.Required {
			return typeName, nil
		}
		return "*" + typeName, nil
	default:
		return goTypeDef(spec, typ)
	}
}

// goTypeDef returns the Go type name for a given IDL type.
func goTypeDef(spec GoSpec, t httpidl.TypeDefinition) (string, error) {
	switch typ := t.(type) {
	case httpidl.BaseType:
		return goBaseType(typ.Name)
	case httpidl.UserType:
		if _, ok := httpidl.FindEnum(spec.Files, typ.Name); ok {
			return typ.Name, nil
		}
		return "*" + typ.Name, nil
	case httpidl.ListType:
		itemType, err := goTypeDef(spec, typ.Item)
		if err != nil {
			return "", err
		}
		return "[]" + itemType, nil
	case httpidl.MapType:
		keyType := "string"
		if typ.Key == "int" {
			keyType = "int64"
		}
		valueType, err := goTypeDef(spec, typ.Value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("map[%s]%s", keyType, valueType), nil
	default:
		return "", errutil.Explain(nil, "unknown type: %s", t.Text())
	}
}

// getTypeKind categorizes a Go type for code generation purposes.
func getTypeKind(spec GoSpec, typeName string) ([]TypeKind, error) {
	typeName, pointer := strings.CutPrefix(typeName, "*")

	switch typeName {
	case "[]byte":
		if pointer {
			return nil, errutil.Explain(nil, "binary type can not be pointer")
		}
		return []TypeKind{TypeKindBytes}, nil
	case "bool":
		if pointer {
			return []TypeKind{TypeKindBoolPtr}, nil
		}
		return []TypeKind{TypeKindBool}, nil
	case "int", "int8", "int16", "int32", "int64":
		if pointer {
			return []TypeKind{TypeKindIntPtr}, nil
		}
		return []TypeKind{TypeKindInt}, nil
	case "uint", "uint8", "uint16", "uint32", "uint64":
		if pointer {
			return []TypeKind{TypeKindUintPtr}, nil
		}
		return []TypeKind{TypeKindUint}, nil
	case "float32", "float64":
		if pointer {
			return []TypeKind{TypeKindFloatPtr}, nil
		}
		return []TypeKind{TypeKindFloat}, nil
	case "string":
		if pointer {
			return []TypeKind{TypeKindStringPtr}, nil
		}
		return []TypeKind{TypeKindString}, nil
	case "interface{}", "any":
		if pointer {
			return nil, errutil.Explain(nil, "any type can not be pointer")
		}
		return []TypeKind{TypeKindAny}, nil
	default: // for linter
	}

	switch {
	case strings.HasPrefix(typeName, "[]"):
		if pointer {
			return nil, errutil.Explain(nil, "list type can not be pointer")
		}
		itemType, err := getTypeKind(spec, typeName[2:])
		if err != nil {
			return nil, err
		}
		return append([]TypeKind{TypeKindList}, itemType...), nil
	case strings.HasPrefix(typeName, "map["):
		if pointer {
			return nil, errutil.Explain(nil, "map type can not be pointer")
		}
		itemInex := strings.Index(typeName, "]")
		keyType, err := getTypeKind(spec, typeName[4:itemInex])
		if err != nil {
			return nil, err
		}
		itemType, err := getTypeKind(spec, typeName[itemInex+1:])
		if err != nil {
			return nil, err
		}
		return append([]TypeKind{TypeKindMap, keyType[0]}, itemType...), nil
	default:
		strType, asString := strings.CutSuffix(typeName, "AsString")
		if _, ok := httpidl.FindEnum(spec.Files, strType); ok {
			if asString {
				if pointer {
					return []TypeKind{TypeKindEnumAsStringPtr}, nil
				}
				return []TypeKind{TypeKindEnumAsString}, nil
			}
			if pointer {
				return []TypeKind{TypeKindEnumPtr}, nil
			}
			return []TypeKind{TypeKindEnum}, nil
		}
		if _, ok := httpidl.FindType(spec.Files, typeName); ok {
			if pointer {
				return []TypeKind{TypeKindStructPtr}, nil
			}
		}
		return nil, errutil.Explain(nil, "unknown type: %s", typeName)
	}
}
