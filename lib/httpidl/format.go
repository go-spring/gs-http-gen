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
	"sort"
	"strconv"
	"strings"
)

const indent = "    "

// docItemKind defines the type of items that can appear in the document.
// This is used to help with ordering and spacing when reconstructing the output.
type docItemKind int

const (
	// Single-line comment (e.g., // comment)
	docItemKindSLComment = docItemKind(iota)

	// Multi-line comment (e.g., /* ... */)
	docItemKindMLComment

	// Constant declaration
	docItemKindConst

	// Enum declaration
	docItemKindEnum

	// Type declaration (struct-like or oneof)
	docItemKindType

	// RPC (remote procedure call) declaration
	docItemKindRPC
)

// docItem represents one piece of a document (comment, const, type, etc.)
// with its source position and pre-rendered buffer content.
type docItem struct {
	kind docItemKind // Kind of the document item
	pos  int         // Position in the original source
	buf  string      // Rendered text of the document item
}

// Format converts a parsed Document AST back into a formatted string.
// It preserves the original order of items (using their source positions)
// and ensures consistent blank lines and comment placement.
func Format(doc Document) string {
	var items []docItem

	// Collect top-level standalone comments
	for _, c := range doc.Comments {
		kind := docItemKindMLComment
		if c.Single {
			kind = docItemKindSLComment
		}
		items = append(items, docItem{
			kind: kind,
			pos:  c.Position.StartLine,
			buf:  strings.Join(c.Text, "\n"),
		})
	}

	// Collect constants
	for _, c := range doc.Consts {
		items = append(items, docItem{
			kind: docItemKindConst,
			pos:  c.Position.StartLine,
			buf:  formatConst(c),
		})
	}

	// Collect enums
	for _, e := range doc.Enums {
		if e.OneOf {
			continue
		}
		items = append(items, docItem{
			kind: docItemKindEnum,
			pos:  e.Position.StartLine,
			buf:  formatEnum(e),
		})
	}

	// Collect types
	for _, t := range doc.Types {
		items = append(items, docItem{
			kind: docItemKindType,
			pos:  t.Position.StartLine,
			buf:  formatType(t),
		})
	}

	// Collect RPCs
	for _, r := range doc.RPCs {
		items = append(items, docItem{
			kind: docItemKindRPC,
			pos:  r.Position.StartLine,
			buf:  formatRPC(r),
		})
	}

	// Sort all collected items by their original source position
	sort.Slice(items, func(i, j int) bool {
		return items[i].pos < items[j].pos
	})

	// Render items with proper spacing depending on their type
	var sb strings.Builder
	lastKind := docItemKindSLComment
	for i, item := range items {
		switch lastKind {
		case docItemKindEnum, docItemKindType, docItemKindRPC:
			sb.WriteString("\n")
		default:
			if i > 0 && (lastKind != item.kind || item.kind == docItemKindMLComment) {
				sb.WriteString("\n")
			}
		}
		sb.WriteString(item.buf)
		sb.WriteString("\n")
		lastKind = item.kind
	}
	return sb.String()
}

// formatAboveComments writes comments that appear above a declaration or field.
// The prefix is typically indentation (e.g., indent for struct fields).
func formatAboveComments(comments []Comment, sb *strings.Builder, prefix string) {
	for _, c := range comments {
		// Single-line comment
		if c.Single {
			sb.WriteString(prefix)
			sb.WriteString(c.Text[0])
			sb.WriteString("\n")
			continue
		}

		// Multi-line comment
		for _, s := range c.Text {
			sb.WriteString(prefix)
			sb.WriteString(s)
			sb.WriteString("\n")
		}
	}
}

// formatRightComment writes a comment that appears on the same line (or to the right)
// of a declaration. Multi-line comments continue on subsequent lines with the given prefix.
func formatRightComment(c *Comment, sb *strings.Builder, prefix string) {
	if c == nil {
		return
	}

	// Inline single-line comment
	if c.Single {
		sb.WriteString(" ")
		sb.WriteString(c.Text[0])
		return
	}

	// Multi-line right-side comment
	for i, s := range c.Text {
		if i == 0 {
			sb.WriteString(" ")
		} else {
			sb.WriteString("\n")
			sb.WriteString(prefix)
		}
		sb.WriteString(s)
	}
}

// formatConst formats a constant declaration, including its leading (above) comments
// and any inline (right-side) comment.
func formatConst(c Const) string {
	var sb strings.Builder
	formatAboveComments(c.Comments.Above, &sb, "")

	sb.WriteString("const ")
	sb.WriteString(c.Type.Name)
	sb.WriteString(" ")
	sb.WriteString(c.Name)
	sb.WriteString(" = ")
	sb.WriteString(c.Value)

	formatRightComment(c.Comments.Right, &sb, "")
	return sb.String()
}

// formatEnum formats an enum declaration and its fields,
// preserving top-level and per-field comments.
func formatEnum(e Enum) string {
	var sb strings.Builder
	formatAboveComments(e.Comments.Above, &sb, "")

	sb.WriteString("enum ")
	sb.WriteString(e.Name)
	sb.WriteString(" {")

	for _, f := range e.Fields {
		sb.WriteString("\n")
		formatAboveComments(f.Comments.Above, &sb, indent)

		sb.WriteString(indent)
		sb.WriteString(f.Name)
		sb.WriteString(" = ")
		sb.WriteString(strconv.FormatInt(f.Value, 10))

		formatFieldAnnotations(f.Annotations, &sb)
		formatRightComment(f.Comments.Right, &sb, indent)
	}

	sb.WriteString("\n}")
	return sb.String()
}

// formatType formats a type (or oneof) declaration, including its generic
// parameters, fields, comments, and potential redefinition.
func formatType(t Type) string {
	var sb strings.Builder
	formatAboveComments(t.Comments.Above, &sb, "")

	if t.OneOf {
		sb.WriteString("oneof ")
	} else {
		sb.WriteString("type ")
	}

	sb.WriteString(t.Name)

	if t.InstType != nil {
		sb.WriteString(" ")
		sb.WriteString(t.InstType.Text())
	} else {
		// If the type has generic parameter(s)
		if t.GenericParam != nil {
			sb.WriteString("<")
			sb.WriteString(*t.GenericParam)
			sb.WriteString(">")
		}

		sb.WriteString(" {")
		for i, f := range t.RawFields {
			if i == 0 && t.OneOf {
				continue
			}
			sb.WriteString("\n")
			formatTypeField(t, f, &sb)
		}
		sb.WriteString("\n}")
	}
	return sb.String()
}

// formatTypeField formats a single field in a type declaration,
// including its type, name, annotations, and comments.
func formatTypeField(t Type, f TypeField, sb *strings.Builder) {
	formatAboveComments(f.Comments.Above, sb, indent)

	sb.WriteString(indent)
	if f.Required {
		sb.WriteString("required ")
	}
	sb.WriteString(f.Type.Text())

	if _, ok := f.Type.(EmbedType); !ok && !t.OneOf {
		sb.WriteString(" ")
		sb.WriteString(f.Name)
	}

	formatFieldAnnotations(f.Annotations, sb)
	formatRightComment(f.Comments.Right, sb, indent)
}

// formatFieldAnnotations formats a field's annotations.
func formatFieldAnnotations(arr []Annotation, sb *strings.Builder) {
	if len(arr) > 0 {
		sb.WriteString(" (")
		for i, a := range arr {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(" ")
			sb.WriteString(a.Key)
			if a.Value != nil {
				sb.WriteString("=")
				sb.WriteString(*a.Value)
			}
		}
		sb.WriteString(" )")
	}
}

// formatRPC formats an RPC declaration including its request/response types,
// annotations, and associated comments.
func formatRPC(r RPC) string {
	var sb strings.Builder
	formatAboveComments(r.Comments.Above, &sb, "")

	if r.SSE {
		sb.WriteString("sse ")
	} else {
		sb.WriteString("rpc ")
	}

	sb.WriteString(r.Name)
	sb.WriteString("(")
	sb.WriteString(r.Request)
	sb.WriteString(") ")
	sb.WriteString(r.Response.Text())
	sb.WriteString(" {")

	// todo group annotations
	for _, a := range r.Annotations {
		sb.WriteString("\n")
		formatAboveComments(a.Comments.Above, &sb, indent)

		sb.WriteString(indent)
		sb.WriteString(a.Key)
		if a.Value != nil {
			sb.WriteString("=")
			sb.WriteString(*a.Value)
		}

		formatRightComment(a.Comments.Right, &sb, indent)
	}

	sb.WriteString("\n}")
	return sb.String()
}
