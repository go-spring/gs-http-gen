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
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// ParseMeta parses meta information from the given input string.
func ParseMeta(s string) (*MetaInfo, error) {
	r := &MetaInfo{}
	if err := json.Unmarshal([]byte(s), r); err != nil {
		return nil, err
	}
	return r, nil
}

// Parse runs the parsing pipeline for the given input string.
func Parse(s string) (doc Document, err error) {
	e := &ErrorListener{}

	defer func() {
		if r := recover(); r != nil {
			doc = Document{}
			err = fmt.Errorf("[PANIC]: %v\n%s", r, debug.Stack())
			if e.err != nil {
				err = fmt.Errorf("%w\n%w", e.err, err)
			}
		}
	}()

	// Step 1. Create lexer and token stream
	input := antlr.NewInputStream(s + "\n")
	lexer := NewTLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(e)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Step 2. Create parser and attach custom error listener
	p := NewTParser(tokens)
	p.RemoveErrorListeners()
	p.AddErrorListener(e)

	// Use SLL mode first (faster, may fall back to LL if needed).
	p.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)

	// Step 3. Walk the parse tree with a custom listener
	l := &ParseTreeListener{
		Tokens:   tokens,
		Attached: make(map[int]struct{}),
	}
	antlr.ParseTreeWalkerDefault.Walk(l, p.Document())

	// Step 4. Return results or error
	if e.err != nil {
		return Document{}, e.err
	}
	return l.Document, nil
}

// ErrorListener implements a custom ANTLR error listener.
type ErrorListener struct {
	*antlr.DefaultErrorListener
	err error
}

// SyntaxError is called by ANTLR when a syntax error is encountered.
func (l *ErrorListener) SyntaxError(_ antlr.Recognizer, _ any, line, column int, msg string, e antlr.RecognitionException) {
	if l.err == nil {
		l.err = fmt.Errorf("line %d:%d %s", line, column, msg)
		return
	}
	l.err = fmt.Errorf("%w\nline %d:%d %s", l.err, line, column, msg)
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
			Top:   l.topComment(ctx.GetStart()),
			Right: l.rightComment(ctx.GetStop()),
		},
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
			Top: l.topComment(ctx.GetStart()),
		},
	}

	for _, f := range ctx.AllEnum_field() {
		v, err := strconv.ParseInt(f.INTEGER().GetText(), 0, 64)
		if err != nil {
			panic(fmt.Errorf("parse enum value on line %d error: %w", f.GetStart().GetLine(), err))
		}
		e.Fields = append(e.Fields, EnumField{
			Name:  f.IDENTIFIER().GetText(),
			Value: v,
			Position: Position{
				Start: f.GetStart().GetLine(),
				Stop:  f.GetStop().GetLine(),
			},
			Comments: Comments{
				Top:   l.topComment(f.GetStart()),
				Right: l.rightComment(f.GetStop()),
			},
		})
	}
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
			Top: l.topComment(ctx.GetStart()),
		},
	}

	if ctx.LEFT_BRACE() != nil {
		l.parseCompleteType(ctx, &t)
	} else {
		l.parseRedefinedType(ctx, &t)
	}

	l.Document.Types = append(l.Document.Types, t)
}

func (l *ParseTreeListener) parseCompleteType(ctx *Type_defContext, t *Type) {

	// Handle generic type parameter (if any)
	if ctx.LESS_THAN() != nil {
		s := ctx.IDENTIFIER(1).GetText()
		t.GenericName = &s
	}

	// Process all type fields
	for _, f := range ctx.AllType_field() {
		typeField := TypeField{
			Position: Position{
				Start: f.GetStart().GetLine(),
				Stop:  f.GetStop().GetLine(),
			},
			Comments: Comments{
				Top:   l.topComment(f.GetStart()),
				Right: l.rightComment(f.GetStop()),
			},
		}

		// Distinguish between embedded fields and normal fields
		if f.Embed_type_field() != nil {
			typeField.FieldType = EmbedType{
				Name:     f.Embed_type_field().User_type().IDENTIFIER().GetText(),
				Optional: f.Embed_type_field().User_type().QUESTION() != nil,
			}
		} else if f.Common_type_field() != nil {
			// Regular field
			typeField.FieldType = l.parseCommonFieldType(f.Common_type_field().Common_field_type())
			typeField.Name = f.Common_type_field().IDENTIFIER().GetText()

			// Default value
			if f.Common_type_field().Const_value() != nil {
				s := f.Common_type_field().Const_value().GetText()
				typeField.Default = &s
			}

			// Annotations
			if f.Common_type_field().Type_annotations() != nil {
				for _, aCtx := range f.Common_type_field().Type_annotations().AllAnnotation() {
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

		t.Fields = append(t.Fields, typeField)
	}
}

// parseRedefinedType handles redefined types, including generic types.
func (l *ParseTreeListener) parseRedefinedType(ctx *Type_defContext, t *Type) {
	t.Redefined = &RedefinedType{
		Name: ctx.IDENTIFIER(1).GetText(),
	}
	g := ctx.Value_type()
	if g.Base_type() != nil {
		t.Redefined.GenericType = BaseType{
			Name:     strings.TrimRight(g.Base_type().GetText(), "?"),
			Optional: g.Base_type().QUESTION() != nil,
		}
	}
	if g.User_type() != nil {
		t.Redefined.GenericType = UserType{
			Name:     g.User_type().IDENTIFIER().GetText(),
			Optional: g.User_type().QUESTION() != nil,
		}
	}
	if g.Container_type() != nil {
		if g.Container_type().Map_type() != nil {
			kt := g.Container_type().Map_type().Key_type().GetText()
			vt := l.parseValueType(g.Container_type().Map_type().Value_type())
			t.Redefined.GenericType = MapType{
				Key:   kt,
				Value: vt,
			}
		} else if g.Container_type().List_type() != nil {
			vt := l.parseValueType(g.Container_type().List_type().Value_type())
			t.Redefined.GenericType = ListType{
				Item: vt,
			}
		}
	}
	if t.Redefined.GenericType != nil {
		return
	}
	panic(fmt.Errorf("unknown type: %s", g.GetText()))
}

// parseCommonFieldType resolves type definitions inside type fields.
// It distinguishes between built-in types, user-defined types, and containers.
func (l *ParseTreeListener) parseCommonFieldType(ctx ICommon_field_typeContext) TypeDefinition {
	if ctx.TYPE_ANY() != nil {
		return AnyType{}
	}
	if ctx.TYPE_BINARY() != nil {
		return BinaryType{}
	}
	if ctx.Base_type() != nil {
		return BaseType{
			Name:     strings.TrimRight(ctx.Base_type().GetText(), "?"),
			Optional: ctx.Base_type().QUESTION() != nil,
		}
	}
	if ctx.User_type() != nil {
		return UserType{
			Name:     ctx.User_type().IDENTIFIER().GetText(),
			Optional: ctx.User_type().QUESTION() != nil,
		}
	}
	if ctx.Container_type() != nil {
		if ctx.Container_type().Map_type() != nil {
			kt := ctx.Container_type().Map_type().Key_type().GetText()
			vt := l.parseValueType(ctx.Container_type().Map_type().Value_type())
			return MapType{
				Key:   kt,
				Value: vt,
			}
		} else if ctx.Container_type().List_type() != nil {
			vt := l.parseValueType(ctx.Container_type().List_type().Value_type())
			return ListType{
				Item: vt,
			}
		}
	}
	panic(fmt.Errorf("unknown type: %s", ctx.GetText()))
}

// parseValueType resolves value types inside container types.
func (l *ParseTreeListener) parseValueType(ctx IValue_typeContext) TypeDefinition {
	if ctx.Base_type() != nil {
		return BaseType{
			Name:     strings.TrimRight(ctx.Base_type().GetText(), "?"),
			Optional: ctx.Base_type().QUESTION() != nil,
		}
	}
	if ctx.User_type() != nil {
		return UserType{
			Name:     strings.TrimRight(ctx.User_type().IDENTIFIER().GetText(), "?"),
			Optional: ctx.User_type().QUESTION() != nil,
		}
	}
	if ctx.Container_type() != nil {
		if ctx.Container_type().Map_type() != nil {
			kt := ctx.Container_type().Map_type().Key_type().GetText()
			vt := l.parseValueType(ctx.Container_type().Map_type().Value_type())
			return MapType{
				Key:   kt,
				Value: vt,
			}
		} else if ctx.Container_type().List_type() != nil {
			vt := l.parseValueType(ctx.Container_type().List_type().Value_type())
			return ListType{
				Item: vt,
			}
		}
	}
	panic(fmt.Errorf("unknown type: %s", ctx.GetText()))
}

func (l *ParseTreeListener) ExitOneof_def(ctx *Oneof_defContext) {
	o := Type{
		Name:  ctx.IDENTIFIER().GetText(),
		OneOf: true,
		Position: Position{
			Start: ctx.GetStart().GetLine(),
			Stop:  ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Top: l.topComment(ctx.GetStart()),
		},
	}

	l.parseOneOfType(ctx, &o)

	l.Document.Types = append(l.Document.Types, o)
}

// parseOneOfType handles oneof types, including fields and annotations.
func (l *ParseTreeListener) parseOneOfType(ctx *Oneof_defContext, o *Type) {

	// Process all oneof fields
	for _, f := range ctx.AllCommon_type_field() {
		typeField := TypeField{
			Position: Position{
				Start: f.GetStart().GetLine(),
				Stop:  f.GetStop().GetLine(),
			},
			Comments: Comments{
				Top:   l.topComment(f.GetStart()),
				Right: l.rightComment(f.GetStop()),
			},
		}

		// Regular field
		typeField.FieldType = l.parseCommonFieldType(f.Common_field_type())
		typeField.Name = f.IDENTIFIER().GetText()

		// Default value
		if f.Const_value() != nil {
			s := f.Const_value().GetText()
			typeField.Default = &s
		}

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

		o.Fields = append(o.Fields, typeField)
	}
}

// ExitRpc_def handles RPC definitions, including request/response
// types and annotations.
func (l *ParseTreeListener) ExitRpc_def(ctx *Rpc_defContext) {
	r := RPC{
		Name:    ctx.IDENTIFIER().GetText(),
		Request: ctx.Rpc_req().GetText(),
		Position: Position{
			Start: ctx.GetStart().GetLine(),
			Stop:  ctx.GetStop().GetLine(),
		},
		Comments: Comments{
			Top: l.topComment(ctx.GetStart()),
		},
	}

	// Handle response type
	if ctx.Rpc_resp().TYPE_STREAM() != nil {
		u := ctx.Rpc_resp().User_type()
		r.Response = RespType{
			Stream: true,
			UserType: &UserType{
				Name:     u.IDENTIFIER().GetText(),
				Optional: u.QUESTION() != nil,
			},
		}
	} else {
		r.Response = RespType{
			Stream:   false,
			TypeName: ctx.Rpc_resp().IDENTIFIER().GetText(),
		}
	}

	// Parse annotations
	for _, aCtx := range ctx.Rpc_annotations().AllAnnotation() {
		a := Annotation{
			Key: aCtx.IDENTIFIER().GetText(),
			Position: Position{
				Start: aCtx.GetStart().GetLine(),
				Stop:  aCtx.GetStop().GetLine(),
			},
			Comments: Comments{
				Top:   l.topComment(aCtx.GetStart()),
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

// isTerminatorToken returns true if the token is a terminator token.
func isTerminatorToken(t antlr.Token) bool {
	return t.GetTokenType() == TLexerNEWLINE || t.GetTokenType() == TLexerSEMI

}

// previousTokenOnChannel returns the previous token on the specified channel.
func (l *ParseTreeListener) previousTokenOnChannel(i int) int {
	tokens := l.Tokens.GetAllTokens()
	for i >= 0 && (isTerminatorToken(tokens[i]) || tokens[i].GetChannel() != antlr.LexerDefaultTokenChannel) {
		i--
	}
	return i
}

// filterForChannel filters tokens for a specific channel.
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

// GetHiddenTokensToLeft returns all hidden tokens to the left of a token.
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

// nextTokenOnChannel returns the next token on the specified channel.
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

// GetHiddenTokensToRight returns all hidden tokens to the right of a token.
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

// topComment extracts comments immediately above a token.
// It supports both single-line (//) and multi-line (/* */) comments.
func (l *ParseTreeListener) topComment(token antlr.Token) []Comment {
	var (
		all []Comment
		ret []Comment
	)

	// Collect single-line comments
	comments := l.GetHiddenTokensToLeft(token.GetTokenIndex(), TLexerSL_COMMENT_CHAN)
	for _, c := range comments {
		if _, ok := l.Attached[c.GetLine()]; ok {
			continue
		}
		all = append(all, Comment{
			Text:   strings.TrimSpace(c.GetText()),
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
		s := strings.TrimSpace(c.GetText())
		ss := strings.Split(s, "\n")
		for i := range ss {
			ss[i] = " " + strings.TrimSpace(ss[i])
		}
		ss[0] = strings.TrimSpace(ss[0])
		s = strings.Join(ss, "\n")
		count := len(ss)
		all = append(all, Comment{
			Text:   s,
			Single: false,
			Position: Position{
				Start: c.GetLine(),
				Stop:  c.GetLine() + count - 1,
			},
		})
	}

	// Sort comments by starting line (descending)
	sort.Slice(all, func(i, j int) bool {
		return all[i].Position.Start >= all[j].Position.Start
	})

	// Select only the contiguous block of comments directly above token
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

	// Remaining comments (not directly attached) go into Document.Comments
	for j := len(all) - 1; j >= i; j-- {
		l.Document.Comments = append(l.Document.Comments, all[j])
	}

	return ret
}

// rightComment extracts comments that appear at the end of the same line
// as a given token.
func (l *ParseTreeListener) rightComment(token antlr.Token) *Comment {
	// Single-line comments
	comments := l.Tokens.GetHiddenTokensToRight(token.GetTokenIndex(), TLexerSL_COMMENT_CHAN)
	for _, c := range comments {
		if c.GetLine() != token.GetLine() {
			continue
		}
		l.Attached[c.GetLine()] = struct{}{}
		return &Comment{
			Text:   strings.TrimSpace(c.GetText()),
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
		s := strings.TrimSpace(c.GetText())
		ss := strings.Split(s, "\n")
		for i := range ss {
			ss[i] = " " + strings.TrimSpace(ss[i])
		}
		ss[0] = strings.TrimSpace(ss[0])
		s = strings.Join(ss, "\n")
		count := len(ss)
		l.Attached[c.GetLine()] = struct{}{}
		return &Comment{
			Text:   s,
			Single: false,
			Position: Position{
				Start: c.GetLine(),
				Stop:  c.GetLine() + count - 1,
			},
		}
	}

	return nil
}
