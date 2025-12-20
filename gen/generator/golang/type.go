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
			sb.WriteString(fmt.Sprintf(`b, err := json.Marshal(%s)
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
			sb.WriteString(fmt.Sprintf(`b, err := json.Marshal(%s)
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
	sb.WriteString(fmt.Sprintf(`if v, ok := values["%s"]; ok {`, formName))

	switch typeKind[0] {
	case TypeKindList:
		valueType := strings.TrimPrefix(typeName, "[]")
		sb.WriteString(fmt.Sprintf(`for _, s := range v {
				var i %s
				parseErr := json.Unmarshal([]byte(s), &i)
				if parseErr != nil {
					err = errutil.Stack(err, "json decode error: %%w", parseErr)
				}
				%s = append(%s, i)
			}`, valueType, fieldName, fieldName))
		sb.WriteString("}")
		return sb.String()
	case TypeKindMap:
		sb.WriteString(fmt.Sprintf(`if len(v) == 1 {
				parseErr := json.Unmarshal([]byte(v[0]), &%s)
				if parseErr != nil {
					err = errutil.Stack(err, "json decode error: %%w", parseErr)
				}
			} else {
				err = errutil.Stack(err, "invalid value for \"%s\"")
			}
		`, fieldName, formName))
		sb.WriteString("}")
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
		sb.WriteString(fmt.Sprintf(`if parseErr := json.Unmarshal([]byte(v[0]), &%s); parseErr != nil {
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

	sb.WriteString("}")
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
		return "jsonutil.DecodeBool"
	case TypeKindBoolPtr:
		return "jsonutil.DecodeBoolPtr"
	case TypeKindInt:
		return "jsonutil.DecodeInt[" + typeName + "]"
	case TypeKindIntPtr:
		return "jsonutil.DecodeIntPtr[" + strings.TrimPrefix(typeName, "*") + "]"
	case TypeKindUint:
		return "jsonutil.DecodeUint[" + typeName + "]"
	case TypeKindUintPtr:
		return "jsonutil.DecodeUintPtr[" + strings.TrimPrefix(typeName, "*") + "]"
	case TypeKindFloat:
		return "jsonutil.DecodeFloat[" + typeName + "]"
	case TypeKindFloatPtr:
		return "jsonutil.DecodeFloatPtr[" + strings.TrimPrefix(typeName, "*") + "]"
	case TypeKindString:
		return "jsonutil.DecodeString"
	case TypeKindStringPtr:
		return "jsonutil.DecodeStringPtr"
	case TypeKindBytes:
		return "jsonutil.DecodeBytes"
	case TypeKindEnum:
		return "jsonutil.DecodeInt[" + typeName + "]"
	case TypeKindEnumPtr:
		return "jsonutil.DecodeIntPtr[" + strings.TrimPrefix(typeName, "*") + "]"
	case TypeKindEnumAsString:
		return "jsonutil.DecodeInt[" + typeName + "]"
	case TypeKindEnumAsStringPtr:
		return "jsonutil.DecodeIntPtr[" + strings.TrimPrefix(typeName, "*") + "]"
	case TypeKindStructPtr:
		return "jsonutil.DecodeObject(New" + strings.TrimPrefix(typeName, "*") + ")"
	case TypeKindList:
		e := genDecodeJSON(strings.TrimPrefix(typeName, "[]"), typeKind[1:])
		return "jsonutil.DecodeArray(" + e + ")"
	case TypeKindMap:
		s := strings.TrimPrefix(typeName, "map[")
		i := strings.Index(s, "]")
		k, v := s[:i], s[i+1:]
		ks := genDecodeJSON(k, typeKind[1:2])
		vs := genDecodeJSON(v, typeKind[2:])
		return "jsonutil.DecodeMap(" + ks + "," + vs + ")"
	default:
		panic("unsupported type")
	}
}

// typeTmpl is a Go template used to generate Go source code from IDL definitions.
var typeTmpl = template.Must(template.New("type").
	Funcs(map[string]any{
		"formatComments":    formatComments,
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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-spring/gs-http-gen/lib/jsonutil"
	"github.com/lvan100/golib/errutil"
)

var _ = json.Marshal
var _ = strings.Contains
var _ = http.NewServeMux
var _ = strconv.FormatInt

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
		{{range $f := $e.Fields}}
			{{- if $f.Comments.Exists }}
				{{formatComments $f.Comments}}
			{{- end}}
			{{$e.Name}}_{{$f.Name}} {{$e.Name}} = {{$f.Value}}
		{{- end}}
	)

	var (
		{{$e.Name}}_name = map[{{$e.Name}}]string{
			{{- range $f := $e.Fields}}
				{{$f.Value}} : "{{$f.Name}}",
			{{- end}}
		}
		{{$e.Name}}_value = map[string]{{$e.Name}}{
			{{- range $f := $e.Fields}}
				"{{$f.Name}}" : {{$f.Value}},
			{{- end}}
		}
		{{- if eq $e.Name "ErrCode"}}
			{{$e.Name}}_message = map[{{$e.Name}}]string{
				{{- range $f := $e.Fields}}
					{{- if $f.ErrorMessage}}
						{{$f.Value}} : "{{$f.ErrorMessage}}",
					{{- else}}
						{{$f.Value}} : "{{$f.Name}}",
					{{- end}}
				{{- end}}
			}
		{{- end}}
	)

	// OneOf{{$e.Name}} is usually used for validation.
	func OneOf{{$e.Name}}(i {{$e.Name}}) bool {
		_, ok := {{$e.Name}}_name[i]
		return ok
	}

	// OneOf{{$e.Name}}AsString is usually used for validation.
	func OneOf{{$e.Name}}AsString(i {{$e.Name}}AsString) bool {
		_, ok := {{$e.Name}}_name[{{$e.Name}}(i)]
		return ok
	}

	// {{$e.Name}}AsString wraps {{$e.Name}} to encode/decode as a JSON string.
	type {{$e.Name}}AsString {{$e.Name}}

	// MarshalJSON implements custom JSON encoding for the enum as a string.
	func (x {{$e.Name}}AsString) MarshalJSON() ([]byte, error) {
		if s, ok := {{$e.Name}}_name[{{$e.Name}}(x)]; ok {
			return []byte(fmt.Sprintf("\"%s\"", s)), nil
		}
		return nil, errutil.Explain(nil,"invalid {{$e.Name}}: %d", x)
	}

	// UnmarshalJSON implements custom JSON decoding for the enum from a string.
	func (x *{{$e.Name}}AsString) UnmarshalJSON(data []byte) error {
		str := strings.Trim(string(data), "\"")
		if v, ok := {{$e.Name}}_value[str]; ok {
			*x = {{$e.Name}}AsString(v)
			return nil
		}
		return errutil.Explain(nil,"invalid {{$e.Name}} value: %q", str)
	}
{{end}}

{{range $s := .Structs}}
	{{- if $s.Comments.Exists }}
		{{formatComments $s.Comments}}
	{{- end}}
	type {{$s.Name}} struct {
		{{- if $s.Request}}
			{{$s.Name}}Body
		{{- end}}
		{{- range $f := $s.Fields}}
			{{- if or $f.Binding (not $s.Request) }}
				{{- if $f.Comments.Exists }}
					{{formatComments $f.Comments}}
				{{- end}}
				{{$f.Name}} {{$f.Type}} {{$f.FieldTag}}
			{{- end}}
		{{- end}}
	}

	// New{{$s.Name}} ...
	func New{{$s.Name}}() *{{$s.Name}} {
		return &{{$s.Name}}{}
	}

	// DecodeJSON ...
	func (r *{{$s.Name}}) DecodeJSON(d jsonutil.Decoder) error {
		const (
			{{- range $f := $s.Fields}}
				{{- if or $f.Binding (not $s.Request) }}
					hash{{$f.Name}} = {{$f.JSONTag.HashKey}} // HashKey("{{$f.JSONTag.Name}}")
				{{- end}}
			{{- end}}
		)

		var (
			{{- range $f := $s.Fields}}
				{{- if or $f.Binding (not $s.Request) }}
					{{- if $f.Required}}
						has{{$f.Name}} bool
					{{- end}}
				{{- end}}
			{{- end}}
		)

		if err := jsonutil.DecodeObjectBegin(d); err != nil {
			return err
		}

		for {
			if d.PeekKind() == '}' {
				break
			}
			key, err := jsonutil.DecodeKey(d)
			if err != nil {
				return err
			}
			switch jsonutil.HashKey(key) {
			{{- range $f := $s.Fields}}
				{{- if or $f.Binding (not $s.Request) }}
					case hash{{$f.Name}}:
						{{- if $f.Required}}
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

		if err := jsonutil.DecodeObjectEnd(d); err != nil {
			return err
		}

		{{range $f := $s.Fields}}
			{{- if or $f.Binding (not $s.Request) }}
				{{- if $f.Required}}
					if !has{{$f.Name}} {
						return errutil.Explain(nil, "missing required field \"{{$f.JSONTag.Name}}\"")
					}
				{{- end}}
			{{- end}}
		{{- end}}

		return nil
	}

	{{if $s.Request}}
		// QueryForm returns the form values of the object.
		func (x *{{$s.Name}}) QueryForm() (string, error) {
			{{- if $s.QueryCount}}
				m := make(url.Values)
				{{- range $f := $s.Fields}}
					{{- if $f.Binding}}
						{{- if eq $f.Binding.Source "query"}}
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
	
		// Binding extracts non-body values (path, query) from *http.Request.
		func (x *{{$s.Name}}) Bind(r *http.Request) error {
			{{- if $s.BindingCount}}
				values, err := url.ParseQuery(r.URL.RawQuery)
				if err != nil {
					return errutil.Explain(err, "parse query error")
				}
				if len(values) == 0 {
					return nil
				}
				{{- range $f := $s.Fields}}
					{{- if and $f.Binding (eq $f.Binding.Source "query")}}
						{{$fieldName := printf "x.%s" $f.Name}}
						{{- decodeFormValue $fieldName $f.Type $f.TypeKind $f.Binding.Field}}
					{{- end}}
				{{- end}}
				return err
			{{- else}}
				return nil
			{{- end}}
		}
	{{end}}

	{{if $s.OnRequest}}
		// Validate checks field values using generated validation expressions.
		func (x *{{$s.Name}}) Validate() (err error) {
			{{- range $f := $s.Fields}}
				{{- if or $f.Binding (not $s.Request) }}
					{{- if $f.ValidateExpr}}
						{{genValidateExpr $s.Name $f.Name $f.Type $f.ValidateExpr}}
					{{- end}}
					{{- if $f.ValidateNested}}
						{{- $fieldName := printf "x.%s" $f.Name}}
						{{genValidateNested $s.Name $f.Name $fieldName $f.TypeKind 0}}
					{{- end}}
				{{- end}}
			{{- end}}
			{{- if $s.Request}}
				if validateErr := x.{{$s.Name}}Body.Validate(); validateErr != nil {
					err = errutil.Stack(err, "validate failed on \"{{$s.Name}}\": %w", validateErr)
				}
			{{- end}}
			return
		}
	{{end}}

	{{if $s.Request}}
		type {{$s.Name}}Body struct {
			{{- range $f := $s.Fields}}
				{{- if not $f.Binding}}
					{{- if $f.Comments.Exists }}
						{{formatComments $f.Comments}}
					{{- end}}
					{{$f.Name}} {{$f.Type}} {{$f.FieldTag}}
				{{- end}}
			{{- end}}
		}


		// New{{$s.Name}}Body ...
		func New{{$s.Name}}Body() *{{$s.Name}}Body {
			return &{{$s.Name}}Body{}
		}

		// DecodeJSON ...
		func (r *{{$s.Name}}Body) DecodeJSON(d jsonutil.Decoder) error {
			const (
				{{- range $f := $s.Fields}}
					{{- if not $f.Binding}}
						hash{{$f.Name}} = {{$f.JSONTag.HashKey}} // HashKey("{{$f.JSONTag.Name}}")
					{{- end}}
				{{- end}}
			)

			var (
				{{- range $f := $s.Fields}}
					{{- if not $f.Binding}}
						{{- if $f.Required}}
							has{{$f.Name}} bool
						{{- end}}
					{{- end}}
				{{- end}}
			)

			if err := jsonutil.DecodeObjectBegin(d); err != nil {
				return err
			}

			for {
				if d.PeekKind() == '}' {
					break
				}
				key, err := jsonutil.DecodeKey(d)
				if err != nil {
					return err
				}
				switch jsonutil.HashKey(key) {
				{{- range $f := $s.Fields}}
					{{- if not $f.Binding}}
						case hash{{$f.Name}}:
							{{- if $f.Required}}
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

			if err := jsonutil.DecodeObjectEnd(d); err != nil {
				return err
			}

			{{range $f := $s.Fields}}
				{{- if not $f.Binding}}
					{{- if $f.Required}}
						if !has{{$f.Name}} {
							return errutil.Explain(nil, "missing required field \"{{$f.JSONTag.Name}}\"")
						}
					{{- end}}
				{{- end}}
			{{- end}}

			return nil
		}
	{{- end}}

	{{if $s.OnForm}}
		// EncodeForm encodes the object to form data.
		func (x *{{$s.Name}}Body) EncodeForm() (string, error) {
			{{- if $s.BodyCount}}
				m := make(url.Values)
				{{- range $f := $s.Fields}}
					{{- if not $f.Binding}}
						{{$fieldName := printf "x.%s" $f.Name}}
						{{- encodeFormValue $fieldName $f.TypeKind $f.FormTag.Name}}
					{{- end}}
				{{- end}}
				return m.Encode(), nil
			{{- else}}
				return "", nil
			{{- end}}
		}

		// DecodeForm decodes the object from form data.
		func (x *{{$s.Name}}Body) DecodeForm(b []byte) error {
			{{- if $s.BodyCount}}
				values, err := url.ParseQuery(string(b))
				if err != nil {
					return errutil.Explain(err, "parse query error")
				}
				if len(values) == 0 {
					return nil
				}
				{{- range $f := $s.Fields}}
					{{- if not $f.Binding}}
						{{$fieldName := printf "x.%s" $f.Name}}
						{{- decodeFormValue $fieldName $f.Type $f.TypeKind $f.FormTag.Name}}
					{{- end}}
				{{- end}}
				return err
			{{- else}}
				return nil
			{{- end}}
		}
	{{end}}

	{{if $s.Request}}
		// Validate checks field values using generated validation expressions.
		func (x *{{$s.Name}}Body) Validate() (err error) {
			{{- range $f := $s.Fields}}
				{{- if not $f.Binding}}
					{{- if $f.ValidateExpr}}
						{{genValidateExpr $s.Name $f.Name $f.Type $f.ValidateExpr}}
					{{- end}}
					{{- if $f.ValidateNested}}
						{{- $fieldName := printf "x.%s" $f.Name}}
						{{genValidateNested $s.Name $f.Name $fieldName $f.TypeKind 0}}
					{{- end}}
				{{- end}}
			{{- end}}
			return
		}
	{{- end}}

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
		return errutil.Explain(nil, "execute template error: %w", err)
	}
	fileName = fileName[:strings.LastIndex(fileName, ".")] + ".go"
	fileName = filepath.Join(config.OutputDir, fileName)
	return formatFile(fileName, buf.Bytes())
}
