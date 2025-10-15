// Code generated from RestPath.g4 by ANTLR 4.13.2. DO NOT EDIT.

package pidl

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type RestPathLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var RestPathLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func restpathlexerLexerInit() {
	staticData := &RestPathLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "", "", "':'", "'/'", "'{'", "'}'", "'*'", "'...'",
	}
	staticData.SymbolicNames = []string{
		"", "STATIC_SEGMENT", "IDENTIFIER", "COLON", "SLASH", "LBRACE", "RBRACE",
		"STAR", "ELLIPSIS",
	}
	staticData.RuleNames = []string{
		"STATIC_SEGMENT", "IDENTIFIER", "COLON", "SLASH", "LBRACE", "RBRACE",
		"STAR", "ELLIPSIS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 8, 43, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 1, 0, 4, 0, 19, 8, 0, 11,
		0, 12, 0, 20, 1, 1, 1, 1, 5, 1, 25, 8, 1, 10, 1, 12, 1, 28, 9, 1, 1, 2,
		1, 2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7,
		1, 7, 0, 0, 8, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 1, 0,
		3, 5, 0, 45, 45, 48, 57, 65, 90, 95, 95, 97, 122, 3, 0, 65, 90, 95, 95,
		97, 122, 4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 44, 0, 1, 1, 0, 0, 0, 0,
		3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0,
		11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 1, 18, 1, 0, 0, 0,
		3, 22, 1, 0, 0, 0, 5, 29, 1, 0, 0, 0, 7, 31, 1, 0, 0, 0, 9, 33, 1, 0, 0,
		0, 11, 35, 1, 0, 0, 0, 13, 37, 1, 0, 0, 0, 15, 39, 1, 0, 0, 0, 17, 19,
		7, 0, 0, 0, 18, 17, 1, 0, 0, 0, 19, 20, 1, 0, 0, 0, 20, 18, 1, 0, 0, 0,
		20, 21, 1, 0, 0, 0, 21, 2, 1, 0, 0, 0, 22, 26, 7, 1, 0, 0, 23, 25, 7, 2,
		0, 0, 24, 23, 1, 0, 0, 0, 25, 28, 1, 0, 0, 0, 26, 24, 1, 0, 0, 0, 26, 27,
		1, 0, 0, 0, 27, 4, 1, 0, 0, 0, 28, 26, 1, 0, 0, 0, 29, 30, 5, 58, 0, 0,
		30, 6, 1, 0, 0, 0, 31, 32, 5, 47, 0, 0, 32, 8, 1, 0, 0, 0, 33, 34, 5, 123,
		0, 0, 34, 10, 1, 0, 0, 0, 35, 36, 5, 125, 0, 0, 36, 12, 1, 0, 0, 0, 37,
		38, 5, 42, 0, 0, 38, 14, 1, 0, 0, 0, 39, 40, 5, 46, 0, 0, 40, 41, 5, 46,
		0, 0, 41, 42, 5, 46, 0, 0, 42, 16, 1, 0, 0, 0, 3, 0, 20, 26, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// RestPathLexerInit initializes any static state used to implement RestPathLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewRestPathLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func RestPathLexerInit() {
	staticData := &RestPathLexerLexerStaticData
	staticData.once.Do(restpathlexerLexerInit)
}

// NewRestPathLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewRestPathLexer(input antlr.CharStream) *RestPathLexer {
	RestPathLexerInit()
	l := new(RestPathLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &RestPathLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "RestPath.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// RestPathLexer tokens.
const (
	RestPathLexerSTATIC_SEGMENT = 1
	RestPathLexerIDENTIFIER     = 2
	RestPathLexerCOLON          = 3
	RestPathLexerSLASH          = 4
	RestPathLexerLBRACE         = 5
	RestPathLexerRBRACE         = 6
	RestPathLexerSTAR           = 7
	RestPathLexerELLIPSIS       = 8
)
