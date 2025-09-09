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

package vidl

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// Expr represents a generic expression node.
type Expr interface {
	Text() string
}

// BinaryExpr represents a binary expression (e.g., a && b, x == y).
type BinaryExpr struct {
	Left  Expr   // Left-hand side expression
	Op    string // Operator (e.g., &&, ||, ==, <)
	Right Expr   // Right-hand side expression
}

func (e BinaryExpr) Text() string {
	if e.Left == nil {
		return ""
	}
	if e.Right == nil {
		return e.Left.Text()
	}
	return fmt.Sprintf("%s %s %s", e.Left.Text(), e.Op, e.Right.Text())
}

// UnaryExpr represents a unary expression (e.g., !x).
type UnaryExpr struct {
	Op   string // Operator (e.g., !)
	Expr Expr   // Operand expression
}

func (e UnaryExpr) Text() string {
	if e.Expr == nil {
		return ""
	}
	return fmt.Sprintf("%s%s", e.Op, e.Expr.Text())
}

// PrimaryExpr represents an atomic expression, which can be a literal,
// identifier, function call, or parenthesized expression.
type PrimaryExpr struct {
	Value string     // Literal value or identifier
	Call  *FuncCall  // Optional function call
	Inner *InnerExpr // Optional parenthesized expression
}

func (e PrimaryExpr) Text() string {
	if e.Inner != nil {
		return e.Inner.Text()
	}
	if e.Call != nil {
		return e.Call.Text()
	}
	return e.Value
}

// FuncCall represents a function call expression with arguments.
type FuncCall struct {
	Name string // Function name
	Args []Expr // Arguments
}

func (f FuncCall) Text() string {
	if len(f.Args) == 0 {
		return f.Name + "()"
	}
	var args []string
	for _, arg := range f.Args {
		args = append(args, arg.Text())
	}
	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ", "))
}

// InnerExpr represents a parenthesized expression.
type InnerExpr struct {
	Expr Expr
}

func (e InnerExpr) Text() string {
	if e.Expr == nil {
		return ""
	}
	return fmt.Sprintf("(%s)", e.Expr.Text())
}

// Parse parses the input string and returns an Expr AST.
func Parse(s string) (expr Expr, err error) {
	e := &ErrorListener{}

	// Recover from panics to return a proper error
	defer func() {
		if r := recover(); r != nil {
			expr = nil
			err = fmt.Errorf("[PANIC]: %v\n%s", r, debug.Stack())
			if e.err != nil {
				err = fmt.Errorf("%w\n%w", e.err, err)
			}
		}
	}()

	// Step 1: Create lexer and token stream
	input := antlr.NewInputStream(s)
	lexer := NewVLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(e)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Step 2: Create parser and attach custom error listener
	p := NewVParser(tokens)
	p.RemoveErrorListeners()
	p.AddErrorListener(e)
	p.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL) // Use faster SLL mode

	// Step 3: Walk parse tree with custom listener
	l := &ParseTreeListener{
		Tokens: tokens,
	}
	antlr.ParseTreeWalkerDefault.Walk(l, p.ValidateExpr())

	// Step 4: Return parsed expression or error
	if e.err != nil {
		return nil, e.err
	}
	return l.Expr, nil
}

// ErrorListener implements a custom ANTLR error listener that records syntax errors.
type ErrorListener struct {
	*antlr.DefaultErrorListener
	err error
}

// SyntaxError is called by ANTLR when a syntax error occurs.
func (l *ErrorListener) SyntaxError(_ antlr.Recognizer, _ any, line, column int, msg string, e antlr.RecognitionException) {
	if l.err == nil {
		l.err = fmt.Errorf("line %d:%d %s", line, column, msg)
		return
	}
	l.err = fmt.Errorf("%w\nline %d:%d %s", l.err, line, column, msg)
}

// ParseTreeListener walks the parse tree and constructs the expression AST.
type ParseTreeListener struct {
	BaseVParserListener
	Tokens *antlr.CommonTokenStream
	Expr   Expr
}

func (l *ParseTreeListener) ExitValidateExpr(ctx *ValidateExprContext) {
	l.Expr = parseValidateExpr(ctx)
}

// parseValidateExpr converts a ValidateExprContext into an Expr.
func parseValidateExpr(ctx IValidateExprContext) Expr {
	if ctx.LogicalOrExpr() == nil {
		return nil
	}
	return parseLogicalOrExpr(ctx.LogicalOrExpr())
}

// parseLogicalOrExpr converts a LogicalOrExprContext into an Expr.
func parseLogicalOrExpr(ctx ILogicalOrExprContext) Expr {
	if ctx.LOGICAL_OR() != nil {
		return BinaryExpr{
			Left:  parseLogicalAndExpr(ctx.LogicalAndExpr(0)),
			Op:    ctx.LOGICAL_OR().GetText(),
			Right: parseLogicalAndExpr(ctx.LogicalAndExpr(1)),
		}
	}
	return parseLogicalAndExpr(ctx.LogicalAndExpr(0))
}

// parseLogicalAndExpr converts a LogicalAndExprContext into an Expr.
func parseLogicalAndExpr(ctx ILogicalAndExprContext) Expr {
	if ctx.LOGICAL_AND() != nil {
		return BinaryExpr{
			Left:  parseEqualityExpr(ctx.EqualityExpr(0)),
			Op:    ctx.LOGICAL_AND().GetText(),
			Right: parseEqualityExpr(ctx.EqualityExpr(1)),
		}
	}
	return parseEqualityExpr(ctx.EqualityExpr(0))
}

// parseEqualityExpr converts an EqualityExprContext into an Expr.
func parseEqualityExpr(ctx IEqualityExprContext) Expr {
	var op antlr.TerminalNode
	if ctx.EQUAL() != nil {
		op = ctx.EQUAL()
	} else if ctx.NOT_EQUAL() != nil {
		op = ctx.NOT_EQUAL()
	}
	if op != nil {
		return BinaryExpr{
			Left:  parseRelationalExpr(ctx.RelationalExpr(0)),
			Op:    op.GetText(),
			Right: parseRelationalExpr(ctx.RelationalExpr(1)),
		}
	}
	return parseRelationalExpr(ctx.RelationalExpr(0))
}

// parseRelationalExpr converts a RelationalExprContext into an Expr.
func parseRelationalExpr(ctx IRelationalExprContext) Expr {
	var op antlr.TerminalNode
	if ctx.LESS_THAN() != nil {
		op = ctx.LESS_THAN()
	} else if ctx.LESS_OR_EQUAL() != nil {
		op = ctx.LESS_OR_EQUAL()
	} else if ctx.GREATER_THAN() != nil {
		op = ctx.GREATER_THAN()
	} else if ctx.GREATER_OR_EQUAL() != nil {
		op = ctx.GREATER_OR_EQUAL()
	}
	if op != nil {
		return BinaryExpr{
			Left:  parseUnaryExpr(ctx.UnaryExpr(0)),
			Op:    op.GetText(),
			Right: parseUnaryExpr(ctx.UnaryExpr(1)),
		}
	}
	return parseUnaryExpr(ctx.UnaryExpr(0))
}

// parseUnaryExpr converts a UnaryExprContext into an Expr.
func parseUnaryExpr(ctx IUnaryExprContext) Expr {
	if ctx.LOGICAL_NOT() != nil {
		return UnaryExpr{
			Op:   ctx.LOGICAL_NOT().GetText(),
			Expr: parseUnaryExpr(ctx.UnaryExpr()),
		}
	}
	return parsePrimaryExpr(ctx.PrimaryExpr())
}

// parsePrimaryExpr converts a PrimaryExprContext into an Expr.
func parsePrimaryExpr(ctx IPrimaryExprContext) Expr {
	if ctx == nil {
		return nil
	}
	if ctx.IDENTIFIER() != nil {
		return PrimaryExpr{
			Value: ctx.IDENTIFIER().GetText(),
		}
	}
	if ctx.KW_DOLLAR() != nil {
		return PrimaryExpr{
			Value: ctx.KW_DOLLAR().GetText(),
		}
	}
	if ctx.KW_NIL() != nil {
		return PrimaryExpr{
			Value: ctx.KW_NIL().GetText(),
		}
	}
	if ctx.INTEGER() != nil {
		return PrimaryExpr{
			Value: ctx.INTEGER().GetText(),
		}
	}
	if ctx.FLOAT() != nil {
		return PrimaryExpr{
			Value: ctx.FLOAT().GetText(),
		}
	}
	if ctx.STRING() != nil {
		return PrimaryExpr{
			Value: ctx.STRING().GetText(),
		}
	}
	if ctx.FunctionCall() != nil {
		return PrimaryExpr{
			Call: parseFunctionCall(ctx.FunctionCall()),
		}
	}
	if ctx.LEFT_PAREN() != nil {
		return PrimaryExpr{
			Inner: &InnerExpr{
				Expr: parseValidateExpr(ctx.ValidateExpr()),
			},
		}
	}
	return nil
}

// parseFunctionCall converts a FunctionCallContext into a FuncCall AST node.
func parseFunctionCall(ctx IFunctionCallContext) *FuncCall {
	var args []Expr
	for _, arg := range ctx.AllValidateExpr() {
		args = append(args, parseValidateExpr(arg))
	}
	return &FuncCall{
		Name: ctx.IDENTIFIER().GetText(),
		Args: args,
	}
}
