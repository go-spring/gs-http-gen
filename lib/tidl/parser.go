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

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// ParseDir scans the specified directory for IDL files (*.idl) and a meta.json file.
// It parses each file into a Document structure and validates cross-file type references.
func ParseDir(dir string) (files map[string]Document, meta *MetaInfo, err error) {
	files = make(map[string]Document)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("read dir %s error: %w", dir, err)
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
				return nil, nil, fmt.Errorf("read file %s error: %w", fileName, err)
			}
			if meta, err = ParseMeta(b); err != nil {
				return nil, nil, fmt.Errorf("parse file %s error: %w", fileName, err)
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
			return nil, nil, fmt.Errorf("read file %s error: %w", fileName, err)
		}
		var doc Document
		if doc, err = Parse(b); err != nil {
			return nil, nil, fmt.Errorf("parse file %s error: %w", fileName, err)
		}
		files[e.Name()] = doc
	}

	// Validate that all used types are defined
	usedTypes := map[string]struct{}{}
	definedTypes := make(map[string]struct{})
	for _, doc := range files {
		maps.Copy(usedTypes, doc.UsedTypes)
		for k := range doc.EnumTypes {
			definedTypes[k] = struct{}{}
		}
		for k := range doc.TypeTypes {
			definedTypes[k] = struct{}{}
		}
	}
	for k := range usedTypes {
		if _, ok := definedTypes[k]; !ok {
			return nil, nil, fmt.Errorf("type %s is used but not defined", k)
		}
	}

	return
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
func Parse(data []byte) (doc Document, err error) {
	if data = bytes.TrimSpace(data); len(data) == 0 {
		return Document{}, nil
	}

	e := &ErrorListener{
		Scanner: bufio.NewScanner(bytes.NewReader(data)),
	}

	// Recover from parser panics to provide better error reporting
	defer func() {
		if r := recover(); r != nil {
			doc = Document{}
			err = fmt.Errorf("[PANIC]: %v\n%s", r, debug.Stack())
			if e.Error != nil {
				err = fmt.Errorf("%w\n%w", e.Error, err)
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
			UsedTypes: make(map[string]struct{}),
		},
		Attached: make(map[int]struct{}),
	}
	antlr.ParseTreeWalkerDefault.Walk(l, p.Document())

	// Step 4: Return result or error
	if e.Error != nil {
		return Document{}, e.Error
	}
	return l.Document, nil
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
		l.Error = fmt.Errorf("line %d:%d %s << text: %q", line, column, msg, text)
		return
	}
	l.Error = fmt.Errorf("%w\nline %d:%d %s << text: %q", l.Error, line, column, msg, text)
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
}

// ExitConst_def handles const definitions in the parse tree.
func (l *ParseTreeListener) ExitConst_def(ctx *Const_defContext) {
	c := Const{
		Type:  ctx.Const_type().GetText(),
		Name:  ctx.IDENTIFIER().GetText(),
		Value: ctx.Const_value().GetText(),
		Position: Position{
			Start: ctx.GetStart().GetLine(),
			Stop:  ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Above: l.aboveComment(ctx.GetStart()),
			Right: l.rightComment(ctx.GetStop()),
		},
	}
	if !IsPascal(c.Name) {
		panic(fmt.Errorf("const name %s is not PascalCase in line %d", c.Name, c.Position.Start))
	}
	l.Document.Consts = append(l.Document.Consts, c)
}

// ExitEnum_def handles enum definitions and their fields.
func (l *ParseTreeListener) ExitEnum_def(ctx *Enum_defContext) {
	e := Enum{
		Name: ctx.IDENTIFIER().GetText(),
		Position: Position{
			Start: ctx.GetStart().GetLine(),
			Stop:  ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Above: l.aboveComment(ctx.GetStart()),
		},
	}
	if !IsPascal(e.Name) {
		panic(fmt.Errorf("enum name %s is not PascalCase in line %d", e.Name, e.Position.Start))
	}

	for _, f := range ctx.AllEnum_field() {
		fieldName := f.IDENTIFIER().GetText()
		if !IsPascal(fieldName) {
			panic(fmt.Errorf("enum field name %s is not PascalCase in line %d", fieldName, f.GetStart().GetLine()))
		}

		// Parse and validate integer value
		fieldValue := f.INTEGER().GetText()
		v, err := strconv.ParseInt(fieldValue, 0, 64)
		if err != nil {
			panic(fmt.Errorf("enum field value %s is not a valid integer in line %d", fieldValue, f.GetStart().GetLine()))
		}

		e.Fields = append(e.Fields, EnumField{
			Name:  fieldName,
			Value: v,
			Position: Position{
				Start: f.GetStart().GetLine(),
				Stop:  f.GetStop().GetLine(),
			},
			Comments: Comments{
				Above: l.aboveComment(f.GetStart()),
				Right: l.rightComment(f.GetStop()),
			},
		})
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
			Start: ctx.GetStart().GetLine(),
			Stop:  ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Above: l.aboveComment(ctx.GetStart()),
		},
	}
	if !IsPascal(t.Name) {
		panic(fmt.Errorf("type name %s is not PascalCase in line %d", t.Name, t.Position.Start))
	}

	// Distinguish between a full struct definition and a type alias
	if ctx.LEFT_BRACE() != nil {
		l.parseCompleteType(ctx, &t)
	} else {
		l.parseRedefinedType(ctx, &t)
	}

	l.Document.TypeTypes[t.Name] = len(l.Document.Types)
	l.Document.Types = append(l.Document.Types, t)
}

// parseCompleteType handles a "struct-like" type with fields and optional generic parameter.
func (l *ParseTreeListener) parseCompleteType(ctx *Type_defContext, t *Type) {

	// Handle generic type parameter (if any)
	if ctx.LESS_THAN() != nil {
		s := ctx.IDENTIFIER(1).GetText()
		t.GenericName = &s
	}

	for _, f := range ctx.AllType_field() {
		typeField := TypeField{
			Position: Position{
				Start: f.GetStart().GetLine(),
				Stop:  f.GetStop().GetLine(),
			},
			Comments: Comments{
				Above: l.aboveComment(f.GetStart()),
				Right: l.rightComment(f.GetStop()),
			},
		}

		// Distinguish between embedded fields and normal fields
		if etf := f.Embed_type_field(); etf != nil {
			u := etf.User_type()
			embedType := EmbedType{
				Name:     u.IDENTIFIER().GetText(),
				Optional: u.QUESTION() != nil,
			}
			if t.GenericName == nil || embedType.Name != *t.GenericName {
				l.Document.UsedTypes[embedType.Name] = struct{}{}
			}
			typeField.FieldType = embedType

		} else if ctf := f.Common_type_field(); ctf != nil {
			l.parseCommonTypeField(ctf, &typeField, t)
		}

		t.Fields = append(t.Fields, typeField)
	}
}

// parseRedefinedType handles redefined types, including generic types.
func (l *ParseTreeListener) parseRedefinedType(ctx *Type_defContext, t *Type) {
	t.Redefined = &RedefinedType{
		Name: ctx.IDENTIFIER(1).GetText(),
	}
	if !IsPascal(t.Redefined.Name) {
		panic(fmt.Errorf("redefined type name %s is not PascalCase in line %d", t.Redefined.Name, t.Position.Start))
	}

	t.Redefined.GenericType = l.parseValueType(ctx.Value_type(), t)
	if t.Redefined.GenericType != nil {
		return
	}

	panic(fmt.Errorf("redefined type %s is not a valid generic type in line %d", t.Redefined.Name, t.Position.Start))
}

// ExitOneof_def handles "oneof" type definitions.
func (l *ParseTreeListener) ExitOneof_def(ctx *Oneof_defContext) {
	o := Type{
		Name:  ctx.IDENTIFIER().GetText(),
		OneOf: true,
		Position: Position{
			Start: ctx.GetStart().GetLine(),
			Stop:  ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Above: l.aboveComment(ctx.GetStart()),
		},
	}
	if !IsPascal(o.Name) {
		panic(fmt.Errorf("oneof name %s is not PascalCase in line %d", o.Name, o.Position.Start))
	}

	for _, f := range ctx.AllCommon_type_field() {
		typeField := TypeField{
			Position: Position{
				Start: f.GetStart().GetLine(),
				Stop:  f.GetStop().GetLine(),
			},
			Comments: Comments{
				Above: l.aboveComment(f.GetStart()),
				Right: l.rightComment(f.GetStop()),
			},
		}
		l.parseCommonTypeField(f, &typeField, &o)
		o.Fields = append(o.Fields, typeField)
	}

	l.Document.TypeTypes[o.Name] = len(l.Document.Types)
	l.Document.Types = append(l.Document.Types, o)
}

// parseCommonTypeField parses a regular field (not embedded) inside a type or oneof.
func (l *ParseTreeListener) parseCommonTypeField(f ICommon_type_fieldContext, typeField *TypeField, t *Type) {
	typeField.FieldType = l.parseCommonFieldType(f.Common_field_type(), t)
	typeField.Name = f.IDENTIFIER().GetText()

	// Annotations
	if f.Type_annotations() != nil {
		for _, aCtx := range f.Type_annotations().AllAnnotation() {
			a := Annotation{
				Key: aCtx.IDENTIFIER().GetText(),
				Position: Position{
					Start: aCtx.GetStart().GetLine(),
					Stop:  aCtx.GetStop().GetLine(),
				},
			}
			if aCtx.Const_value() != nil {
				s := aCtx.Const_value().GetText()
				a.Value = &s
			}
			typeField.Annotations = append(typeField.Annotations, a)
		}
	}
}

// parseCommonFieldType distinguishes between built-in, user-defined, or container types.
func (l *ParseTreeListener) parseCommonFieldType(ctx ICommon_field_typeContext, t *Type) TypeDefinition {
	if ctx.TYPE_ANY() != nil {
		return AnyType{}
	}
	if ctx.TYPE_BINARY() != nil {
		return BinaryType{}
	}
	return l.parseValueType(ctx, t)
}

// parseValueType resolves value types inside container types.
func (l *ParseTreeListener) parseValueType(ctx interface {
	GetText() string
	GetStart() antlr.Token
	GetStop() antlr.Token
	Base_type() IBase_typeContext
	User_type() IUser_typeContext
	Container_type() IContainer_typeContext
}, t *Type) TypeDefinition {

	// Built-in primitive type
	if b := ctx.Base_type(); b != nil {
		return BaseType{
			Name:     strings.TrimRight(b.GetText(), "?"),
			Optional: b.QUESTION() != nil,
		}
	}

	// Reference to a user-defined type
	if u := ctx.User_type(); u != nil {
		typ := UserType{
			Name:     u.IDENTIFIER().GetText(),
			Optional: u.QUESTION() != nil,
		}
		if t.GenericName == nil || typ.Name != *t.GenericName {
			l.Document.UsedTypes[typ.Name] = struct{}{}
		}
		return typ
	}

	// Container types (map / list)
	if c := ctx.Container_type(); c != nil {
		if c.Map_type() != nil {
			kt := c.Map_type().Key_type().GetText()
			vt := l.parseValueType(c.Map_type().Value_type(), t)
			return MapType{
				Key:   kt,
				Value: vt,
			}
		} else if c.List_type() != nil {
			vt := l.parseValueType(c.List_type().Value_type(), t)
			return ListType{
				Item: vt,
			}
		}
	}

	panic(fmt.Errorf("invalid type %s in line %d", ctx.GetText(), ctx.GetStart().GetLine()))
}

// ExitRpc_def handles RPC definitions, including request/response
// types and annotations.
func (l *ParseTreeListener) ExitRpc_def(ctx *Rpc_defContext) {
	r := RPC{
		Name: ctx.IDENTIFIER().GetText(),
		Position: Position{
			Start: ctx.GetStart().GetLine(),
			Stop:  ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Above: l.aboveComment(ctx.GetStart()),
		},
	}
	if !IsPascal(r.Name) {
		panic(fmt.Errorf("RPC name %s is not PascalCase in line %d", r.Name, r.Position.Start))
	}

	// Request
	reqType := ctx.Rpc_req().User_type()
	r.Request = UserType{
		Name:     reqType.IDENTIFIER().GetText(),
		Optional: reqType.QUESTION() != nil,
	}
	if !IsPascal(r.Request.Name) {
		panic(fmt.Errorf("RPC request type %s is not PascalCase in line %d", r.Request.Name, r.Position.Start))
	}
	l.Document.UsedTypes[r.Request.Name] = struct{}{}

	// Response
	respType := ctx.Rpc_resp().User_type()
	if ctx.Rpc_resp().TYPE_STREAM() != nil {
		r.Response.Stream = true
	}
	r.Response.UserType = UserType{
		Name:     respType.IDENTIFIER().GetText(),
		Optional: respType.QUESTION() != nil,
	}
	if !IsPascal(r.Response.UserType.Name) {
		panic(fmt.Errorf("RPC response type %s is not PascalCase in line %d", r.Response.UserType.Name, r.Position.Start))
	}
	l.Document.UsedTypes[r.Response.UserType.Name] = struct{}{}

	// Annotations
	for _, aCtx := range ctx.Rpc_annotations().AllAnnotation() {
		a := Annotation{
			Key: aCtx.IDENTIFIER().GetText(),
			Position: Position{
				Start: aCtx.GetStart().GetLine(),
				Stop:  aCtx.GetStop().GetLine(),
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

	l.Document.RPCs = append(l.Document.RPCs, r)
}

// isTerminatorToken returns true if the given token is considered a statement terminator.
// In this parser, a newline or semicolon marks the end of a statement.
func isTerminatorToken(t antlr.Token) bool {
	return t.GetTokenType() == TLexerNEWLINE || t.GetTokenType() == TLexerSEMICOLON
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
				Start: c.GetLine(),
				Stop:  c.GetLine(),
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
				Start: c.GetLine(),
				Stop:  c.GetLine() + len(lines) - 1,
			},
		})
	}

	// Sort comments by starting line in descending order
	sort.Slice(all, func(i, j int) bool {
		return all[i].Position.Start >= all[j].Position.Start
	})

	// Select only the contiguous block directly above the token
	i := 0
	lastLine := token.GetLine()
	for ; i < len(all); i++ {
		c := all[i]
		if c.Position.Stop != lastLine-1 {
			break
		}
		ret = append([]Comment{c}, ret...)
		lastLine = c.Position.Start
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
				Start: c.GetLine(),
				Stop:  c.GetLine(),
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
				Start: c.GetLine(),
				Stop:  c.GetLine() + len(lines) - 1,
			},
		}
	}

	return nil
}
