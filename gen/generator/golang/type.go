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
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-spring/gs-http-gen/gen/generator"
	"github.com/go-spring/gs-http-gen/lib/httpidl"
	"github.com/go-spring/gs-http-gen/lib/validate"
	"github.com/lvan100/golib/errutil"
)

// formatFile formats Go source code using `go format`
// and writes the formatted code to the given file.
func formatFile(fileName string, b []byte) error {
	b, err := format.Source(b)
	if err != nil {
		return errutil.Explain(nil, "format source for file %s error: %w", fileName, err)
	}
	err = os.WriteFile(fileName, b, os.ModePerm)
	if err != nil {
		return errutil.Explain(nil, "write file %s error: %w", fileName, err)
	}
	return nil
}

// formatComments converts a tidl.Comments into Go comments.
func formatComments(c httpidl.Comments) string {
	var lines []string
	for _, s := range c.Above {
		lines = append(lines, s.Text...)
	}
	if c.Right != nil {
		lines = append(lines, c.Right.Text...)
	}
	return strings.Join(lines, "\n")
}

// genDefaultValue generates Go code to generate a default value for a field.
func genDefaultValue(typeKind []TypeKind, defaultValue string) string {
	switch typeKind[0] {
	case TypeKindInt:
		_, err := strconv.ParseInt(defaultValue, 10, 64)
		if err != nil {
			panic("parse error")
		}
		return defaultValue
	case TypeKindUint:
		_, err := strconv.ParseUint(defaultValue, 10, 64)
		if err != nil {
			panic("parse error")
		}
		return defaultValue
	case TypeKindString:
		return strconv.Quote(defaultValue)
	default:
		panic("unsupported type")
	}
}

// decodePathValue generates Go code to decode a field value from path parameter.
func decodePathValue(fieldName string, typeKind []TypeKind, pathName string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`if s := r.PathValue("%s"); s == "" {
		err = errutil.Stack(err, "required field \"%s\" is missing")
	} else {
	`, pathName, pathName))

	switch typeKind[0] {
	case TypeKindInt:
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseInt(s, 10, 64); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				%s = i
			}`, pathName, fieldName))
	case TypeKindUint:
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseUint(s, 10, 64); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				%s = i
			}`, pathName, fieldName))
	case TypeKindString:
		sb.WriteString(fmt.Sprintf(`%s = s`, fieldName))
	default:
		panic("unsupported type")
	}
	sb.WriteString(fmt.Sprintf("\n}"))
	return sb.String()
}

// encodeFormValue generates Go code to encode a field value to form data.
func encodeFormValue(fieldName string, typeKind []TypeKind, formName string) string {
	var sb strings.Builder
	if IsPointer(typeKind[0]) {
		sb.WriteString(fmt.Sprintf("if %s != nil {\n", fieldName))
		switch typeKind[0] {
		case TypeKindBoolPtr:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", strconv.FormatBool(%s))`, formName, "*"+fieldName))
		case TypeKindIntPtr:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", strconv.FormatInt(int64(%s), 10))`, formName, "*"+fieldName))
		case TypeKindUintPtr:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", strconv.FormatUint(uint64(%s), 10))`, formName, "*"+fieldName))
		case TypeKindFloatPtr:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", strconv.FormatFloat(float64(%s), 'f', -1, 64))`, formName, "*"+fieldName))
		case TypeKindStringPtr:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", %s)`, formName, "*"+fieldName))
		case TypeKindBytes:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", %s)`, formName, "*"+fieldName)) // todo
		case TypeKindEnumPtr:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", strconv.FormatInt(int64(%s), 10))`, formName, "*"+fieldName))
		case TypeKindEnumAsStringPtr:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", string(%s))`, formName, "*"+fieldName))
		case TypeKindStructPtr:
			sb.WriteString(fmt.Sprintf(`b, err := jsonflow.Marshal(%s)
				if err != nil {
					return "", err
				}
				m.Add("%s", string(b))`, "*"+fieldName, formName))
		default:
			panic("unsupported type")
		}
		sb.WriteString(fmt.Sprintf("\n}"))
	} else {
		switch typeKind[0] {
		case TypeKindBool:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", strconv.FormatBool(%s))`, formName, fieldName))
		case TypeKindInt:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", strconv.FormatInt(int64(%s), 10))`, formName, fieldName))
		case TypeKindUint:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", strconv.FormatUint(uint64(%s), 10))`, formName, fieldName))
		case TypeKindFloat:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", strconv.FormatFloat(float64(%s), 'f', -1, 64))`, formName, fieldName))
		case TypeKindString:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", %s)`, formName, fieldName))
		case TypeKindEnum:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", strconv.FormatInt(int64(%s), 10))`, formName, fieldName))
		case TypeKindEnumAsString:
			sb.WriteString(fmt.Sprintf(`m.Add("%s", string(%s))`, formName, fieldName))
		case TypeKindMap:
			sb.WriteString(fmt.Sprintf(`b, err := jsonflow.Marshal(%s)
				if err != nil {
					return "", err
				}
				m.Add("%s", string(b))`, fieldName, formName))
		case TypeKindList:
			sb.WriteString(fmt.Sprintf("for i := range len(%s) {\n", fieldName))
			sb.WriteString(encodeFormValue(fieldName+"[i]", typeKind[1:], formName))
			sb.WriteString(fmt.Sprintf("\n}"))
		default:
			panic("unsupported type")
		}
	}
	return sb.String()
}

// decodeFormValue generates Go code to decode a field value from form data.
func decodeFormValue(fieldName string, typeName string, typeKind []TypeKind, formName string) string {
	var sb strings.Builder

	switch typeKind[0] {
	case TypeKindList:
		valueType := strings.TrimPrefix(typeName, "[]")
		sb.WriteString(fmt.Sprintf(`for _, s := range v {
				var i %s
				if parseErr = jsonflow.Unmarshal([]byte(s), &i); parseErr != nil {
					err = errutil.Stack(err, "json decode error: %%w", parseErr)
				} else {
					%s = append(%s, i)
				}
			}`, valueType, fieldName, fieldName))
		return sb.String()
	case TypeKindMap:
		sb.WriteString(fmt.Sprintf(`if len(v) == 1 {
				parseErr := jsonflow.Unmarshal([]byte(v[0]), &%s)
				if parseErr != nil {
					err = errutil.Stack(err, "json decode error: %%w", parseErr)
				}
			} else {
				err = errutil.Stack(err, "invalid value for \"%s\"")
			}
		`, fieldName, formName))
		return sb.String()
	default: // for linter
	}

	sb.WriteString(fmt.Sprintf(`if len(v) == 1 {`))
	switch typeKind[0] {
	case TypeKindBool:
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseBool(v[0]); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				%s = i
			}`, formName, fieldName))
	case TypeKindBoolPtr:
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseBool(v[0]); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				%s = &i
			}`, formName, fieldName))
	case TypeKindInt:
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseInt(v[0], 10, 64); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				%s = i
			}`, formName, fieldName))
	case TypeKindIntPtr:
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseInt(v[0], 10, 64); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				%s = &i
			}`, formName, fieldName))
	case TypeKindUint:
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseUint(v[0], 10, 64); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				%s = i
			}`, formName, fieldName))
	case TypeKindUintPtr:
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseUint(v[0], 10, 64); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				%s = &i
			}`, formName, fieldName))
	case TypeKindFloat:
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseFloat(v[0], 64); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				%s = i
			}`, formName, fieldName))
	case TypeKindFloatPtr:
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseFloat(v[0], 64); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				%s = &i
			}`, formName, fieldName))
	case TypeKindString:
		sb.WriteString(fmt.Sprintf(`%s = v[0]`, fieldName))
	case TypeKindStringPtr:
		sb.WriteString(fmt.Sprintf(`%s = &v[0]`, fieldName))
	case TypeKindStructPtr:
		sb.WriteString(fmt.Sprintf(`if parseErr := jsonflow.Unmarshal([]byte(v[0]), &%s); parseErr != nil {
				err = errutil.Stack(err, "json decode error: %%w", parseErr)
			}`, fieldName))
	case TypeKindBytes: // todo
	case TypeKindEnum:
		valueType := strings.TrimPrefix(typeName, "*")
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseInt(v[0], 10, 64); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				if e := %s(i); !OneOf%s(e) {
					err = errutil.Stack(err, "invalid value for \"%s\"")
				} else{
					%s = e
				}
			}`, formName, valueType, valueType, formName, fieldName))
	case TypeKindEnumPtr:
		valueType := strings.TrimPrefix(typeName, "*")
		sb.WriteString(fmt.Sprintf(`if i, parseErr := strconv.ParseInt(v[0], 10, 64); parseErr != nil {
				err = errutil.Stack(err, "parse \"%s\" error: %%w", parseErr)
			} else {
				if e := %s(i); !OneOf%s(e) {
					err = errutil.Stack(err, "invalid value for \"%s\"")
				} else{
					%s = &e
				}
			}`, formName, valueType, valueType, formName, fieldName))
	case TypeKindEnumAsString:
	case TypeKindEnumAsStringPtr:
	default:
		panic("unsupported type")
	}

	sb.WriteString(fmt.Sprintf(`} else {
			err = errutil.Stack(err, "invalid value for \"%s\"")
		}`, formName))
	return sb.String()
}

// genValidateExpr generates the Go code for a validation expression
func genValidateExpr(receiverType, fieldName, fieldType string, expr validate.Expr) (string, error) {
	receiverType = strings.TrimSuffix(receiverType, "Body") // todo

	// 对于结构体而言，只应当验证字段非空，其内部字段的验证应当由自己完成
	fieldType, pointer := strings.CutPrefix(fieldType, "*")
	dollar := "x." + fieldName
	if pointer {
		dollar = "*" + dollar
	}

	// Generate the Go expression for validation
	str, err := compileValidateExpr(dollar, fieldType, expr)
	if err != nil {
		return "", errutil.Explain(err, `failed to generate validate expression for %s.%s`, receiverType, fieldName)
	}

	// Wrap in an if statement returning an error on failure
	str = fmt.Sprintf(`if !(%s) {
		err = errutil.Stack(err,"validate failed on \"%s.%s\"")
	}`, str, receiverType, fieldName)

	if pointer {
		str = fmt.Sprintf(`if x.%s != nil { %s }`, fieldName, str)
	}
	return str, nil
}

// compileValidateExpr recursively generates Go code for a validation expression
func compileValidateExpr(fieldName, fieldType string, expr validate.Expr) (string, error) {
	switch x := expr.(type) {
	case validate.BinaryExpr:
		left, err := compileValidateExpr(fieldName, fieldType, x.Left)
		if err != nil {
			return "", err
		}
		right, err := compileValidateExpr(fieldName, fieldType, x.Right)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s %s %s", left, x.Op, right), nil

	case validate.UnaryExpr:
		str, err := compileValidateExpr(fieldName, fieldType, x.Expr)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s%s", x.Op, str), nil

	case *validate.FuncCall:
		if len(x.Args) == 0 {
			return x.Name + "()", nil
		}
		var args []string
		for _, arg := range x.Args {
			str, err := compileValidateExpr(fieldName, fieldType, arg)
			if err != nil {
				return "", err
			}
			args = append(args, str)
		}
		return fmt.Sprintf("%s(%s)", x.Name, strings.Join(args, ", ")), nil

	case *validate.InnerExpr:
		str, err := compileValidateExpr(fieldName, fieldType, x.Expr)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s)", str), nil

	case validate.PrimaryExpr:
		if x.Inner != nil {
			return compileValidateExpr(fieldName, fieldType, x.Inner)
		}
		if x.Call != nil {
			return compileValidateExpr(fieldName, fieldType, x.Call)
		}
		if x.Value == "$" {
			return fieldName, nil
		}
		return x.Value, nil

	default:
		return "", errutil.Explain(nil, "unknown expression type: %s", x.Text())
	}
}

// genValidateNested generates the nested validation code for a Go struct field.
func genValidateNested(receiverType, fieldName string, itemName string, typeKind []TypeKind, depth int) string {
	receiverType = strings.TrimSuffix(receiverType, "Body") // todo
	childName := fmt.Sprintf("v%d", depth)
	switch typeKind[0] {
	case TypeKindList, TypeKindMap:
		str := genValidateNested(receiverType, fieldName, childName, typeKind[1:], depth+1)
		if str == "" {
			return ""
		}
		str = fmt.Sprintf(`for _, %s := range %s {
				%s
			}`, childName, itemName, str)
		return str
	case TypeKindStructPtr:
		str := fmt.Sprintf(`if %s != nil {
				if validateErr := %s.Validate(); validateErr != nil {
					err = errutil.Stack(err, "validate failed on \"%s.%s\": %%w", validateErr)
				}
			}`, itemName, itemName, receiverType, fieldName)
		return str
	default:
		return ""
	}
}

// genDecodeJSON generates the JSON decoding code for a Go struct field.
func genDecodeJSON(typeName string, typeKind []TypeKind) string {
	switch typeKind[0] {
	case TypeKindBool:
		return "jsonflow.DecodeBool"
	case TypeKindBoolPtr:
		return "jsonflow.DecodeBoolPtr"
	case TypeKindInt:
		return "jsonflow.DecodeInt[" + typeName + "]"
	case TypeKindIntPtr:
		return "jsonflow.DecodeIntPtr[" + strings.TrimPrefix(typeName, "*") + "]"
	case TypeKindUint:
		return "jsonflow.DecodeUint[" + typeName + "]"
	case TypeKindUintPtr:
		return "jsonflow.DecodeUintPtr[" + strings.TrimPrefix(typeName, "*") + "]"
	case TypeKindFloat:
		return "jsonflow.DecodeFloat[" + typeName + "]"
	case TypeKindFloatPtr:
		return "jsonflow.DecodeFloatPtr[" + strings.TrimPrefix(typeName, "*") + "]"
	case TypeKindString:
		return "jsonflow.DecodeString"
	case TypeKindStringPtr:
		return "jsonflow.DecodeStringPtr"
	case TypeKindBytes:
		return "jsonflow.DecodeBytes"
	case TypeKindEnum:
		return "jsonflow.DecodeInt[" + typeName + "]"
	case TypeKindEnumPtr:
		return "jsonflow.DecodeIntPtr[" + strings.TrimPrefix(typeName, "*") + "]"
	case TypeKindEnumAsString:
		return "jsonflow.DecodeAny[" + typeName + "]"
	case TypeKindEnumAsStringPtr:
		return "jsonflow.DecodeAny[" + typeName + "]"
	case TypeKindStructPtr:
		return "jsonflow.DecodeObject(New" + strings.TrimPrefix(typeName, "*") + ")"
	case TypeKindList:
		e := genDecodeJSON(strings.TrimPrefix(typeName, "[]"), typeKind[1:])
		return "jsonflow.DecodeArray(" + e + ")"
	case TypeKindMap:
		s := strings.TrimPrefix(typeName, "map[")
		i := strings.Index(s, "]")
		k, v := s[:i], s[i+1:]
		ks := genDecodeJSON(k, typeKind[1:2])
		vs := genDecodeJSON(v, typeKind[2:])
		return "jsonflow.DecodeMap(" + ks + "," + vs + ")"
	default:
		panic("unsupported type")
	}
}

// typeTmpl is a Go template used to generate Go source code from IDL definitions.
var typeTmpl = template.Must(template.New("type").
	Funcs(map[string]any{
		"formatComments":    formatComments,
		"genDefaultValue":   genDefaultValue,
		"decodePathValue":   decodePathValue,
		"encodeFormValue":   encodeFormValue,
		"decodeFormValue":   decodeFormValue,
		"genValidateExpr":   genValidateExpr,
		"genValidateNested": genValidateNested,
		"genDecodeJSON":     genDecodeJSON,
	}).
	Parse(`
// Code generated by gs-http-gen compiler. DO NOT EDIT.

package {{.Package}}

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/lvan100/golib/errutil"
	"github.com/lvan100/golib/hashutil"
	"github.com/lvan100/golib/jsonflow"
)

var _ = strings.Trim
var _ = strconv.Itoa
var _ = http.StatusOK

{{ range $c := .Consts }}
	{{- if $c.Comments.Exists }}
		{{formatComments $c.Comments}}
	{{- end}}
	const {{$c.Name}} {{$c.Type}} = {{$c.Value}}
{{end}}

{{ range $e := .Enums }}
	{{- if $e.Comments.Exists }}
		{{formatComments $e.Comments}}
	{{- end}}
	type {{$e.Name}} int32

	const (
		{{ range $f := $e.Fields }}
			{{- if $f.Comments.Exists }}
				{{formatComments $f.Comments}}
			{{- end}}
			{{$e.Name}}_{{$f.Name}} {{$e.Name}} = {{$f.Value}}
		{{- end}}
	)

	var (
		{{$e.Name}}_name = map[{{$e.Name}}]string{
			{{- range $f := $e.Fields }}
				{{$f.Value}} : "{{$f.Name}}",
			{{- end}}
		}
		{{$e.Name}}_value = map[string]{{$e.Name}}{
			{{- range $f := $e.Fields }}
				"{{$f.Name}}" : {{$f.Value}},
			{{- end}}
		}
		{{- if $e.KindError }} {{- /* only for error */}}
			{{$e.Name}}_message = map[{{$e.Name}}]string{
				{{- range $f := $e.Fields }}
					{{- if $f.ErrorMessage }}
						{{$f.Value}} : "{{$f.ErrorMessage}}",
					{{- else}}
						{{$f.Value}} : "{{$f.Name}}",
					{{- end}}
				{{- end}}
			}
		{{- end}}
	)

	// OneOf{{$e.Name}} reports whether it's a valid {{$e.Name}}.
	func OneOf{{$e.Name}}(i {{$e.Name}}) bool {
		_, ok := {{$e.Name}}_name[i]
		return ok
	}

	// OneOf{{$e.Name}}AsString reports whether it's a valid {{$e.Name}}AsString.
	func OneOf{{$e.Name}}AsString(i {{$e.Name}}AsString) bool {
		_, ok := {{$e.Name}}_name[{{$e.Name}}(i)]
		return ok
	}

	// {{$e.Name}}AsString wraps {{$e.Name}} to encode/decode as a JSON string.
	type {{$e.Name}}AsString {{$e.Name}}

	// MarshalJSON encodes the enum value as its string name.
	func (x {{$e.Name}}AsString) MarshalJSON() ([]byte, error) {
		if s, ok := {{$e.Name}}_name[{{$e.Name}}(x)]; ok {
			return []byte(fmt.Sprintf("\"%s\"", s)), nil
		}
		return nil, errutil.Explain(nil,"invalid {{$e.Name}}AsString: %d", x)
	}

	// UnmarshalJSON decodes the enum value from its string name.
	func (x *{{$e.Name}}AsString) UnmarshalJSON(data []byte) error {
		str := strings.Trim(string(data), "\"")
		if v, ok := {{$e.Name}}_value[str]; ok {
			*x = {{$e.Name}}AsString(v)
			return nil
		}
		return errutil.Explain(nil,"invalid {{$e.Name}}AsString value: %q", str)
	}
{{end}}

{{ range $s := .Structs }}
	{{- if not $s.Request }}
		{{- if $s.Comments.Exists }}
			{{formatComments $s.Comments}}
		{{- end}}
		type {{$s.Name}} struct {
			{{- range $f := $s.Fields }}
				{{- if $f.Comments.Exists }}
					{{formatComments $f.Comments}}
				{{- end}}
				{{$f.Name}} {{$f.Type}} {{$f.FieldTag}}
			{{- end}}
		}

		// New{{$s.Name}} creates a new {{$s.Name}} instance
		// and initializes fields with default values.
		func New{{$s.Name}}() *{{$s.Name}} {
			return &{{$s.Name}}{
				{{- range $f := $s.Fields }}
					{{- if $f.CompatDefault}}
						{{$f.Name}}: {{genDefaultValue $f.TypeKind $f.CompatDefault}},
					{{- end}}
				{{- end}}
			}
		}

		// DecodeJSON decodes a JSON object into {{$s.Name}} using a hash-based
		// field dispatch mechanism for high-performance parsing.
		func (r *{{$s.Name}}) DecodeJSON(d jsonflow.Decoder) (err error) {
			{{- if $s.FieldCount }}
				const (
					{{- range $f := $s.Fields }}
						hash{{$f.Name}} = {{$f.JSONTag.HashKey}} // HashKey("{{$f.JSONTag.Name}}")
					{{- end}}
				)
			{{- end}}

			{{- $HasRequired := false }}
			{{- range $f := $s.Fields }}
				{{- if and $f.Required (not $f.CompatDefault) }}
					{{ $HasRequired = true }}
				{{- end}}
			{{- end}}

			{{- if $HasRequired }}
				var (
					{{- range $f := $s.Fields }}
						{{- if and $f.Required (not $f.CompatDefault) }}
							has{{$f.Name}} bool
						{{- end}}
					{{- end}}
				)
			{{- end}}

			if err = jsonflow.DecodeObjectBegin(d); err != nil {
				return err
			}

			for {
				if d.PeekKind() == '}' {
					break
				}

				var key string
				key, err = jsonflow.DecodeString(d)
				if err != nil {
					return err
				}

				switch hashutil.FNV1a64(key) {
				{{- range $f := $s.Fields }}
					case hash{{$f.Name}}:
						{{- if and $f.Required (not $f.CompatDefault) }}
							has{{$f.Name}} = true
						{{- end}}
						if r.{{$f.Name}}, err = {{genDecodeJSON $f.Type $f.TypeKind}}(d); err != nil {
							return err
						}
				{{- end}}
				default:
					if err = d.SkipValue(); err != nil {
						return err
					}
				}
			}

			if err = jsonflow.DecodeObjectEnd(d); err != nil {
				return err
			}

			{{- if $HasRequired }}
				{{ range $f := $s.Fields }}
					{{- if and $f.Required (not $f.CompatDefault) }}
						if !has{{$f.Name}} {
							err = errutil.Stack(err, "missing required field \"{{$f.JSONTag.Name}}\"")
						}
					{{- end}}
				{{- end}}
			{{- end}}
			return 
		}

		{{- $Validate := false}}
		{{- range $f := $s.Fields }}
			{{- if $f.ValidateExpr }}
				{{- $Validate = true}}
			{{- end}}
			{{- if $f.ValidateNested }}
				{{- $Validate = true}}
			{{- end}}
		{{- end}}

		{{- if $Validate}}
			// Validate checks field values using generated validation expressions.
			func (x *{{$s.Name}}) Validate() (err error) {
				{{- range $f := $s.Fields }}
					{{- if $f.ValidateExpr }}
						{{genValidateExpr $s.Name $f.Name $f.Type $f.ValidateExpr}}
					{{- end}}
					{{- if $f.ValidateNested }}
						{{- $fieldName := printf "x.%s" $f.Name}}
						{{genValidateNested $s.Name $f.Name $fieldName $f.TypeKind 0}}
					{{- end}}
				{{- end}}
				return
			}
		{{- end}}
	{{- end}} {{- /* end of struct (not request) */}}

	{{- if $s.Request }}
		{{- if $s.Comments.Exists }}
			{{formatComments $s.Comments}}
		{{- end}}
		type {{$s.Name}} struct {
			{{$s.Name}}Body
			{{- range $f := $s.Fields }}
				{{- if $f.Binding }}
					{{- if $f.Comments.Exists }}
						{{formatComments $f.Comments}}
					{{- end}}
					{{$f.Name}} {{$f.Type}} {{$f.FieldTag}}
				{{- end}}
			{{- end}}
		}

		// New{{$s.Name}} creates a new {{$s.Name}} instance
		// and initializes fields with default values.
		func New{{$s.Name}}() *{{$s.Name}} {
			return &{{$s.Name}}{
				{{- range $f := $s.Fields }}
					{{- if $f.Binding }}
						{{- if $f.CompatDefault}}
							{{$f.Name}}: {{genDefaultValue $f.TypeKind $f.CompatDefault}},
						{{- end}}
					{{- end}}
				{{- end}}
			}
		}

		// QueryForm encodes query-bound fields into URL-encoded form data.
		func (x *{{$s.Name}}) QueryForm() (string, error) {
			{{- if $s.QueryCount }}
				m := make(url.Values)
				{{- range $f := $s.Fields }}
					{{- if $f.Binding }}
						{{- if eq $f.Binding.Source "query" }}
							{{$fieldName := printf "x.%s" $f.Name}}
							{{- encodeFormValue $fieldName $f.TypeKind $f.Binding.Field}}
						{{- end}}
					{{- end}}
				{{- end}}
				return m.Encode(), nil
			{{- else}}
				return "", nil
			{{- end}}
		}

		// Bind extracts path and query parameters from the HTTP request
		// and assigns them to the corresponding struct fields.
		func (x *{{$s.Name}}) Bind(r *http.Request) (err error) {
			{{- if $s.BindingCount }}
				{{- range $f := $s.Fields }}
					{{- if and $f.Binding (eq $f.Binding.Source "path") }}
						{{- $fieldName := printf "x.%s" $f.Name}}
						{{decodePathValue $fieldName $f.TypeKind $f.Binding.Field}}
					{{- end}}
				{{- end}}

				{{- if $s.QueryCount }}
					values, parseErr := url.ParseQuery(r.URL.RawQuery)
					if parseErr != nil {
						err = errutil.Explain(err, "parse query error: %w", parseErr)
						return
					}

					{{- $HasRequired := false }}
					{{- range $f := $s.Fields }}
						{{- if and $f.Binding (eq $f.Binding.Source "query") }}
							{{- if and $f.Required (not $f.CompatDefault) }}
								{{ $HasRequired = true }}
							{{- end}}
						{{- end}}
					{{- end}}

					{{- if $HasRequired }}
						var (
							{{- range $f := $s.Fields }}
								{{- if and $f.Binding (eq $f.Binding.Source "query") }}
									{{- if and $f.Required (not $f.CompatDefault) }}
										has{{$f.Name}} bool
									{{- end}}
								{{- end}}
							{{- end}}
						)
					{{end}}

					{{range $f := $s.Fields }}
						{{- if and $f.Binding (eq $f.Binding.Source "query") }}
							if v, ok := values["{{$f.Binding.Field}}"]; ok {
								{{- if and $f.Required (not $f.CompatDefault) }}
									has{{$f.Name}} = true
								{{- end}}
								{{- $fieldName := printf "x.%s" $f.Name}}
								{{decodeFormValue $fieldName $f.Type $f.TypeKind $f.Binding.Field}}
							}
						{{- end}}
					{{- end}}

					{{- if $HasRequired }}
						{{- range $f := $s.Fields }}
							{{- if and $f.Binding (eq $f.Binding.Source "query") }}
								{{- if and $f.Required (not $f.CompatDefault) }}
									if !has{{$f.Name}} {
										err = errutil.Explain(err, "missing required field \"{{$f.Binding.Field}}\"")
									}
								{{- end}}
							{{- end}}
						{{- end}}
					{{- end}}
				{{- end}}
			{{- end}}
			return
		}

		// Validate validates both bound parameters and request body fields.
		func (x *{{$s.Name}}) Validate() (err error) {
			{{- range $f := $s.Fields }}
				{{- if $f.Binding }}
					{{- if $f.ValidateExpr }}
						{{genValidateExpr $s.Name $f.Name $f.Type $f.ValidateExpr}}
					{{- end}}
					{{- if $f.ValidateNested }}
						{{- $fieldName := printf "x.%s" $f.Name}}
						{{genValidateNested $s.Name $f.Name $fieldName $f.TypeKind 0}}
					{{- end}}
				{{- end}}
			{{- end}}
			if validateErr := x.{{$s.Name}}Body.Validate(); validateErr != nil {
				err = errutil.Stack(err, "validate failed on \"{{$s.Name}}\": %w", validateErr)
			}
			return
		}

		// {{$s.Name}}Body represents the request body payload,
		// excluding path and query parameters.
		type {{$s.Name}}Body struct {
			{{- range $f := $s.Fields }}
				{{- if not $f.Binding }}
					{{- if $f.Comments.Exists }}
						{{formatComments $f.Comments}}
					{{- end}}
					{{$f.Name}} {{$f.Type}} {{$f.FieldTag}}
				{{- end}}
			{{- end}}
		}

		// New{{$s.Name}}Body creates a new {{$s.Name}}Body instance
		// and initializes fields with default values.
		func New{{$s.Name}}Body() *{{$s.Name}}Body {
			return &{{$s.Name}}Body{
				{{- range $f := $s.Fields }}
					{{- if not $f.Binding }}
						{{- if $f.CompatDefault}}
							{{$f.Name}}: {{genDefaultValue $f.TypeKind $f.CompatDefault}},
						{{- end}}
					{{- end}}
				{{- end}}
			}
		}

		{{- if $s.FormEncoded }}
			// EncodeForm encodes the request body as application/x-www-form-urlencoded data.
			func (x *{{$s.Name}}Body) EncodeForm() (string, error) {
				{{- if $s.BodyCount }}
					m := make(url.Values)
					{{- range $f := $s.Fields }}
						{{- if not $f.Binding }}
							{{$fieldName := printf "x.%s" $f.Name}}
							{{- encodeFormValue $fieldName $f.TypeKind $f.FormTag.Name}}
						{{- end}}
					{{- end}}
					return m.Encode(), nil
				{{- else}}
					return "", nil
				{{- end}}
			}

			// DecodeForm decodes application/x-www-form-urlencoded data into the request body.
			func (x *{{$s.Name}}Body) DecodeForm(b []byte) (err error) {
				{{- if $s.BodyCount }}
					values, parseErr := url.ParseQuery(string(b))
					if parseErr != nil {
						err = errutil.Explain(err, "parse query error: %w", parseErr)
						return
					}

					{{- $HasRequired := false }}
					{{- range $f := $s.Fields }}
						{{- if not $f.Binding }}
							{{- if and $f.Required (not $f.CompatDefault) }}
								{{ $HasRequired = true }}
							{{- end}}
						{{- end}}
					{{- end}}

					{{- if $HasRequired }}
						var (
							{{- range $f := $s.Fields }}
								{{- if not $f.Binding }}
									{{- if and $f.Required (not $f.CompatDefault) }}
										has{{$f.Name}} bool
									{{- end}}
								{{- end}}
							{{- end}}
						)
					{{end}}

					{{range $f := $s.Fields }}
						{{- if not $f.Binding }}
							if v, ok := values["{{$f.FormTag.Name}}"]; ok {
								{{- if and $f.Required (not $f.CompatDefault) }}
									has{{$f.Name}} = true
								{{- end}}
								{{- $fieldName := printf "x.%s" $f.Name}}
								{{decodeFormValue $fieldName $f.Type $f.TypeKind $f.FormTag.Name}}
							}
						{{- end}}
					{{- end}}

					{{- if $HasRequired }}
						{{- range $f := $s.Fields }}
							{{- if not $f.Binding }}
								{{- if and $f.Required (not $f.CompatDefault) }}
									if !has{{$f.Name}} {
										err = errutil.Explain(err, "missing required field \"{{$f.FormTag.Name}}\"")
									}
								{{- end}}
							{{- end}}
						{{- end}}
					{{- end}}
				{{- end}}

				return
			}
		{{- end}} {{- /* end of form encoded */}}

		{{- if $s.JSONEncoded }}
			// DecodeJSON decodes a JSON object into {{$s.Name}}Body using a hash-based
			// field dispatch mechanism for high-performance parsing.
			func (r *{{$s.Name}}Body) DecodeJSON(d jsonflow.Decoder) (err error) {
				{{- if $s.BodyCount }}
					const (
						{{- range $f := $s.Fields }}
							{{- if not $f.Binding }}
								hash{{$f.Name}} = {{$f.JSONTag.HashKey}} // HashKey("{{$f.JSONTag.Name}}")
							{{- end}}
						{{- end}}
					)
				{{- end}}

				{{- $HasRequired := false }}
				{{- range $f := $s.Fields }}
					{{- if not $f.Binding }}
						{{- if and $f.Required (not $f.CompatDefault) }}
							{{ $HasRequired = true }}
						{{- end}}
					{{- end}}
				{{- end}}

				{{- if $HasRequired }}
					var (
						{{- range $f := $s.Fields }}
							{{- if not $f.Binding }}
								{{- if and $f.Required (not $f.CompatDefault) }}
									has{{$f.Name}} bool
								{{- end}}
							{{- end}}
						{{- end}}
					)
				{{- end}}

				if err = jsonflow.DecodeObjectBegin(d); err != nil {
					return err
				}

				for {
					if d.PeekKind() == '}' {
						break
					}

					var key string
					key, err = jsonflow.DecodeString(d)
					if err != nil {
						return err
					}

					switch hashutil.FNV1a64(key) {
					{{- range $f := $s.Fields }}
						{{- if not $f.Binding }}
							case hash{{$f.Name}}:
								{{- if and $f.Required (not $f.CompatDefault) }}
									has{{$f.Name}} = true
								{{- end}}
								if r.{{$f.Name}}, err = {{genDecodeJSON $f.Type $f.TypeKind}}(d); err != nil {
									return err
								}
						{{- end}}
					{{- end}}
					default:
						if err = d.SkipValue(); err != nil {
							return err
						}
					}
				}

				if err = jsonflow.DecodeObjectEnd(d); err != nil {
					return err
				}

				{{- if $HasRequired }}
					{{ range $f := $s.Fields }}
						{{- if not $f.Binding }}
							{{- if and $f.Required (not $f.CompatDefault) }}
								if !has{{$f.Name}} {
									return errutil.Stack(err, "missing required field \"{{$f.JSONTag.Name}}\"")
								}
							{{- end}}
						{{- end}}
					{{- end}}
				{{- end}}
				return 
			}
		{{end}} {{- /* end of json encoded */}}

		// Validate checks field values using generated validation expressions.
		func (x *{{$s.Name}}Body) Validate() (err error) {
			{{- range $f := $s.Fields }}
				{{- if not $f.Binding }}
					{{- if $f.ValidateExpr }}
						{{genValidateExpr $s.Name $f.Name $f.Type $f.ValidateExpr}}
					{{- end}}
					{{- if $f.ValidateNested }}
						{{- $fieldName := printf "x.%s" $f.Name}}
						{{genValidateNested $s.Name $f.Name $fieldName $f.TypeKind 0}}
					{{- end}}
				{{- end}}
			{{- end}}
			return
		}

	{{end}} {{- /* end of struct (request) */}}
{{end}}
`))

// genType generates a Go source file corresponding to the IDL file.
// It includes constants, enums, and struct types.
func (g *Generator) genType(config *generator.Config, fileName string, spec GoSpec) error {
	buf := &bytes.Buffer{}
	err := typeTmpl.Execute(buf, map[string]any{
		"Package": config.GoPackage,
		"Consts":  spec.Consts[fileName],
		"Enums":   spec.Enums[fileName],
		"Structs": spec.Types[fileName],
	})
	if err != nil {
		return errutil.Explain(nil, "execute type template error: %w", err)
	}
	fileName = fileName[:strings.LastIndex(fileName, ".")] + ".go"
	fileName = filepath.Join(config.OutputDir, fileName)
	return formatFile(fileName, buf.Bytes())
}
