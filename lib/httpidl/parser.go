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
	"bufio"
	"bytes"
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/go-spring/gs-http-gen/lib/pathidl"
	"github.com/go-spring/gs-http-gen/lib/validate"
	"github.com/lvan100/golib/errutil"
)

// BuiltinFuncs is a set of built-in validation functions
var BuiltinFuncs = map[string]struct{}{
	"len": {},
}

// Project represents a collection of IDL files and their associated meta-information.
type Project struct {
	Meta  *MetaInfo
	Files map[string]Document
	Funcs map[string]ValidateFunc
}

// RequestMeta represents the metadata of a request type.
type RequestMeta struct {
	OnForm bool
}

// ValidateFunc represents a validate function.
type ValidateFunc struct {
	Name string
	Type string
}

// ParseDir scans the specified directory for IDL files (*.idl) and a meta.json file.
// It parses each file into a Document structure and validates cross-file type references.
func ParseDir(dir string) (Project, error) {
	var meta *MetaInfo
	files := make(map[string]Document)
	reqs := make(map[string]RequestMeta)
	funcs := make(map[string]ValidateFunc)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return Project{}, errutil.Explain(nil, "read dir %s error: %w", dir, err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		fileName := e.Name()

		// Parse meta.json file if found
		if fileName == "meta.json" {
			var b []byte
			fileName = filepath.Join(dir, fileName)
			if b, err = os.ReadFile(fileName); err != nil {
				return Project{}, errutil.Explain(nil, "read file %s error: %w", fileName, err)
			}
			if meta, err = ParseMeta(b); err != nil {
				return Project{}, errutil.Explain(nil, "parse file %s error: %w", fileName, err)
			}
			continue
		}

		// Skip non-IDL files
		if !strings.HasSuffix(fileName, ".idl") {
			continue
		}

		var b []byte
		fileName = filepath.Join(dir, fileName)
		if b, err = os.ReadFile(fileName); err != nil {
			return Project{}, errutil.Explain(nil, "read file %s error: %w", fileName, err)
		}

		doc, validateFuncs, err := Parse(b)
		if err != nil {
			return Project{}, errutil.Explain(nil, "parse file %s error: %w", fileName, err)
		}
		files[e.Name()] = doc

		// validate request type
		for _, r := range doc.RPCs {
			reqs[r.Request] = RequestMeta{
				OnForm: strings.HasPrefix(r.ContentType, "application/x-www-form-urlencoded"),
			}
		}

		// record validate func
		for name, f := range validateFuncs {
			if v, ok := funcs[name]; !ok {
				funcs[name] = f
			} else if v.Type != f.Type {
				return Project{}, errutil.Explain(nil, "validate func %s is defined multiple times", name)
			}
		}
	}

	if meta == nil {
		return Project{}, errutil.Explain(nil, "no meta file")
	}
	if len(files) == 0 {
		return Project{}, errutil.Explain(nil, "no idl file")
	}

	// Validate that all used types are defined
	userTypes := map[string]struct{}{}
	definedTypes := make(map[string]struct{})
	for _, doc := range files {
		for k := range doc.EnumTypes {
			definedTypes[k] = struct{}{}
		}
		for k := range doc.TypeTypes {
			definedTypes[k] = struct{}{}
		}
		maps.Copy(userTypes, doc.UserTypes)
	}
	for k := range userTypes {
		if _, ok := definedTypes[k]; !ok {
			return Project{}, errutil.Explain(nil, "type %s is used but not defined", k)
		}
	}

	for _, doc := range files {
		for i := range doc.Types {
			t := doc.Types[i]
			t.Fields = t.RawFields
			if t.GenericParam != nil { // generic type, need instance
				// do nothing ...
			} else if t.InstType != nil { // generic type instance
				srcType, ok := GetType(files, t.InstType.BaseName)
				if !ok {
					return Project{}, errutil.Explain(nil, "type %s is used but not defined", t.InstType.BaseName)
				}
				var fields []TypeField
				for _, f := range srcType.Fields {
					f.Type = replaceGenericType(f.Type, *srcType.GenericParam, t.InstType.GenericType)
					fields = append(fields, f)
				}
				t.Fields = fields
			} else if t.Embedded {
				var fields []TypeField
				for _, f := range t.Fields {
					if e, ok := f.Type.(EmbedType); ok {
						srcType, ok := GetType(files, e.Name)
						if !ok {
							return Project{}, errutil.Explain(nil, "type %s is used but not defined", e.Name)
						}
						fields = append(fields, srcType.Fields...)
					} else {
						fields = append(fields, f)
					}
				}
				t.Fields = fields
			}

			if v, ok := reqs[t.Name]; ok {
				t.Request = true
				t.OnRequest = true
				t.OnForm = v.OnForm
			}
			doc.Types[i] = t // update
		}
	}

	// 一般来说，我们只需要最 request 类型进行 validate 操作
	for _, doc := range files {
		for _, t := range doc.Types {
			if t.Request {
				if _, err = getAndUpdateTypeValidate(files, t); err != nil {
					return Project{}, errutil.Explain(err, `failed to get type attr of type %s`, t.Name)
				}
			}
		}
	}

	for _, doc := range files {
		for i := range doc.RPCs {
			rpc := doc.RPCs[i]
			segments, err := pathidl.Parse(rpc.Path)
			if err != nil {
				return Project{}, errutil.Explain(err, `failed to parse path %s`, rpc.Path)
			}
			params := make(map[string]string)
			for _, seg := range segments {
				if seg.Type == pathidl.Static {
					continue
				}
				params[seg.Value] = ""
			}
			srcType, ok := GetType(files, rpc.Request)
			if !ok {
				return Project{}, errutil.Explain(nil, "type %s is used but not defined", rpc.Request)
			}
			for _, f := range srcType.Fields {
				if f.Binding == nil || f.Binding.Source != "path" {
					continue
				}
				if _, ok = params[f.Binding.Field]; !ok {
					err = errutil.Explain(nil, "path parameter %s not found in request type %s", f.Binding.Field, rpc.Request)
					return Project{}, err
				}
				params[f.Binding.Field] = f.Name
			}
			for k, s := range params {
				if s == "" {
					err = errutil.Explain(nil, "path parameter %s not found in request type %s", k, rpc.Request)
					return Project{}, err
				}
			}
			rpc.PathSegments = segments
			rpc.PathParams = params
			doc.RPCs[i] = rpc
		}
	}

	return Project{
		Meta:  meta,
		Files: files,
		Funcs: funcs,
	}, nil
}

func getAndUpdateTypeValidate(files map[string]Document, t Type) (bool, error) {
	t.OnRequest = true
	for i, f := range t.Fields {
		b, err := getTypeValidate(files, f.Type)
		if err != nil {
			return false, err
		}
		if b {
			f.ValidateNested = true
			t.Fields[i] = f
		}
		if f.Required || f.ValidateExpr != nil || f.ValidateNested {
			t.Validate = true
		}
	}
	fileName, index := FindType(files, t.Name)
	files[fileName].Types[index] = t // update
	return t.Validate, nil
}

func getTypeValidate(files map[string]Document, t TypeDefinition) (bool, error) {
	switch x := t.(type) {
	case UserType:
		srcType, ok := GetType(files, x.Name)
		if !ok {
			if _, ok = GetEnum(files, x.Name); !ok {
				return false, errutil.Explain(nil, "type %s is used but not defined", x.Name)
			}
			return false, nil
		}
		return getAndUpdateTypeValidate(files, srcType)
	case ListType:
		return getTypeValidate(files, x.Item)
	case MapType:
		return getTypeValidate(files, x.Value)
	default: // for linter
	}
	return false, nil
}

// replaceGenericType replaces a generic type with a concrete type.
func replaceGenericType(t TypeDefinition, genericName string, genericType TypeDefinition) TypeDefinition {
	switch u := t.(type) {
	case UserType:
		if u.Name == genericName {
			return genericType
		}
		return u
	case ListType:
		u.Item = replaceGenericType(u.Item, genericName, genericType)
		return u
	case MapType:
		u.Value = replaceGenericType(u.Value, genericName, genericType)
		return u
	default:
		return t
	}
}

// ParseMeta parses the JSON meta-information file.
func ParseMeta(data []byte) (*MetaInfo, error) {
	r := &MetaInfo{}
	if err := json.Unmarshal(data, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Parse runs the parsing pipeline for a single IDL input.
func Parse(data []byte) (doc Document, funcs map[string]ValidateFunc, err error) {
	if data = bytes.TrimSpace(data); len(data) == 0 {
		return Document{}, nil, nil
	}

	e := &ErrorListener{
		Scanner: bufio.NewScanner(bytes.NewReader(data)),
	}

	// Recover from parser panics to provide better error reporting
	defer func() {
		if r := recover(); r != nil {
			doc = Document{}
			err = errutil.Explain(nil, "[PANIC]: %v\n%s", r, debug.Stack())
			if e.Error != nil {
				err = errutil.Explain(nil, "%w\n%w", e.Error, err)
			}
		}
	}()

	// Step 1: Set up lexer and token stream
	reader := bytes.NewReader(append(data, '\n'))
	input := antlr.NewIoStream(reader)
	lexer := NewTLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(e)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Step 2: Set up parser
	p := NewTParser(tokens)
	p.RemoveErrorListeners()
	p.AddErrorListener(e)

	// Use faster SLL prediction first (fallback to LL on failure)
	p.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)

	// Step 3: Walk the parse tree with a custom listener
	l := &ParseTreeListener{
		Tokens: tokens,
		Document: Document{
			EnumTypes: make(map[string]int),
			TypeTypes: make(map[string]int),
			UserTypes: make(map[string]struct{}),
		},
		Attached: make(map[int]struct{}),
		Funcs:    make(map[string]ValidateFunc),
	}
	antlr.ParseTreeWalkerDefault.Walk(l, p.Document())

	// Step 4: Return result or error
	if e.Error != nil {
		return Document{}, nil, e.Error
	}
	return l.Document, l.Funcs, nil
}

// ErrorListener implements a custom ANTLR error listener.
type ErrorListener struct {
	*antlr.DefaultErrorListener
	Error   error
	Scanner *bufio.Scanner
	Line    int
}

// SyntaxError is called by ANTLR when a syntax error is encountered.
func (l *ErrorListener) SyntaxError(_ antlr.Recognizer, _ any, line, column int, msg string, e antlr.RecognitionException) {
	var text string
	for l.Scanner.Scan() {
		l.Line++
		if l.Line == line {
			text = l.Scanner.Text()
			break
		}
	}
	if l.Error == nil {
		l.Error = errutil.Explain(nil, "line %d:%d %s << text: %q", line, column, msg, text)
		return
	}
	l.Error = errutil.Explain(nil, "%w\nline %d:%d %s << text: %q", l.Error, line, column, msg, text)
}

// ParseTreeListener extends the auto-generated base listener.
// It captures parsed constructs (const, enum, type, rpc, etc.)
// and collects associated comments.
type ParseTreeListener struct {
	BaseTParserListener
	Tokens   *antlr.CommonTokenStream
	Document Document

	// Attached stores lines that already have "right-side" comments
	// to prevent re-using them as "top" comments.
	Attached map[int]struct{}

	// Funcs stores validate functions
	Funcs map[string]ValidateFunc
}

// ExitConst_def handles const definitions in the parse tree.
func (l *ParseTreeListener) ExitConst_def(ctx *Const_defContext) {
	c := Const{
		Type: BaseType{
			Name: ctx.Base_type().GetText(),
		},
		Name:  ctx.IDENTIFIER().GetText(),
		Value: ctx.Const_value().GetText(),
		Position: Position{
			StartLine: ctx.GetStart().GetLine(),
			EndLine:   ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Above: l.aboveComment(ctx.GetStart()),
			Right: l.rightComment(ctx.GetStop()),
		},
	}
	if !IsPascal(c.Name) {
		panic(errutil.Explain(nil, "const name %s is not PascalCase in line %d", c.Name, c.Position.StartLine))
	}
	l.Document.Consts = append(l.Document.Consts, c)
}

// ExitEnum_def handles enum definitions and their fields.
func (l *ParseTreeListener) ExitEnum_def(ctx *Enum_defContext) {
	e := Enum{
		Name: ctx.IDENTIFIER().GetText(),
		Position: Position{
			StartLine: ctx.GetStart().GetLine(),
			EndLine:   ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Above: l.aboveComment(ctx.GetStart()),
		},
	}
	if !IsPascal(e.Name) {
		panic(errutil.Explain(nil, "enum name %s is not PascalCase in line %d", e.Name, e.Position.StartLine))
	}

	for _, f := range ctx.AllEnum_field() {
		fieldName := f.IDENTIFIER().GetText()
		if !IsPascal(fieldName) {
			panic(errutil.Explain(nil, "enum field name %s is not PascalCase in line %d", fieldName, f.GetStart().GetLine()))
		}

		// Parse and validate integer value
		fieldValue := f.INTEGER().GetText()
		v, err := strconv.ParseInt(fieldValue, 0, 64)
		if err != nil {
			panic(errutil.Explain(nil, "enum field value %s is not a valid integer in line %d", fieldValue, f.GetStart().GetLine()))
		}

		enumField := EnumField{
			Name:  fieldName,
			Value: v,
			Position: Position{
				StartLine: f.GetStart().GetLine(),
				EndLine:   f.GetStop().GetLine(),
			},
			Comments: Comments{
				Above: l.aboveComment(f.GetStart()),
				Right: l.rightComment(f.GetStop()),
			},
		}

		// Annotations
		enumField.Annotations = l.parseFieldAnnotations(f.Field_annotations())

		// Error message
		if errmsg, ok := GetAnnotation(enumField.Annotations, "errmsg"); ok {
			if errmsg.Value == nil {
				panic(errutil.Explain(nil, `annotation "errmsg" value is nil in field %s of enum %s`, fieldName, e.Name))
			}
			s := strings.TrimSpace(strings.Trim(*errmsg.Value, `"`))
			enumField.ErrorMessage = &s
		}

		e.Fields = append(e.Fields, enumField)
	}

	l.Document.EnumTypes[e.Name] = len(l.Document.Enums)
	l.Document.Enums = append(l.Document.Enums, e)
}

// ExitType_def handles type definitions, including generic parameters,
// fields, and annotations.
func (l *ParseTreeListener) ExitType_def(ctx *Type_defContext) {
	t := Type{
		Name: ctx.IDENTIFIER(0).GetText(),
		Position: Position{
			StartLine: ctx.GetStart().GetLine(),
			EndLine:   ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Above: l.aboveComment(ctx.GetStart()),
		},
	}
	if !IsPascal(t.Name) {
		panic(errutil.Explain(nil, "type name %s is not PascalCase in line %d", t.Name, t.Position.StartLine))
	}

	if ctx.LEFT_BRACE() != nil {
		l.parseCompleteType(ctx, &t)
	} else {
		l.parseInstantiatedType(ctx, &t)
	}

	l.Document.TypeTypes[t.Name] = len(l.Document.Types)
	l.Document.Types = append(l.Document.Types, t)
}

// parseValueType resolves value types inside container types.
func (l *ParseTreeListener) parseValueType(ctx IValue_typeContext, t *Type) TypeDefinition {

	// Built-in bytes type
	if ctx.TYPE_BYTES() != nil {
		return BytesType{}
	}

	// Built-in primitive type
	if ctx.Base_type() != nil {
		return BaseType{
			Name: ctx.Base_type().GetText(),
		}
	}

	// Reference to a user-defined type
	if ctx.User_type() != nil {
		ut := UserType{
			Name: ctx.User_type().IDENTIFIER().GetText(),
		}
		// Track user-defined types
		if t == nil || t.GenericParam == nil || ut.Name != *t.GenericParam {
			l.Document.UserTypes[ut.Name] = struct{}{}
		}
		return ut
	}

	// Container types (map / list)
	if c := ctx.Container_type(); c != nil {
		if c.Map_type() != nil {
			kt := c.Map_type().Key_type().GetText()
			vt := l.parseValueType(c.Map_type().Value_type(), t)
			return MapType{Key: kt, Value: vt}
		}
		if c.List_type() != nil {
			vt := l.parseValueType(c.List_type().Value_type(), t)
			return ListType{Item: vt}
		}
	}

	panic(errutil.Explain(nil, "invalid type %s in line %d", ctx.GetText(), ctx.GetStart().GetLine()))
}

// parseInstantiatedType handles instantiated types, including generic types.
func (l *ParseTreeListener) parseInstantiatedType(ctx *Type_defContext, t *Type) {
	t.InstType = &InstType{
		BaseName: ctx.IDENTIFIER(1).GetText(),
	}
	if !IsPascal(t.InstType.BaseName) {
		panic(errutil.Explain(nil, "type name %s is not PascalCase in line %d", t.InstType.BaseName, t.Position.StartLine))
	}

	t.InstType.GenericType = l.parseValueType(ctx.Value_type(), t)
}

// parseCompleteType handles a "struct-like" type with fields and optional generic parameter.
func (l *ParseTreeListener) parseCompleteType(ctx *Type_defContext, t *Type) {

	// Handle generic type parameter (if any)
	if ctx.LESS_THAN() != nil {
		s := ctx.IDENTIFIER(1).GetText()
		t.GenericParam = &s
	}

	for _, f := range ctx.AllType_field() {
		typeField := TypeField{
			Position: Position{
				StartLine: f.GetStart().GetLine(),
				EndLine:   f.GetStop().GetLine(),
			},
			Comments: Comments{
				Above: l.aboveComment(f.GetStart()),
				Right: l.rightComment(f.GetStop()),
			},
		}

		// Distinguish between embedded fields and normal fields
		if f.Embed_type_field() != nil {
			t.Embedded = true
			embedType := EmbedType{
				Name: f.Embed_type_field().User_type().IDENTIFIER().GetText(),
			}
			// Track user-defined types
			if t.GenericParam == nil || embedType.Name != *t.GenericParam {
				l.Document.UserTypes[embedType.Name] = struct{}{}
			}
			typeField.Type = embedType

		} else if f.Common_type_field() != nil {
			l.parseCommonTypeField(f.Common_type_field(), &typeField, t)
		}

		t.RawFields = append(t.RawFields, typeField)
	}
}

// parseCommonTypeField parses a regular field (not embedded) inside a type or oneof.
func (l *ParseTreeListener) parseCommonTypeField(f ICommon_type_fieldContext, typeField *TypeField, t *Type) {
	typeField.Type = l.parseValueType(f.Value_type(), t)
	typeField.Name = f.IDENTIFIER().GetText()
	typeField.Required = f.KW_REQUIRED() != nil
	typeField.Annotations = l.parseFieldAnnotations(f.Field_annotations())

	if opt, ok := GetAnnotation(typeField.Annotations, "compat_absent"); ok {
		if !typeField.Required {
			panic(errutil.Explain(nil, "field %s is required but has compat_absent annotation in line %d", typeField.Name, typeField.Position.StartLine))
		}
		if opt.Value == nil {
			panic(errutil.Explain(nil, "annotation compat_absent for field %s is missing value in line %d", typeField.Name, typeField.Position.StartLine))
		}
		s := strings.TrimSpace(*opt.Value)
		if s == "" {
			panic(errutil.Explain(nil, "annotation compat_absent for field %s is empty in line %d", typeField.Name, typeField.Position.StartLine))
		}
		s = strings.TrimSpace(strings.Trim(s, "\"")) // Remove quotes
		typeField.CompatAbsent = &s
	}

	_, typeField.EnumAsString = GetAnnotation(typeField.Annotations, "enum_as_string")

	typeField.JSONTag = JSONTag{Name: typeField.Name, OmitEmpty: true}
	if opt, ok := GetAnnotation(typeField.Annotations, "json"); ok {
		if opt.Value == nil {
			panic(errutil.Explain(nil, "annotation json for field %s is missing value in line %d", typeField.Name, typeField.Position.StartLine))
		}
		s := strings.TrimSpace(*opt.Value)
		if s == "" {
			panic(errutil.Explain(nil, "annotation json for field %s is empty in line %d", typeField.Name, typeField.Position.StartLine))
		}
		s = strings.TrimSpace(strings.Trim(s, "\"")) // Remove quotes
		for i, v := range strings.Split(s, ",") {
			v = strings.TrimSpace(v)
			if i == 0 {
				if v != "" {
					typeField.JSONTag.Name = v
				}
				continue
			}
			switch v {
			case "omitempty":
				typeField.JSONTag.OmitEmpty = true
			case "non-omitempty":
				typeField.JSONTag.OmitEmpty = false
			default: // for linter
			}
		}
	}

	typeField.FormTag = FormTag{Name: typeField.JSONTag.Name}
	if opt, ok := GetAnnotation(typeField.Annotations, "form"); ok {
		if opt.Value == nil {
			panic(errutil.Explain(nil, "annotation form for field %s is missing value in line %d", typeField.Name, typeField.Position.StartLine))
		}
		s := strings.TrimSpace(*opt.Value)
		if s == "" {
			panic(errutil.Explain(nil, "annotation form for field %s is empty in line %d", typeField.Name, typeField.Position.StartLine))
		}
		s = strings.TrimSpace(strings.Trim(s, "\"")) // Remove quotes
		for i, v := range strings.Split(s, ",") {
			v = strings.TrimSpace(v)
			if i == 0 {
				if v != "" {
					typeField.FormTag.Name = v
				}
				continue
			}
			// ...
		}
	}

	if opt, ok := GetAnnotation(typeField.Annotations, "path", "query"); ok {
		if opt.Key == "path" {
			if s := typeField.Type.Text(); s != "string" && s != "int" {
				panic(errutil.Explain(nil, "annotation path for field %s is not string or int in line %d", typeField.Name, typeField.Position.StartLine))
			}
		}
		if opt.Value == nil {
			panic(errutil.Explain(nil, "annotation %s for field %s is missing value in line %d", opt.Key, typeField.Name, typeField.Position.StartLine))
		}
		s := strings.TrimSpace(*opt.Value)
		if s == "" {
			panic(errutil.Explain(nil, "annotation %s for field %s is empty in line %d", opt.Key, typeField.Name, typeField.Position.StartLine))
		}
		s = strings.TrimSpace(strings.Trim(s, "\"")) // Remove quotes
		typeField.Binding = &Binding{Source: opt.Key, Field: s}
	}

	if opt, ok := GetAnnotation(typeField.Annotations, "validate"); ok {
		if opt.Value == nil {
			panic(errutil.Explain(nil, "annotation validate for field %s is missing value in line %d", typeField.Name, typeField.Position.StartLine))
		}
		s := strings.TrimSpace(*opt.Value)
		if s == "" {
			panic(errutil.Explain(nil, "annotation validate for field %s is empty in line %d", typeField.Name, typeField.Position.StartLine))
		}
		s, err := strconv.Unquote(s) // Remove quotes
		if err != nil {
			panic(errutil.Explain(nil, `annotation validate for field %s value is not properly quoted in line %d`, typeField.Name, typeField.Position.StartLine))
		}
		if s = strings.TrimSpace(s); s == "" {
			panic(errutil.Explain(nil, `annotation validate for field %s value is empty in line %d`, typeField.Name, typeField.Position.StartLine))
		}
		typeField.ValidateExpr, err = validate.Parse(s)
		if err != nil {
			panic(errutil.Explain(err, `failed to parse validate expression %s`, *opt.Value))
		}
		l.collectValidateFuncs(typeField.Type.Text(), typeField.ValidateExpr)
	}
}

// parseFieldAnnotations parses field annotations.
func (l *ParseTreeListener) parseFieldAnnotations(ctx IField_annotationsContext) []Annotation {
	if ctx == nil {
		return nil
	}
	var ret []Annotation
	for _, aCtx := range ctx.AllAnnotation() {
		a := Annotation{
			Key: aCtx.IDENTIFIER().GetText(),
			Position: Position{
				StartLine: aCtx.GetStart().GetLine(),
				EndLine:   aCtx.GetStop().GetLine(),
			},
		}
		if aCtx.Const_value() != nil {
			s := aCtx.Const_value().GetText()
			a.Value = &s
		}
		ret = append(ret, a)
	}
	return ret
}

// collectValidateFuncs collects validate functions from the given expression.
func (l *ParseTreeListener) collectValidateFuncs(fieldType string, expr validate.Expr) {
	switch x := expr.(type) {
	case validate.PrimaryExpr:
		if x.Inner != nil {
			l.collectValidateFuncs(fieldType, x.Inner)
		} else if x.Call != nil {
			l.collectValidateFuncs(fieldType, x.Call)
		}
	case *validate.InnerExpr:
		l.collectValidateFuncs(fieldType, x.Expr)
	case validate.UnaryExpr:
		l.collectValidateFuncs(fieldType, x.Expr)
	case validate.BinaryExpr:
		l.collectValidateFuncs(fieldType, x.Left)
		l.collectValidateFuncs(fieldType, x.Right)
	case *validate.FuncCall:
		if _, ok := BuiltinFuncs[x.Name]; !ok {
			if v, ok := l.Funcs[x.Name]; !ok {
				l.Funcs[x.Name] = ValidateFunc{
					Name: x.Name,
					Type: fieldType,
				}
			} else if v.Type != fieldType {
				panic(errutil.Explain(nil, "validate function %s is used with different types", x.Name))
			}
		}
		for _, arg := range x.Args {
			l.collectValidateFuncs(fieldType, arg)
		}
	default:
		panic(errutil.Explain(nil, "unexpected validate expression type %T", x))
	}
}

// ExitOneof_def handles "oneof" type definitions.
func (l *ParseTreeListener) ExitOneof_def(ctx *Oneof_defContext) {
	o := Type{
		Name:  ctx.IDENTIFIER().GetText(),
		OneOf: true,
		Position: Position{
			StartLine: ctx.GetStart().GetLine(),
			EndLine:   ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Above: l.aboveComment(ctx.GetStart()),
		},
	}
	if !IsPascal(o.Name) {
		panic(errutil.Explain(nil, "oneof name %s is not PascalCase in line %d", o.Name, o.Position.StartLine))
	}

	e := Enum{Name: o.Name + "Type", OneOf: true}
	o.RawFields = append(o.RawFields, TypeField{
		Name: "FieldType",
		Type: UserType{Name: e.Name},
		Annotations: []Annotation{
			{Key: "enum_as_string"},
		},
		JSONTag: JSONTag{
			Name: "FieldType",
		},
		FormTag: FormTag{
			Name: "FieldType",
		},
		Required:     true,
		EnumAsString: true,
	})

	for i, f := range ctx.AllUser_type() {

		// add enum fields
		e.Fields = append(e.Fields, EnumField{
			Name:  f.IDENTIFIER().GetText(),
			Value: int64(i),
		})

		typeField := TypeField{
			Name: f.IDENTIFIER().GetText(),
			Type: UserType{
				Name: f.IDENTIFIER().GetText(),
			},
			Position: Position{
				StartLine: f.GetStart().GetLine(),
				EndLine:   f.GetStop().GetLine(),
			},
			Comments: Comments{
				Above: l.aboveComment(f.GetStart()),
				Right: l.rightComment(f.GetStop()),
			},
			JSONTag: JSONTag{
				Name:      f.IDENTIFIER().GetText(),
				OmitEmpty: true,
			},
			FormTag: FormTag{
				Name: f.IDENTIFIER().GetText(),
			},
		}

		o.RawFields = append(o.RawFields, typeField)
	}

	l.Document.EnumTypes[e.Name] = len(l.Document.EnumTypes)
	l.Document.Enums = append(l.Document.Enums, e)

	l.Document.TypeTypes[o.Name] = len(l.Document.Types)
	l.Document.Types = append(l.Document.Types, o)
}

// ExitRpc_def handles RPC definitions, including request/response
// types and annotations.
func (l *ParseTreeListener) ExitRpc_def(ctx *Rpc_defContext) {
	var err error

	sse := false
	if ctx.KW_SSE() != nil {
		sse = true
	}

	r := RPC{
		SSE:  sse,
		Name: ctx.IDENTIFIER().GetText(),
		Position: Position{
			StartLine: ctx.GetStart().GetLine(),
			EndLine:   ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Above: l.aboveComment(ctx.GetStart()),
		},
	}
	if !IsPascal(r.Name) {
		panic(errutil.Explain(nil, "RPC name %s is not PascalCase in line %d", r.Name, r.Position.StartLine))
	}

	// Request
	r.Request = ctx.Rpc_req().User_type().IDENTIFIER().GetText()
	l.Document.UserTypes[r.Request] = struct{}{}

	// Response
	r.Response = l.parseValueType(ctx.Rpc_resp().Value_type(), nil)

	// Annotations
	for _, aCtx := range ctx.Rpc_annotations().AllAnnotation() {
		a := Annotation{
			Key: aCtx.IDENTIFIER().GetText(),
			Position: Position{
				StartLine: aCtx.GetStart().GetLine(),
				EndLine:   aCtx.GetStop().GetLine(),
			},
			Comments: Comments{
				Above: l.aboveComment(aCtx.GetStart()),
				Right: l.rightComment(aCtx.GetStop()),
			},
		}
		if aCtx.Const_value() != nil {
			s := aCtx.Const_value().GetText()
			a.Value = &s
		}
		r.Annotations = append(r.Annotations, a)
	}

	// Retrieve the "path" annotation
	path, ok := GetAnnotation(r.Annotations, "path")
	if !ok {
		panic(errutil.Explain(nil, `annotation "path" not found in rpc %s`, r.Name))
	}
	if path.Value == nil {
		panic(errutil.Explain(nil, `annotation "path" value is nil in rpc %s`, r.Name))
	}

	// Retrieve the "method" annotation
	method, ok := GetAnnotation(r.Annotations, "method")
	if !ok {
		panic(errutil.Explain(nil, `annotation "method" not found in rpc %s`, r.Name))
	}
	if method.Value == nil {
		panic(errutil.Explain(nil, `annotation "method" value is nil in rpc %s`, r.Name))
	}

	// Retrieve the "contentType" annotation
	ct, ok := GetAnnotation(r.Annotations, "contentType")
	if !ok {
		panic(errutil.Explain(nil, `annotation "contentType" not found in rpc %s`, r.Name))
	}
	if ct.Value == nil {
		panic(errutil.Explain(nil, `annotation "contentType" value is nil in rpc %s`, r.Name))
	}

	var contentType string
	switch s := strings.TrimSpace(strings.Trim(*ct.Value, `"`)); s {
	case "form":
		contentType = "application/x-www-form-urlencoded"
	case "json":
		contentType = "application/json"
	default:
		contentType = s
	}

	// Retrieve the "connTimeout" annotation
	connTimeout, ok := GetAnnotation(r.Annotations, "connTimeout")
	if !ok {
		panic(errutil.Explain(nil, `annotation "connTimeout" not found in rpc %s`, r.Name))
	}
	if connTimeout.Value == nil {
		panic(errutil.Explain(nil, `annotation "connTimeout" value is nil in rpc %s`, r.Name))
	}

	// Retrieve the "readTimeout" annotation
	readTimeout, ok := GetAnnotation(r.Annotations, "readTimeout")
	if !ok {
		panic(errutil.Explain(nil, `annotation "readTimeout" not found in rpc %s`, r.Name))
	}
	if readTimeout.Value == nil {
		panic(errutil.Explain(nil, `annotation "readTimeout" value is nil in rpc %s`, r.Name))
	}

	// Retrieve the "writeTimeout" annotation
	writeTimeout, ok := GetAnnotation(r.Annotations, "writeTimeout")
	if !ok {
		panic(errutil.Explain(nil, `annotation "writeTimeout" not found in rpc %s`, r.Name))
	}
	if writeTimeout.Value == nil {
		panic(errutil.Explain(nil, `annotation "writeTimeout" value is nil in rpc %s`, r.Name))
	}

	r.Path = strings.Trim(*path.Value, `"`)
	r.Method = strings.ToUpper(strings.Trim(*method.Value, `"`))
	r.ContentType = contentType

	r.ConnTimeout, err = strconv.Atoi(strings.Trim(*connTimeout.Value, `"`))
	if err != nil {
		panic(errutil.Explain(nil, "invalid connTimeout value in rpc %s", r.Name))
	}
	r.ReadTimeout, err = strconv.Atoi(strings.Trim(*readTimeout.Value, `"`))
	if err != nil {
		panic(errutil.Explain(nil, "invalid readTimeout value in rpc %s", r.Name))
	}
	r.WriteTimeout, err = strconv.Atoi(strings.Trim(*writeTimeout.Value, `"`))
	if err != nil {
		panic(errutil.Explain(nil, "invalid writeTimeout value in rpc %s", r.Name))
	}

	l.Document.RPCs = append(l.Document.RPCs, r)
}

// isTerminatorToken returns true if the given token is considered a statement terminator.
// In this parser, a newline or semicolon marks the end of a statement.
func isTerminatorToken(t antlr.Token) bool {
	return t.GetTokenType() == TLexerNEWLINE
}

// previousTokenOnChannel finds the index of the previous token that is on the default channel.
// It skips terminator tokens (newline/semicolon) and tokens on hidden channels.
func (l *ParseTreeListener) previousTokenOnChannel(i int) int {
	tokens := l.Tokens.GetAllTokens()
	for i >= 0 && (isTerminatorToken(tokens[i]) || tokens[i].GetChannel() != antlr.LexerDefaultTokenChannel) {
		i--
	}
	return i
}

// filterForChannel returns a slice of tokens between indices [left, right] that belong to the given channel.
// channel = -1 means "all hidden channels".
func (l *ParseTreeListener) filterForChannel(left, right, channel int) []antlr.Token {
	tokens := l.Tokens.GetAllTokens()
	hidden := make([]antlr.Token, 0)
	for i := left; i < right+1; i++ {
		t := tokens[i]
		if channel == -1 {
			if t.GetChannel() != antlr.LexerDefaultTokenChannel {
				hidden = append(hidden, t)
			}
		} else if t.GetChannel() == channel {
			hidden = append(hidden, t)
		}
	}
	if len(hidden) == 0 {
		return nil
	}
	return hidden
}

// GetHiddenTokensToLeft returns all hidden tokens to the left of a given token index
// that belong to the specified channel.
func (l *ParseTreeListener) GetHiddenTokensToLeft(tokenIndex, channel int) []antlr.Token {
	tokens := l.Tokens.GetAllTokens()
	if tokenIndex < 0 || tokenIndex >= len(tokens) {
		panic(strconv.Itoa(tokenIndex) + " not in 0.." + strconv.Itoa(len(tokens)-1))
	}

	prevOnChannel := l.previousTokenOnChannel(tokenIndex - 1)
	if prevOnChannel == tokenIndex-1 {
		return nil
	}

	// If there are none on channel to the left and prevOnChannel == -1 then from = 0
	from := prevOnChannel + 1
	to := tokenIndex - 1
	return l.filterForChannel(from, to, channel)
}

// nextTokenOnChannel finds the next token index on the default channel,
// skipping terminators and hidden tokens.
// Returns -1 if no such token exists.
func (l *ParseTreeListener) nextTokenOnChannel(i int) int {
	tokens := l.Tokens.GetAllTokens()
	if i >= len(tokens) {
		return -1
	}
	token := tokens[i]
	for isTerminatorToken(tokens[i]) || token.GetChannel() != antlr.LexerDefaultTokenChannel {
		if token.GetTokenType() == antlr.TokenEOF {
			return -1
		}
		i++
		token = tokens[i]
	}
	return i
}

// GetHiddenTokensToRight returns all hidden tokens to the right of a given token index
// that belong to the specified channel.
func (l *ParseTreeListener) GetHiddenTokensToRight(tokenIndex, channel int) []antlr.Token {
	tokens := l.Tokens.GetAllTokens()
	if tokenIndex < 0 || tokenIndex >= len(tokens) {
		panic(strconv.Itoa(tokenIndex) + " not in 0.." + strconv.Itoa(len(tokens)-1))
	}

	nextOnChannel := l.nextTokenOnChannel(tokenIndex + 1)
	from := tokenIndex + 1

	// If no onChannel to the right, then nextOnChannel == -1, so set 'to' to the last token
	var to int
	if nextOnChannel == -1 {
		to = len(tokens) - 1
	} else {
		to = nextOnChannel
	}
	return l.filterForChannel(from, to, channel)
}

// formatSingleLineComment trims and normalizes a single-line comment text.
// It ensures the comment starts with "// " and removes extra whitespace.
func formatSingleLineComment(text string) string {
	s := strings.TrimSpace(text)
	s = strings.TrimSpace(strings.TrimPrefix(s, "//"))
	return "// " + s
}

// formatMultiLineComment splits a multi-line comment (/* ... */) into normalized lines.
// Each line is trimmed, and leading '*' is standardized.
func formatMultiLineComment(text string) []string {
	var lines []string
	for i, s := range strings.Split(text, "\n") {
		s = strings.TrimSpace(s)
		if i == 0 {
			s = strings.TrimSpace(strings.TrimPrefix(s, "/*"))
			if s == "" {
				s = "/*"
			} else {
				s = "/* " + s
			}
		} else {
			if strings.HasSuffix(s, "*/") {
				s = strings.TrimSpace(s[:len(s)-2]) + " */"
			}
			if strings.HasPrefix(s, "*") {
				s = " * " + strings.TrimSpace(s[1:])
			}
		}
		lines = append(lines, s)
	}
	return lines
}

// aboveComment extracts all comments immediately above a token.
// Supports both single-line (//) and multi-line (/* */) comments.
// Returns only the contiguous block directly attached to the token.
func (l *ParseTreeListener) aboveComment(token antlr.Token) []Comment {
	var (
		all []Comment // all collected comments
		ret []Comment // contiguous block directly above token
	)

	// Collect single-line comments
	comments := l.GetHiddenTokensToLeft(token.GetTokenIndex(), TLexerSL_COMMENT_CHAN)
	for _, c := range comments {
		if _, ok := l.Attached[c.GetLine()]; ok {
			continue
		}
		line := formatSingleLineComment(c.GetText())
		all = append(all, Comment{
			Text:   []string{line},
			Single: true,
			Position: Position{
				StartLine: c.GetLine(),
				EndLine:   c.GetLine(),
			},
		})
	}

	// Collect multi-line comments
	comments = l.GetHiddenTokensToLeft(token.GetTokenIndex(), TLexerML_COMMENT_CHAN)
	for _, c := range comments {
		if _, ok := l.Attached[c.GetLine()]; ok {
			continue
		}
		lines := formatMultiLineComment(c.GetText())
		all = append(all, Comment{
			Text:   lines,
			Single: false,
			Position: Position{
				StartLine: c.GetLine(),
				EndLine:   c.GetLine() + len(lines) - 1,
			},
		})
	}

	// Sort comments by starting line in descending order
	sort.Slice(all, func(i, j int) bool {
		return all[i].Position.StartLine >= all[j].Position.StartLine
	})

	// Select only the contiguous block directly above the token
	i := 0
	lastLine := token.GetLine()
	for ; i < len(all); i++ {
		c := all[i]
		if c.Position.EndLine != lastLine-1 {
			break
		}
		ret = append([]Comment{c}, ret...)
		lastLine = c.Position.StartLine
	}

	// Remaining comments are stored as detached comments in the Document
	for j := len(all) - 1; j >= i; j-- {
		l.Document.Comments = append(l.Document.Comments, all[j])
	}

	return ret
}

// rightComment extracts a comment that appears on the same line as a token.
// Supports both single-line and multi-line comments.
func (l *ParseTreeListener) rightComment(token antlr.Token) *Comment {
	// Single-line comments
	comments := l.Tokens.GetHiddenTokensToRight(token.GetTokenIndex(), TLexerSL_COMMENT_CHAN)
	for _, c := range comments {
		if c.GetLine() != token.GetLine() {
			continue
		}
		l.Attached[c.GetLine()] = struct{}{}
		line := formatSingleLineComment(c.GetText())
		return &Comment{
			Text:   []string{line},
			Single: true,
			Position: Position{
				StartLine: c.GetLine(),
				EndLine:   c.GetLine(),
			},
		}
	}

	// Multi-line comments
	comments = l.Tokens.GetHiddenTokensToRight(token.GetTokenIndex(), TLexerML_COMMENT_CHAN)
	for _, c := range comments {
		if c.GetLine() != token.GetLine() {
			continue
		}
		l.Attached[c.GetLine()] = struct{}{}
		lines := formatMultiLineComment(c.GetText())
		return &Comment{
			Text:   lines,
			Single: false,
			Position: Position{
				StartLine: c.GetLine(),
				EndLine:   c.GetLine() + len(lines) - 1,
			},
		}
	}

	return nil
}
