// Code generated from RestPath.g4 by ANTLR 4.13.2. DO NOT EDIT.

package pidl // RestPath
import "github.com/antlr4-go/antlr/v4"

// RestPathListener is a complete listener for a parse tree produced by RestPathParser.
type RestPathListener interface {
	antlr.ParseTreeListener

	// EnterPath is called when entering the path production.
	EnterPath(c *PathContext)

	// EnterSegment is called when entering the segment production.
	EnterSegment(c *SegmentContext)

	// ExitPath is called when exiting the path production.
	ExitPath(c *PathContext)

	// ExitSegment is called when exiting the segment production.
	ExitSegment(c *SegmentContext)
}
