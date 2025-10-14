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
		"", "'/'", "':'", "'*'", "'{'", "'...'", "'}'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "", "STATIC_SEGMENT", "IDENTIFIER",
	}
	staticData.RuleNames = []string{
		"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "STATIC_SEGMENT", "IDENTIFIER",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 8, 43, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 1, 0, 1, 0, 1, 1, 1, 1, 1,
		2, 1, 2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 6, 4, 6, 33,
		8, 6, 11, 6, 12, 6, 34, 1, 7, 1, 7, 5, 7, 39, 8, 7, 10, 7, 12, 7, 42, 9,
		7, 0, 0, 8, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 1, 0, 3,
		5, 0, 45, 45, 48, 57, 65, 90, 95, 95, 97, 122, 3, 0, 65, 90, 95, 95, 97,
		122, 4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 44, 0, 1, 1, 0, 0, 0, 0, 3,
		1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11,
		1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 1, 17, 1, 0, 0, 0, 3,
		19, 1, 0, 0, 0, 5, 21, 1, 0, 0, 0, 7, 23, 1, 0, 0, 0, 9, 25, 1, 0, 0, 0,
		11, 29, 1, 0, 0, 0, 13, 32, 1, 0, 0, 0, 15, 36, 1, 0, 0, 0, 17, 18, 5,
		47, 0, 0, 18, 2, 1, 0, 0, 0, 19, 20, 5, 58, 0, 0, 20, 4, 1, 0, 0, 0, 21,
		22, 5, 42, 0, 0, 22, 6, 1, 0, 0, 0, 23, 24, 5, 123, 0, 0, 24, 8, 1, 0,
		0, 0, 25, 26, 5, 46, 0, 0, 26, 27, 5, 46, 0, 0, 27, 28, 5, 46, 0, 0, 28,
		10, 1, 0, 0, 0, 29, 30, 5, 125, 0, 0, 30, 12, 1, 0, 0, 0, 31, 33, 7, 0,
		0, 0, 32, 31, 1, 0, 0, 0, 33, 34, 1, 0, 0, 0, 34, 32, 1, 0, 0, 0, 34, 35,
		1, 0, 0, 0, 35, 14, 1, 0, 0, 0, 36, 40, 7, 1, 0, 0, 37, 39, 7, 2, 0, 0,
		38, 37, 1, 0, 0, 0, 39, 42, 1, 0, 0, 0, 40, 38, 1, 0, 0, 0, 40, 41, 1,
		0, 0, 0, 41, 16, 1, 0, 0, 0, 42, 40, 1, 0, 0, 0, 3, 0, 34, 40, 0,
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
	RestPathLexerT__0           = 1
	RestPathLexerT__1           = 2
	RestPathLexerT__2           = 3
	RestPathLexerT__3           = 4
	RestPathLexerT__4           = 5
	RestPathLexerT__5           = 6
	RestPathLexerSTATIC_SEGMENT = 7
	RestPathLexerIDENTIFIER     = 8
)
