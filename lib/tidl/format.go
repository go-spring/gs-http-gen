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

// Format parses the input source string into a Document AST
// and returns its normalized pretty-printed representation.
func Format(data []byte) (string, error) {
	doc, err := Parse(data)
	if err != nil {
		return "", err
	}
	return Dump(doc), nil
}

// Dump converts a parsed Document AST back into a formatted string.
// It preserves the original order of items (using their source positions)
// and ensures consistent blank lines and comment placement.
func Dump(doc Document) string {
	var items []docItem

	// Collect top-level standalone comments
	for _, c := range doc.Comments {
		kind := docItemKindMLComment
		if c.Single {
			kind = docItemKindSLComment
		}
		items = append(items, docItem{
			kind: kind,
			pos:  c.Position.Start,
			buf:  strings.Join(c.Text, "\n"),
		})
	}

	// Collect constants
	for _, c := range doc.Consts {
		items = append(items, docItem{
			kind: docItemKindConst,
			pos:  c.Position.Start,
			buf:  dumpConst(c),
		})
	}

	// Collect enums
	for _, e := range doc.Enums {
		items = append(items, docItem{
			kind: docItemKindEnum,
			pos:  e.Position.Start,
			buf:  dumpEnum(e),
		})
	}

	// Collect types
	for _, t := range doc.Types {
		items = append(items, docItem{
			kind: docItemKindType,
			pos:  t.Position.Start,
			buf:  dumpType(t),
		})
	}

	// Collect RPCs
	for _, r := range doc.RPCs {
		items = append(items, docItem{
			kind: docItemKindRPC,
			pos:  r.Position.Start,
			buf:  dumpRPC(r),
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

// dumpAboveComments writes comments that appear above a declaration or field.
// The prefix is typically indentation (e.g., indent for struct fields).
func dumpAboveComments(comments []Comment, sb *strings.Builder, prefix string) {
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

// dumpRightComment writes a comment that appears on the same line (or to the right)
// of a declaration. Multi-line comments continue on subsequent lines with the given prefix.
func dumpRightComment(c *Comment, sb *strings.Builder, prefix string) {
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

// dumpConst formats a constant declaration, including its leading (above) comments
// and any inline (right-side) comment.
func dumpConst(c Const) string {
	var sb strings.Builder
	dumpAboveComments(c.Comments.Above, &sb, "")

	sb.WriteString("const ")
	sb.WriteString(c.Type)
	sb.WriteString(" ")
	sb.WriteString(c.Name)
	sb.WriteString(" = ")
	sb.WriteString(c.Value)

	dumpRightComment(c.Comments.Right, &sb, "")
	return sb.String()
}

// dumpEnum formats an enum declaration and its fields,
// preserving top-level and per-field comments.
func dumpEnum(e Enum) string {
	var sb strings.Builder
	dumpAboveComments(e.Comments.Above, &sb, "")

	sb.WriteString("enum ")
	sb.WriteString(e.Name)
	sb.WriteString(" {")

	for _, f := range e.Fields {
		sb.WriteString("\n")
		dumpAboveComments(f.Comments.Above, &sb, indent)

		sb.WriteString(indent)
		sb.WriteString(f.Name)
		sb.WriteString(" = ")
		sb.WriteString(strconv.FormatInt(f.Value, 10))

		dumpRightComment(f.Comments.Right, &sb, indent)
	}

	sb.WriteString("\n}")
	return sb.String()
}

// dumpType formats a type (or oneof) declaration, including its generic
// parameters, fields, comments, and potential redefinition.
func dumpType(t Type) string {
	var sb strings.Builder
	dumpAboveComments(t.Comments.Above, &sb, "")

	if t.OneOf {
		sb.WriteString("oneof ")
	} else {
		sb.WriteString("type ")
	}

	sb.WriteString(t.Name)

	if t.Redefined != nil {
		sb.WriteString(" ")
		sb.WriteString(t.Redefined.Text())
	} else {
		// If the type has generic parameter(s)
		if t.GenericName != nil {
			sb.WriteString("<")
			sb.WriteString(*t.GenericName)
			sb.WriteString(">")
		}

		sb.WriteString(" {")
		for _, f := range t.Fields {
			sb.WriteString("\n")
			dumpTypeField(f, &sb)
		}
		sb.WriteString("\n}")
	}
	return sb.String()
}

// dumpTypeField formats a single field in a type declaration,
// including its type, name, default value, annotations, and comments.
func dumpTypeField(f TypeField, sb *strings.Builder) {
	dumpAboveComments(f.Comments.Above, sb, indent)

	sb.WriteString(indent)
	sb.WriteString(f.FieldType.Text())

	if _, ok := f.FieldType.(EmbedType); !ok {
		sb.WriteString(" ")
		sb.WriteString(f.Name)

		// Default value
		if f.Default != nil {
			sb.WriteString(" = ")
			sb.WriteString(*f.Default)
		}

		// Annotations
		if len(f.Annotations) > 0 {
			sb.WriteString(" (")
			for i, a := range f.Annotations {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(" ")
				sb.WriteString(a.Key)
				if a.Value != nil {
					sb.WriteString(" = ")
					sb.WriteString(*a.Value)
				}
			}
			sb.WriteString(" )")
		}
	}

	dumpRightComment(f.Comments.Right, sb, indent)
}

// dumpRPC formats an RPC declaration including its request/response types,
// annotations, and associated comments.
func dumpRPC(r RPC) string {
	var sb strings.Builder
	dumpAboveComments(r.Comments.Above, &sb, "")

	sb.WriteString("rpc ")
	sb.WriteString(r.Name)
	sb.WriteString("(")
	sb.WriteString(r.Request)
	sb.WriteString(") ")
	sb.WriteString(r.Response.Text())
	sb.WriteString(" {")

	for _, a := range r.Annotations {
		sb.WriteString("\n")
		dumpAboveComments(a.Comments.Above, &sb, indent)

		sb.WriteString(indent)
		sb.WriteString(a.Key)
		if a.Value != nil {
			sb.WriteString(" = ")
			sb.WriteString(*a.Value)
		}

		dumpRightComment(a.Comments.Right, &sb, indent)
	}

	sb.WriteString("\n}")
	return sb.String()
}
