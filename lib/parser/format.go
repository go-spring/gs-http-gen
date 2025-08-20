package parser

import (
	"sort"
	"strconv"
	"strings"
)

// docItemKind defines the type of items that can appear in the document.
// This is used to help with ordering and spacing when reconstructing the output.
type docItemKind int

const (
	// Single-line comment (e.g., // comment)
	docItemKindSLComment = docItemKind(iota)

	// Multi-line comment (e.g., /* ... */)
	docItemKindMLComment

	// Top-level comment attached to a const declaration
	docItemKindConstTop

	// Constant declaration
	docItemKindConst

	// Enum declaration
	docItemKindEnum

	// Type declaration (struct-like)
	docItemKindType

	// RPC (remote procedure call) declaration
	docItemKindRPC
)

// docItem represents one piece of a document (comment, const, type, etc.)
// with its source position and pre-rendered buffer content.
type docItem struct {
	kind docItemKind // The kind of document item
	pos  int         // Position in the original source
	buf  string      // The rendered text for this item
}

// Format parses the input string into a Document and then
// pretty-prints (dumps) it back into a normalized string form.
func Format(s string) (string, error) {
	doc, err := Parse(s)
	if err != nil {
		return "", err
	}
	return Dump(doc), nil
}

// Dump takes a Document AST and converts it into a formatted string.
// It preserves ordering and attaches comments correctly.
func Dump(doc Document) string {
	var sb strings.Builder
	dumpDocument(doc, &sb)
	return sb.String()
}

// dumpDocument collects all items (comments, consts, enums, types, RPCs),
// sorts them by position, and then writes them into the builder with proper spacing.
func dumpDocument(doc Document, sb *strings.Builder) {
	var items []docItem

	// Process standalone comments
	for _, c := range doc.Comments {
		kind := docItemKindMLComment
		if c.Single {
			kind = docItemKindSLComment
		}
		items = append(items, docItem{
			kind: kind,
			pos:  c.Position.Start,
			buf:  c.Text,
		})
	}

	// Process constants
	for _, c := range doc.Consts {
		items = append(items, docItem{
			kind: docItemKindConst,
			pos:  c.Position.Start,
			buf:  dumpConst(c),
		})
	}

	// Process enums
	for _, e := range doc.Enums {
		items = append(items, docItem{
			kind: docItemKindEnum,
			pos:  e.Position.Start,
			buf:  dumpEnum(e),
		})
	}

	// Process types
	for _, t := range doc.Types {
		items = append(items, docItem{
			kind: docItemKindType,
			pos:  t.Position.Start,
			buf:  dumpType(t),
		})
	}

	// Process RPCs
	for _, r := range doc.RPCs {
		items = append(items, docItem{
			kind: docItemKindRPC,
			pos:  r.Position.Start,
			buf:  dumpRPC(r),
		})
	}

	// Sort items by source position so they are written in the original order
	sort.Slice(items, func(i, j int) bool {
		return items[i].pos < items[j].pos
	})

	// Render items with proper spacing depending on their type
	lastKind := docItemKindSLComment
	for i, item := range items {
		switch lastKind {
		case docItemKindConstTop, docItemKindEnum, docItemKindType, docItemKindRPC:
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
}

// dumpConst converts a Const node into its textual representation,
// including its top comments and right-side comment if present.
func dumpConst(c *Const) string {
	var sb strings.Builder
	for _, s := range c.Comments.Top {
		sb.WriteString(s.Text)
		sb.WriteString("\n")
	}
	sb.WriteString("const ")
	sb.WriteString(c.Type)
	sb.WriteString(" ")
	sb.WriteString(c.Name)
	sb.WriteString(" = ")
	sb.WriteString(c.Value)
	if c.Comments.Right != nil {
		sb.WriteString(" ")
		sb.WriteString(c.Comments.Right.Text)
	}
	return sb.String()
}

// dumpEnum converts an Enum node into its textual representation,
// including field comments and inline comments.
func dumpEnum(e *Enum) string {
	var sb strings.Builder
	for _, s := range e.Comments.Top {
		sb.WriteString(s.Text)
		sb.WriteString("\n")
	}
	sb.WriteString("enum ")
	sb.WriteString(e.Name)
	sb.WriteString(" {")
	for _, f := range e.Fields {
		for _, s := range f.Comments.Top {
			sb.WriteString(s.Text)
			sb.WriteString("\n")
		}
		sb.WriteString("\n\t")
		sb.WriteString(f.Name)
		sb.WriteString(" = ")
		sb.WriteString(strconv.FormatInt(f.Value, 10))
		if f.Comments.Right != nil {
			sb.WriteString(" ")
			sb.WriteString(f.Comments.Right.Text)
		}
	}
	sb.WriteString("\n}")
	return sb.String()
}

// dumpType converts a Type node into its textual representation,
// including field definitions, generic parameter (if present),
// and top-level comments.
func dumpType(t *Type) string {
	var sb strings.Builder
	for _, s := range t.Comments.Top {
		sb.WriteString(s.Text)
		sb.WriteString("\n")
	}
	sb.WriteString("type ")
	sb.WriteString(t.Name)
	if t.Redefined != nil {
		sb.WriteString(" ")
		sb.WriteString(t.Redefined.Text())
	} else {
		if t.GenericName != nil {
			sb.WriteString("<")
			sb.WriteString(*t.GenericName)
			sb.WriteString(">")
		}
		sb.WriteString(" {")
		for _, f := range t.Fields {
			sb.WriteString("\n\t")
			dumpTypeField(f, &sb)
		}
		sb.WriteString("\n}")
	}
	return sb.String()
}

// dumpTypeField converts a TypeField node into textual representation,
// including annotations, default value, and comments.
func dumpTypeField(f TypeField, sb *strings.Builder) {
	for _, s := range f.Comments.Top {
		sb.WriteString(s.Text)
		sb.WriteString("\n\t")
	}
	sb.WriteString(f.FieldType.Text())
	if _, ok := f.FieldType.(EmbedType); !ok {
		sb.WriteString(" ")
		sb.WriteString(f.Name)
		if f.Default != nil {
			sb.WriteString(" = ")
			sb.WriteString(*f.Default)
		}
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
	if f.Comments.Right != nil {
		sb.WriteString(" ")
		sb.WriteString(f.Comments.Right.Text)
	}
}

// dumpRPC converts an RPC node into textual representation,
// including its annotations and comments.
func dumpRPC(r *RPC) string {
	var sb strings.Builder
	for _, s := range r.Comments.Top {
		sb.WriteString(s.Text)
		sb.WriteString("\n")
	}
	sb.WriteString("rpc ")
	sb.WriteString(r.Name)
	sb.WriteString("(")
	sb.WriteString(r.Request)
	sb.WriteString(") ")
	sb.WriteString(r.Response.Text())
	sb.WriteString(" {")
	for _, a := range r.Annotations {
		sb.WriteString("\n\t")
		for _, s := range a.Comments.Top {
			sb.WriteString("\n")
			sb.WriteString(s.Text)
		}
		sb.WriteString(a.Key)
		if a.Value != nil {
			sb.WriteString(" = ")
			sb.WriteString(*a.Value)
		}
		if a.Comments.Right != nil {
			sb.WriteString(" ")
			sb.WriteString(a.Comments.Right.Text)
		}
	}
	sb.WriteString("\n}")
	return sb.String()
}
