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
		"", "'/'",
	}
	staticData.SymbolicNames = []string{
		"", "", "STATIC_SEGMENT", "PARAM_SEGMENT", "BRACED_PARAM", "IDENTIFIER",
	}
	staticData.RuleNames = []string{
		"T__0", "STATIC_SEGMENT", "PARAM_SEGMENT", "BRACED_PARAM", "IDENTIFIER",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 5, 39, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 1, 0, 1, 0, 1, 1, 4, 1, 15, 8, 1, 11, 1, 12, 1, 16, 1, 2, 1, 2,
		1, 2, 3, 2, 22, 8, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 3, 3, 29, 8, 3, 1,
		3, 1, 3, 1, 4, 1, 4, 5, 4, 35, 8, 4, 10, 4, 12, 4, 38, 9, 4, 0, 0, 5, 1,
		1, 3, 2, 5, 3, 7, 4, 9, 5, 1, 0, 3, 5, 0, 45, 45, 48, 57, 65, 90, 95, 95,
		97, 122, 3, 0, 65, 90, 95, 95, 97, 122, 4, 0, 48, 57, 65, 90, 95, 95, 97,
		122, 42, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1,
		0, 0, 0, 0, 9, 1, 0, 0, 0, 1, 11, 1, 0, 0, 0, 3, 14, 1, 0, 0, 0, 5, 18,
		1, 0, 0, 0, 7, 23, 1, 0, 0, 0, 9, 32, 1, 0, 0, 0, 11, 12, 5, 47, 0, 0,
		12, 2, 1, 0, 0, 0, 13, 15, 7, 0, 0, 0, 14, 13, 1, 0, 0, 0, 15, 16, 1, 0,
		0, 0, 16, 14, 1, 0, 0, 0, 16, 17, 1, 0, 0, 0, 17, 4, 1, 0, 0, 0, 18, 19,
		5, 58, 0, 0, 19, 21, 3, 9, 4, 0, 20, 22, 5, 42, 0, 0, 21, 20, 1, 0, 0,
		0, 21, 22, 1, 0, 0, 0, 22, 6, 1, 0, 0, 0, 23, 24, 5, 123, 0, 0, 24, 28,
		3, 9, 4, 0, 25, 26, 5, 46, 0, 0, 26, 27, 5, 46, 0, 0, 27, 29, 5, 46, 0,
		0, 28, 25, 1, 0, 0, 0, 28, 29, 1, 0, 0, 0, 29, 30, 1, 0, 0, 0, 30, 31,
		5, 125, 0, 0, 31, 8, 1, 0, 0, 0, 32, 36, 7, 1, 0, 0, 33, 35, 7, 2, 0, 0,
		34, 33, 1, 0, 0, 0, 35, 38, 1, 0, 0, 0, 36, 34, 1, 0, 0, 0, 36, 37, 1,
		0, 0, 0, 37, 10, 1, 0, 0, 0, 38, 36, 1, 0, 0, 0, 5, 0, 16, 21, 28, 36,
		0,
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
	RestPathLexerSTATIC_SEGMENT = 2
	RestPathLexerPARAM_SEGMENT  = 3
	RestPathLexerBRACED_PARAM   = 4
	RestPathLexerIDENTIFIER     = 5
)
