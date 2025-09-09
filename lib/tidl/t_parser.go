// Code generated from TParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package tidl // TParser
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type TParser struct {
	*antlr.BaseParser
}

var TParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func tparserParserInit() {
	staticData := &TParserParserStaticData
	staticData.LiteralNames = []string{
		"", "'const'", "'enum'", "'type'", "'rpc'", "'true'", "'false'", "'any'",
		"'bool'", "'int'", "'float'", "'string'", "'binary'", "'stream'", "'map'",
		"'list'", "'<'", "'>'", "'('", "')'", "'{'", "'}'", "'='", "','", "'?'",
	}
	staticData.SymbolicNames = []string{
		"", "KW_CONST", "KW_ENUM", "KW_TYPE", "KW_RPC", "KW_TRUE", "KW_FALSE",
		"TYPE_ANY", "TYPE_BOOL", "TYPE_INT", "TYPE_FLOAT", "TYPE_STRING", "TYPE_BINARY",
		"TYPE_STREAM", "TYPE_MAP", "TYPE_LIST", "LESS_THAN", "GREATER_THAN",
		"LEFT_PAREN", "RIGHT_PAREN", "LEFT_BRACE", "RIGHT_BRACE", "EQUAL", "COMMA",
		"QUESTION", "STRING", "IDENTIFIER", "INTEGER", "FLOAT", "WHITESPACE",
		"SINGLE_LINE_COMMENT", "MULTI_LINE_COMMENT",
	}
	staticData.RuleNames = []string{
		"document", "definition", "const_def", "const_type", "enum_def", "enum_field",
		"type_def", "type_field", "embed_type_field", "common_type_field", "common_field_type",
		"generic_type", "type_annotations", "rpc_def", "rpc_req", "rpc_resp",
		"rpc_annotations", "annotation", "base_type", "user_type", "container_type",
		"map_type", "key_type", "list_type", "value_type", "const_value",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 31, 217, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 1, 0,
		5, 0, 54, 8, 0, 10, 0, 12, 0, 57, 9, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1,
		1, 3, 1, 65, 8, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 4,
		1, 4, 1, 4, 1, 4, 5, 4, 79, 8, 4, 10, 4, 12, 4, 82, 9, 4, 1, 4, 1, 4, 1,
		5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 3, 6, 95, 8, 6, 1, 6,
		1, 6, 5, 6, 99, 8, 6, 10, 6, 12, 6, 102, 9, 6, 1, 6, 1, 6, 1, 6, 1, 6,
		1, 6, 1, 6, 1, 6, 1, 6, 3, 6, 112, 8, 6, 1, 7, 1, 7, 3, 7, 116, 8, 7, 1,
		8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9, 124, 8, 9, 1, 9, 3, 9, 127, 8, 9,
		1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 3, 10, 134, 8, 10, 1, 11, 1, 11, 1,
		11, 3, 11, 139, 8, 11, 1, 12, 1, 12, 1, 12, 1, 12, 5, 12, 145, 8, 12, 10,
		12, 12, 12, 148, 9, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13,
		1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1,
		15, 3, 15, 168, 8, 15, 1, 16, 1, 16, 5, 16, 172, 8, 16, 10, 16, 12, 16,
		175, 9, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17, 3, 17, 182, 8, 17, 1, 18,
		1, 18, 3, 18, 186, 8, 18, 1, 19, 1, 19, 3, 19, 190, 8, 19, 1, 20, 1, 20,
		3, 20, 194, 8, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1,
		22, 1, 22, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 24, 1, 24, 1, 24, 3, 24,
		213, 8, 24, 1, 25, 1, 25, 1, 25, 0, 0, 26, 0, 2, 4, 6, 8, 10, 12, 14, 16,
		18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 0,
		3, 1, 0, 8, 11, 2, 0, 9, 9, 11, 11, 2, 0, 5, 6, 25, 28, 216, 0, 55, 1,
		0, 0, 0, 2, 64, 1, 0, 0, 0, 4, 66, 1, 0, 0, 0, 6, 72, 1, 0, 0, 0, 8, 74,
		1, 0, 0, 0, 10, 85, 1, 0, 0, 0, 12, 111, 1, 0, 0, 0, 14, 115, 1, 0, 0,
		0, 16, 117, 1, 0, 0, 0, 18, 119, 1, 0, 0, 0, 20, 133, 1, 0, 0, 0, 22, 138,
		1, 0, 0, 0, 24, 140, 1, 0, 0, 0, 26, 151, 1, 0, 0, 0, 28, 159, 1, 0, 0,
		0, 30, 167, 1, 0, 0, 0, 32, 169, 1, 0, 0, 0, 34, 178, 1, 0, 0, 0, 36, 183,
		1, 0, 0, 0, 38, 187, 1, 0, 0, 0, 40, 193, 1, 0, 0, 0, 42, 195, 1, 0, 0,
		0, 44, 202, 1, 0, 0, 0, 46, 204, 1, 0, 0, 0, 48, 212, 1, 0, 0, 0, 50, 214,
		1, 0, 0, 0, 52, 54, 3, 2, 1, 0, 53, 52, 1, 0, 0, 0, 54, 57, 1, 0, 0, 0,
		55, 53, 1, 0, 0, 0, 55, 56, 1, 0, 0, 0, 56, 58, 1, 0, 0, 0, 57, 55, 1,
		0, 0, 0, 58, 59, 5, 0, 0, 1, 59, 1, 1, 0, 0, 0, 60, 65, 3, 4, 2, 0, 61,
		65, 3, 8, 4, 0, 62, 65, 3, 12, 6, 0, 63, 65, 3, 26, 13, 0, 64, 60, 1, 0,
		0, 0, 64, 61, 1, 0, 0, 0, 64, 62, 1, 0, 0, 0, 64, 63, 1, 0, 0, 0, 65, 3,
		1, 0, 0, 0, 66, 67, 5, 1, 0, 0, 67, 68, 3, 6, 3, 0, 68, 69, 5, 26, 0, 0,
		69, 70, 5, 22, 0, 0, 70, 71, 3, 50, 25, 0, 71, 5, 1, 0, 0, 0, 72, 73, 7,
		0, 0, 0, 73, 7, 1, 0, 0, 0, 74, 75, 5, 2, 0, 0, 75, 76, 5, 26, 0, 0, 76,
		80, 5, 20, 0, 0, 77, 79, 3, 10, 5, 0, 78, 77, 1, 0, 0, 0, 79, 82, 1, 0,
		0, 0, 80, 78, 1, 0, 0, 0, 80, 81, 1, 0, 0, 0, 81, 83, 1, 0, 0, 0, 82, 80,
		1, 0, 0, 0, 83, 84, 5, 21, 0, 0, 84, 9, 1, 0, 0, 0, 85, 86, 5, 26, 0, 0,
		86, 87, 5, 22, 0, 0, 87, 88, 5, 27, 0, 0, 88, 11, 1, 0, 0, 0, 89, 90, 5,
		3, 0, 0, 90, 94, 5, 26, 0, 0, 91, 92, 5, 16, 0, 0, 92, 93, 5, 26, 0, 0,
		93, 95, 5, 17, 0, 0, 94, 91, 1, 0, 0, 0, 94, 95, 1, 0, 0, 0, 95, 96, 1,
		0, 0, 0, 96, 100, 5, 20, 0, 0, 97, 99, 3, 14, 7, 0, 98, 97, 1, 0, 0, 0,
		99, 102, 1, 0, 0, 0, 100, 98, 1, 0, 0, 0, 100, 101, 1, 0, 0, 0, 101, 103,
		1, 0, 0, 0, 102, 100, 1, 0, 0, 0, 103, 112, 5, 21, 0, 0, 104, 105, 5, 3,
		0, 0, 105, 106, 5, 26, 0, 0, 106, 107, 5, 26, 0, 0, 107, 108, 5, 16, 0,
		0, 108, 109, 3, 22, 11, 0, 109, 110, 5, 17, 0, 0, 110, 112, 1, 0, 0, 0,
		111, 89, 1, 0, 0, 0, 111, 104, 1, 0, 0, 0, 112, 13, 1, 0, 0, 0, 113, 116,
		3, 18, 9, 0, 114, 116, 3, 16, 8, 0, 115, 113, 1, 0, 0, 0, 115, 114, 1,
		0, 0, 0, 116, 15, 1, 0, 0, 0, 117, 118, 3, 38, 19, 0, 118, 17, 1, 0, 0,
		0, 119, 120, 3, 20, 10, 0, 120, 123, 5, 26, 0, 0, 121, 122, 5, 22, 0, 0,
		122, 124, 3, 50, 25, 0, 123, 121, 1, 0, 0, 0, 123, 124, 1, 0, 0, 0, 124,
		126, 1, 0, 0, 0, 125, 127, 3, 24, 12, 0, 126, 125, 1, 0, 0, 0, 126, 127,
		1, 0, 0, 0, 127, 19, 1, 0, 0, 0, 128, 134, 5, 7, 0, 0, 129, 134, 3, 36,
		18, 0, 130, 134, 3, 38, 19, 0, 131, 134, 3, 40, 20, 0, 132, 134, 5, 12,
		0, 0, 133, 128, 1, 0, 0, 0, 133, 129, 1, 0, 0, 0, 133, 130, 1, 0, 0, 0,
		133, 131, 1, 0, 0, 0, 133, 132, 1, 0, 0, 0, 134, 21, 1, 0, 0, 0, 135, 139,
		3, 36, 18, 0, 136, 139, 3, 38, 19, 0, 137, 139, 3, 40, 20, 0, 138, 135,
		1, 0, 0, 0, 138, 136, 1, 0, 0, 0, 138, 137, 1, 0, 0, 0, 139, 23, 1, 0,
		0, 0, 140, 141, 5, 18, 0, 0, 141, 146, 3, 34, 17, 0, 142, 143, 5, 23, 0,
		0, 143, 145, 3, 34, 17, 0, 144, 142, 1, 0, 0, 0, 145, 148, 1, 0, 0, 0,
		146, 144, 1, 0, 0, 0, 146, 147, 1, 0, 0, 0, 147, 149, 1, 0, 0, 0, 148,
		146, 1, 0, 0, 0, 149, 150, 5, 19, 0, 0, 150, 25, 1, 0, 0, 0, 151, 152,
		5, 4, 0, 0, 152, 153, 5, 26, 0, 0, 153, 154, 5, 18, 0, 0, 154, 155, 3,
		28, 14, 0, 155, 156, 5, 19, 0, 0, 156, 157, 3, 30, 15, 0, 157, 158, 3,
		32, 16, 0, 158, 27, 1, 0, 0, 0, 159, 160, 5, 26, 0, 0, 160, 29, 1, 0, 0,
		0, 161, 168, 5, 26, 0, 0, 162, 163, 5, 13, 0, 0, 163, 164, 5, 16, 0, 0,
		164, 165, 3, 38, 19, 0, 165, 166, 5, 17, 0, 0, 166, 168, 1, 0, 0, 0, 167,
		161, 1, 0, 0, 0, 167, 162, 1, 0, 0, 0, 168, 31, 1, 0, 0, 0, 169, 173, 5,
		20, 0, 0, 170, 172, 3, 34, 17, 0, 171, 170, 1, 0, 0, 0, 172, 175, 1, 0,
		0, 0, 173, 171, 1, 0, 0, 0, 173, 174, 1, 0, 0, 0, 174, 176, 1, 0, 0, 0,
		175, 173, 1, 0, 0, 0, 176, 177, 5, 21, 0, 0, 177, 33, 1, 0, 0, 0, 178,
		181, 5, 26, 0, 0, 179, 180, 5, 22, 0, 0, 180, 182, 3, 50, 25, 0, 181, 179,
		1, 0, 0, 0, 181, 182, 1, 0, 0, 0, 182, 35, 1, 0, 0, 0, 183, 185, 7, 0,
		0, 0, 184, 186, 5, 24, 0, 0, 185, 184, 1, 0, 0, 0, 185, 186, 1, 0, 0, 0,
		186, 37, 1, 0, 0, 0, 187, 189, 5, 26, 0, 0, 188, 190, 5, 24, 0, 0, 189,
		188, 1, 0, 0, 0, 189, 190, 1, 0, 0, 0, 190, 39, 1, 0, 0, 0, 191, 194, 3,
		42, 21, 0, 192, 194, 3, 46, 23, 0, 193, 191, 1, 0, 0, 0, 193, 192, 1, 0,
		0, 0, 194, 41, 1, 0, 0, 0, 195, 196, 5, 14, 0, 0, 196, 197, 5, 16, 0, 0,
		197, 198, 3, 44, 22, 0, 198, 199, 5, 23, 0, 0, 199, 200, 3, 48, 24, 0,
		200, 201, 5, 17, 0, 0, 201, 43, 1, 0, 0, 0, 202, 203, 7, 1, 0, 0, 203,
		45, 1, 0, 0, 0, 204, 205, 5, 15, 0, 0, 205, 206, 5, 16, 0, 0, 206, 207,
		3, 48, 24, 0, 207, 208, 5, 17, 0, 0, 208, 47, 1, 0, 0, 0, 209, 213, 3,
		36, 18, 0, 210, 213, 3, 38, 19, 0, 211, 213, 3, 40, 20, 0, 212, 209, 1,
		0, 0, 0, 212, 210, 1, 0, 0, 0, 212, 211, 1, 0, 0, 0, 213, 49, 1, 0, 0,
		0, 214, 215, 7, 2, 0, 0, 215, 51, 1, 0, 0, 0, 19, 55, 64, 80, 94, 100,
		111, 115, 123, 126, 133, 138, 146, 167, 173, 181, 185, 189, 193, 212,
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

// TParserInit initializes any static state used to implement TParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewTParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func TParserInit() {
	staticData := &TParserParserStaticData
	staticData.once.Do(tparserParserInit)
}

// NewTParser produces a new parser instance for the optional input antlr.TokenStream.
func NewTParser(input antlr.TokenStream) *TParser {
	TParserInit()
	this := new(TParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &TParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "TParser.g4"

	return this
}

// TParser tokens.
const (
	TParserEOF                 = antlr.TokenEOF
	TParserKW_CONST            = 1
	TParserKW_ENUM             = 2
	TParserKW_TYPE             = 3
	TParserKW_RPC              = 4
	TParserKW_TRUE             = 5
	TParserKW_FALSE            = 6
	TParserTYPE_ANY            = 7
	TParserTYPE_BOOL           = 8
	TParserTYPE_INT            = 9
	TParserTYPE_FLOAT          = 10
	TParserTYPE_STRING         = 11
	TParserTYPE_BINARY         = 12
	TParserTYPE_STREAM         = 13
	TParserTYPE_MAP            = 14
	TParserTYPE_LIST           = 15
	TParserLESS_THAN           = 16
	TParserGREATER_THAN        = 17
	TParserLEFT_PAREN          = 18
	TParserRIGHT_PAREN         = 19
	TParserLEFT_BRACE          = 20
	TParserRIGHT_BRACE         = 21
	TParserEQUAL               = 22
	TParserCOMMA               = 23
	TParserQUESTION            = 24
	TParserSTRING              = 25
	TParserIDENTIFIER          = 26
	TParserINTEGER             = 27
	TParserFLOAT               = 28
	TParserWHITESPACE          = 29
	TParserSINGLE_LINE_COMMENT = 30
	TParserMULTI_LINE_COMMENT  = 31
)

// TParser rules.
const (
	TParserRULE_document          = 0
	TParserRULE_definition        = 1
	TParserRULE_const_def         = 2
	TParserRULE_const_type        = 3
	TParserRULE_enum_def          = 4
	TParserRULE_enum_field        = 5
	TParserRULE_type_def          = 6
	TParserRULE_type_field        = 7
	TParserRULE_embed_type_field  = 8
	TParserRULE_common_type_field = 9
	TParserRULE_common_field_type = 10
	TParserRULE_generic_type      = 11
	TParserRULE_type_annotations  = 12
	TParserRULE_rpc_def           = 13
	TParserRULE_rpc_req           = 14
	TParserRULE_rpc_resp          = 15
	TParserRULE_rpc_annotations   = 16
	TParserRULE_annotation        = 17
	TParserRULE_base_type         = 18
	TParserRULE_user_type         = 19
	TParserRULE_container_type    = 20
	TParserRULE_map_type          = 21
	TParserRULE_key_type          = 22
	TParserRULE_list_type         = 23
	TParserRULE_value_type        = 24
	TParserRULE_const_value       = 25
)

// IDocumentContext is an interface to support dynamic dispatch.
type IDocumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOF() antlr.TerminalNode
	AllDefinition() []IDefinitionContext
	Definition(i int) IDefinitionContext

	// IsDocumentContext differentiates from other interfaces.
	IsDocumentContext()
}

type DocumentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDocumentContext() *DocumentContext {
	var p = new(DocumentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_document
	return p
}

func InitEmptyDocumentContext(p *DocumentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_document
}

func (*DocumentContext) IsDocumentContext() {}

func NewDocumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DocumentContext {
	var p = new(DocumentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_document

	return p
}

func (s *DocumentContext) GetParser() antlr.Parser { return s.parser }

func (s *DocumentContext) EOF() antlr.TerminalNode {
	return s.GetToken(TParserEOF, 0)
}

func (s *DocumentContext) AllDefinition() []IDefinitionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDefinitionContext); ok {
			len++
		}
	}

	tst := make([]IDefinitionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDefinitionContext); ok {
			tst[i] = t.(IDefinitionContext)
			i++
		}
	}

	return tst
}

func (s *DocumentContext) Definition(i int) IDefinitionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDefinitionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDefinitionContext)
}

func (s *DocumentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DocumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DocumentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterDocument(s)
	}
}

func (s *DocumentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitDocument(s)
	}
}

func (p *TParser) Document() (localctx IDocumentContext) {
	localctx = NewDocumentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, TParserRULE_document)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(55)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&30) != 0 {
		{
			p.SetState(52)
			p.Definition()
		}

		p.SetState(57)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(58)
		p.Match(TParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDefinitionContext is an interface to support dynamic dispatch.
type IDefinitionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Const_def() IConst_defContext
	Enum_def() IEnum_defContext
	Type_def() IType_defContext
	Rpc_def() IRpc_defContext

	// IsDefinitionContext differentiates from other interfaces.
	IsDefinitionContext()
}

type DefinitionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDefinitionContext() *DefinitionContext {
	var p = new(DefinitionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_definition
	return p
}

func InitEmptyDefinitionContext(p *DefinitionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_definition
}

func (*DefinitionContext) IsDefinitionContext() {}

func NewDefinitionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DefinitionContext {
	var p = new(DefinitionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_definition

	return p
}

func (s *DefinitionContext) GetParser() antlr.Parser { return s.parser }

func (s *DefinitionContext) Const_def() IConst_defContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConst_defContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConst_defContext)
}

func (s *DefinitionContext) Enum_def() IEnum_defContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnum_defContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnum_defContext)
}

func (s *DefinitionContext) Type_def() IType_defContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IType_defContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IType_defContext)
}

func (s *DefinitionContext) Rpc_def() IRpc_defContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRpc_defContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRpc_defContext)
}

func (s *DefinitionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DefinitionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DefinitionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterDefinition(s)
	}
}

func (s *DefinitionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitDefinition(s)
	}
}

func (p *TParser) Definition() (localctx IDefinitionContext) {
	localctx = NewDefinitionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, TParserRULE_definition)
	p.SetState(64)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TParserKW_CONST:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(60)
			p.Const_def()
		}

	case TParserKW_ENUM:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(61)
			p.Enum_def()
		}

	case TParserKW_TYPE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(62)
			p.Type_def()
		}

	case TParserKW_RPC:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(63)
			p.Rpc_def()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConst_defContext is an interface to support dynamic dispatch.
type IConst_defContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	KW_CONST() antlr.TerminalNode
	Const_type() IConst_typeContext
	IDENTIFIER() antlr.TerminalNode
	EQUAL() antlr.TerminalNode
	Const_value() IConst_valueContext

	// IsConst_defContext differentiates from other interfaces.
	IsConst_defContext()
}

type Const_defContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConst_defContext() *Const_defContext {
	var p = new(Const_defContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_const_def
	return p
}

func InitEmptyConst_defContext(p *Const_defContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_const_def
}

func (*Const_defContext) IsConst_defContext() {}

func NewConst_defContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Const_defContext {
	var p = new(Const_defContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_const_def

	return p
}

func (s *Const_defContext) GetParser() antlr.Parser { return s.parser }

func (s *Const_defContext) KW_CONST() antlr.TerminalNode {
	return s.GetToken(TParserKW_CONST, 0)
}

func (s *Const_defContext) Const_type() IConst_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConst_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConst_typeContext)
}

func (s *Const_defContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(TParserIDENTIFIER, 0)
}

func (s *Const_defContext) EQUAL() antlr.TerminalNode {
	return s.GetToken(TParserEQUAL, 0)
}

func (s *Const_defContext) Const_value() IConst_valueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConst_valueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConst_valueContext)
}

func (s *Const_defContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Const_defContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Const_defContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterConst_def(s)
	}
}

func (s *Const_defContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitConst_def(s)
	}
}

func (p *TParser) Const_def() (localctx IConst_defContext) {
	localctx = NewConst_defContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, TParserRULE_const_def)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(66)
		p.Match(TParserKW_CONST)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(67)
		p.Const_type()
	}
	{
		p.SetState(68)
		p.Match(TParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(69)
		p.Match(TParserEQUAL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(70)
		p.Const_value()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConst_typeContext is an interface to support dynamic dispatch.
type IConst_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPE_BOOL() antlr.TerminalNode
	TYPE_INT() antlr.TerminalNode
	TYPE_FLOAT() antlr.TerminalNode
	TYPE_STRING() antlr.TerminalNode

	// IsConst_typeContext differentiates from other interfaces.
	IsConst_typeContext()
}

type Const_typeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConst_typeContext() *Const_typeContext {
	var p = new(Const_typeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_const_type
	return p
}

func InitEmptyConst_typeContext(p *Const_typeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_const_type
}

func (*Const_typeContext) IsConst_typeContext() {}

func NewConst_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Const_typeContext {
	var p = new(Const_typeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_const_type

	return p
}

func (s *Const_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Const_typeContext) TYPE_BOOL() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_BOOL, 0)
}

func (s *Const_typeContext) TYPE_INT() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_INT, 0)
}

func (s *Const_typeContext) TYPE_FLOAT() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_FLOAT, 0)
}

func (s *Const_typeContext) TYPE_STRING() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_STRING, 0)
}

func (s *Const_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Const_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Const_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterConst_type(s)
	}
}

func (s *Const_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitConst_type(s)
	}
}

func (p *TParser) Const_type() (localctx IConst_typeContext) {
	localctx = NewConst_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, TParserRULE_const_type)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(72)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&3840) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEnum_defContext is an interface to support dynamic dispatch.
type IEnum_defContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	KW_ENUM() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LEFT_BRACE() antlr.TerminalNode
	RIGHT_BRACE() antlr.TerminalNode
	AllEnum_field() []IEnum_fieldContext
	Enum_field(i int) IEnum_fieldContext

	// IsEnum_defContext differentiates from other interfaces.
	IsEnum_defContext()
}

type Enum_defContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEnum_defContext() *Enum_defContext {
	var p = new(Enum_defContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_enum_def
	return p
}

func InitEmptyEnum_defContext(p *Enum_defContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_enum_def
}

func (*Enum_defContext) IsEnum_defContext() {}

func NewEnum_defContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Enum_defContext {
	var p = new(Enum_defContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_enum_def

	return p
}

func (s *Enum_defContext) GetParser() antlr.Parser { return s.parser }

func (s *Enum_defContext) KW_ENUM() antlr.TerminalNode {
	return s.GetToken(TParserKW_ENUM, 0)
}

func (s *Enum_defContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(TParserIDENTIFIER, 0)
}

func (s *Enum_defContext) LEFT_BRACE() antlr.TerminalNode {
	return s.GetToken(TParserLEFT_BRACE, 0)
}

func (s *Enum_defContext) RIGHT_BRACE() antlr.TerminalNode {
	return s.GetToken(TParserRIGHT_BRACE, 0)
}

func (s *Enum_defContext) AllEnum_field() []IEnum_fieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IEnum_fieldContext); ok {
			len++
		}
	}

	tst := make([]IEnum_fieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IEnum_fieldContext); ok {
			tst[i] = t.(IEnum_fieldContext)
			i++
		}
	}

	return tst
}

func (s *Enum_defContext) Enum_field(i int) IEnum_fieldContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnum_fieldContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnum_fieldContext)
}

func (s *Enum_defContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Enum_defContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Enum_defContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterEnum_def(s)
	}
}

func (s *Enum_defContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitEnum_def(s)
	}
}

func (p *TParser) Enum_def() (localctx IEnum_defContext) {
	localctx = NewEnum_defContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, TParserRULE_enum_def)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(74)
		p.Match(TParserKW_ENUM)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(75)
		p.Match(TParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(76)
		p.Match(TParserLEFT_BRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(80)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == TParserIDENTIFIER {
		{
			p.SetState(77)
			p.Enum_field()
		}

		p.SetState(82)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(83)
		p.Match(TParserRIGHT_BRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEnum_fieldContext is an interface to support dynamic dispatch.
type IEnum_fieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	EQUAL() antlr.TerminalNode
	INTEGER() antlr.TerminalNode

	// IsEnum_fieldContext differentiates from other interfaces.
	IsEnum_fieldContext()
}

type Enum_fieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEnum_fieldContext() *Enum_fieldContext {
	var p = new(Enum_fieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_enum_field
	return p
}

func InitEmptyEnum_fieldContext(p *Enum_fieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_enum_field
}

func (*Enum_fieldContext) IsEnum_fieldContext() {}

func NewEnum_fieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Enum_fieldContext {
	var p = new(Enum_fieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_enum_field

	return p
}

func (s *Enum_fieldContext) GetParser() antlr.Parser { return s.parser }

func (s *Enum_fieldContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(TParserIDENTIFIER, 0)
}

func (s *Enum_fieldContext) EQUAL() antlr.TerminalNode {
	return s.GetToken(TParserEQUAL, 0)
}

func (s *Enum_fieldContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(TParserINTEGER, 0)
}

func (s *Enum_fieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Enum_fieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Enum_fieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterEnum_field(s)
	}
}

func (s *Enum_fieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitEnum_field(s)
	}
}

func (p *TParser) Enum_field() (localctx IEnum_fieldContext) {
	localctx = NewEnum_fieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, TParserRULE_enum_field)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(85)
		p.Match(TParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(86)
		p.Match(TParserEQUAL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(87)
		p.Match(TParserINTEGER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IType_defContext is an interface to support dynamic dispatch.
type IType_defContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	KW_TYPE() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	LEFT_BRACE() antlr.TerminalNode
	RIGHT_BRACE() antlr.TerminalNode
	LESS_THAN() antlr.TerminalNode
	GREATER_THAN() antlr.TerminalNode
	AllType_field() []IType_fieldContext
	Type_field(i int) IType_fieldContext
	Generic_type() IGeneric_typeContext

	// IsType_defContext differentiates from other interfaces.
	IsType_defContext()
}

type Type_defContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyType_defContext() *Type_defContext {
	var p = new(Type_defContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_type_def
	return p
}

func InitEmptyType_defContext(p *Type_defContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_type_def
}

func (*Type_defContext) IsType_defContext() {}

func NewType_defContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Type_defContext {
	var p = new(Type_defContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_type_def

	return p
}

func (s *Type_defContext) GetParser() antlr.Parser { return s.parser }

func (s *Type_defContext) KW_TYPE() antlr.TerminalNode {
	return s.GetToken(TParserKW_TYPE, 0)
}

func (s *Type_defContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(TParserIDENTIFIER)
}

func (s *Type_defContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(TParserIDENTIFIER, i)
}

func (s *Type_defContext) LEFT_BRACE() antlr.TerminalNode {
	return s.GetToken(TParserLEFT_BRACE, 0)
}

func (s *Type_defContext) RIGHT_BRACE() antlr.TerminalNode {
	return s.GetToken(TParserRIGHT_BRACE, 0)
}

func (s *Type_defContext) LESS_THAN() antlr.TerminalNode {
	return s.GetToken(TParserLESS_THAN, 0)
}

func (s *Type_defContext) GREATER_THAN() antlr.TerminalNode {
	return s.GetToken(TParserGREATER_THAN, 0)
}

func (s *Type_defContext) AllType_field() []IType_fieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IType_fieldContext); ok {
			len++
		}
	}

	tst := make([]IType_fieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IType_fieldContext); ok {
			tst[i] = t.(IType_fieldContext)
			i++
		}
	}

	return tst
}

func (s *Type_defContext) Type_field(i int) IType_fieldContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IType_fieldContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IType_fieldContext)
}

func (s *Type_defContext) Generic_type() IGeneric_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGeneric_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGeneric_typeContext)
}

func (s *Type_defContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Type_defContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Type_defContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterType_def(s)
	}
}

func (s *Type_defContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitType_def(s)
	}
}

func (p *TParser) Type_def() (localctx IType_defContext) {
	localctx = NewType_defContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, TParserRULE_type_def)
	var _la int

	p.SetState(111)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(89)
			p.Match(TParserKW_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(90)
			p.Match(TParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(94)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == TParserLESS_THAN {
			{
				p.SetState(91)
				p.Match(TParserLESS_THAN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(92)
				p.Match(TParserIDENTIFIER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(93)
				p.Match(TParserGREATER_THAN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(96)
			p.Match(TParserLEFT_BRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(100)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&67166080) != 0 {
			{
				p.SetState(97)
				p.Type_field()
			}

			p.SetState(102)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(103)
			p.Match(TParserRIGHT_BRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(104)
			p.Match(TParserKW_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(105)
			p.Match(TParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(106)
			p.Match(TParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(107)
			p.Match(TParserLESS_THAN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(108)
			p.Generic_type()
		}
		{
			p.SetState(109)
			p.Match(TParserGREATER_THAN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IType_fieldContext is an interface to support dynamic dispatch.
type IType_fieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Common_type_field() ICommon_type_fieldContext
	Embed_type_field() IEmbed_type_fieldContext

	// IsType_fieldContext differentiates from other interfaces.
	IsType_fieldContext()
}

type Type_fieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyType_fieldContext() *Type_fieldContext {
	var p = new(Type_fieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_type_field
	return p
}

func InitEmptyType_fieldContext(p *Type_fieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_type_field
}

func (*Type_fieldContext) IsType_fieldContext() {}

func NewType_fieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Type_fieldContext {
	var p = new(Type_fieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_type_field

	return p
}

func (s *Type_fieldContext) GetParser() antlr.Parser { return s.parser }

func (s *Type_fieldContext) Common_type_field() ICommon_type_fieldContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICommon_type_fieldContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICommon_type_fieldContext)
}

func (s *Type_fieldContext) Embed_type_field() IEmbed_type_fieldContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEmbed_type_fieldContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEmbed_type_fieldContext)
}

func (s *Type_fieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Type_fieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Type_fieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterType_field(s)
	}
}

func (s *Type_fieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitType_field(s)
	}
}

func (p *TParser) Type_field() (localctx IType_fieldContext) {
	localctx = NewType_fieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, TParserRULE_type_field)
	p.SetState(115)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(113)
			p.Common_type_field()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(114)
			p.Embed_type_field()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEmbed_type_fieldContext is an interface to support dynamic dispatch.
type IEmbed_type_fieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	User_type() IUser_typeContext

	// IsEmbed_type_fieldContext differentiates from other interfaces.
	IsEmbed_type_fieldContext()
}

type Embed_type_fieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEmbed_type_fieldContext() *Embed_type_fieldContext {
	var p = new(Embed_type_fieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_embed_type_field
	return p
}

func InitEmptyEmbed_type_fieldContext(p *Embed_type_fieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_embed_type_field
}

func (*Embed_type_fieldContext) IsEmbed_type_fieldContext() {}

func NewEmbed_type_fieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Embed_type_fieldContext {
	var p = new(Embed_type_fieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_embed_type_field

	return p
}

func (s *Embed_type_fieldContext) GetParser() antlr.Parser { return s.parser }

func (s *Embed_type_fieldContext) User_type() IUser_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUser_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUser_typeContext)
}

func (s *Embed_type_fieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Embed_type_fieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Embed_type_fieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterEmbed_type_field(s)
	}
}

func (s *Embed_type_fieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitEmbed_type_field(s)
	}
}

func (p *TParser) Embed_type_field() (localctx IEmbed_type_fieldContext) {
	localctx = NewEmbed_type_fieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, TParserRULE_embed_type_field)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(117)
		p.User_type()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICommon_type_fieldContext is an interface to support dynamic dispatch.
type ICommon_type_fieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Common_field_type() ICommon_field_typeContext
	IDENTIFIER() antlr.TerminalNode
	EQUAL() antlr.TerminalNode
	Const_value() IConst_valueContext
	Type_annotations() IType_annotationsContext

	// IsCommon_type_fieldContext differentiates from other interfaces.
	IsCommon_type_fieldContext()
}

type Common_type_fieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommon_type_fieldContext() *Common_type_fieldContext {
	var p = new(Common_type_fieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_common_type_field
	return p
}

func InitEmptyCommon_type_fieldContext(p *Common_type_fieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_common_type_field
}

func (*Common_type_fieldContext) IsCommon_type_fieldContext() {}

func NewCommon_type_fieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Common_type_fieldContext {
	var p = new(Common_type_fieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_common_type_field

	return p
}

func (s *Common_type_fieldContext) GetParser() antlr.Parser { return s.parser }

func (s *Common_type_fieldContext) Common_field_type() ICommon_field_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICommon_field_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICommon_field_typeContext)
}

func (s *Common_type_fieldContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(TParserIDENTIFIER, 0)
}

func (s *Common_type_fieldContext) EQUAL() antlr.TerminalNode {
	return s.GetToken(TParserEQUAL, 0)
}

func (s *Common_type_fieldContext) Const_value() IConst_valueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConst_valueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConst_valueContext)
}

func (s *Common_type_fieldContext) Type_annotations() IType_annotationsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IType_annotationsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IType_annotationsContext)
}

func (s *Common_type_fieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Common_type_fieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Common_type_fieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterCommon_type_field(s)
	}
}

func (s *Common_type_fieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitCommon_type_field(s)
	}
}

func (p *TParser) Common_type_field() (localctx ICommon_type_fieldContext) {
	localctx = NewCommon_type_fieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, TParserRULE_common_type_field)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(119)
		p.Common_field_type()
	}
	{
		p.SetState(120)
		p.Match(TParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(123)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == TParserEQUAL {
		{
			p.SetState(121)
			p.Match(TParserEQUAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(122)
			p.Const_value()
		}

	}
	p.SetState(126)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == TParserLEFT_PAREN {
		{
			p.SetState(125)
			p.Type_annotations()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICommon_field_typeContext is an interface to support dynamic dispatch.
type ICommon_field_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPE_ANY() antlr.TerminalNode
	Base_type() IBase_typeContext
	User_type() IUser_typeContext
	Container_type() IContainer_typeContext
	TYPE_BINARY() antlr.TerminalNode

	// IsCommon_field_typeContext differentiates from other interfaces.
	IsCommon_field_typeContext()
}

type Common_field_typeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommon_field_typeContext() *Common_field_typeContext {
	var p = new(Common_field_typeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_common_field_type
	return p
}

func InitEmptyCommon_field_typeContext(p *Common_field_typeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_common_field_type
}

func (*Common_field_typeContext) IsCommon_field_typeContext() {}

func NewCommon_field_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Common_field_typeContext {
	var p = new(Common_field_typeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_common_field_type

	return p
}

func (s *Common_field_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Common_field_typeContext) TYPE_ANY() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_ANY, 0)
}

func (s *Common_field_typeContext) Base_type() IBase_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBase_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBase_typeContext)
}

func (s *Common_field_typeContext) User_type() IUser_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUser_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUser_typeContext)
}

func (s *Common_field_typeContext) Container_type() IContainer_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IContainer_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IContainer_typeContext)
}

func (s *Common_field_typeContext) TYPE_BINARY() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_BINARY, 0)
}

func (s *Common_field_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Common_field_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Common_field_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterCommon_field_type(s)
	}
}

func (s *Common_field_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitCommon_field_type(s)
	}
}

func (p *TParser) Common_field_type() (localctx ICommon_field_typeContext) {
	localctx = NewCommon_field_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, TParserRULE_common_field_type)
	p.SetState(133)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TParserTYPE_ANY:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(128)
			p.Match(TParserTYPE_ANY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case TParserTYPE_BOOL, TParserTYPE_INT, TParserTYPE_FLOAT, TParserTYPE_STRING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(129)
			p.Base_type()
		}

	case TParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(130)
			p.User_type()
		}

	case TParserTYPE_MAP, TParserTYPE_LIST:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(131)
			p.Container_type()
		}

	case TParserTYPE_BINARY:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(132)
			p.Match(TParserTYPE_BINARY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IGeneric_typeContext is an interface to support dynamic dispatch.
type IGeneric_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Base_type() IBase_typeContext
	User_type() IUser_typeContext
	Container_type() IContainer_typeContext

	// IsGeneric_typeContext differentiates from other interfaces.
	IsGeneric_typeContext()
}

type Generic_typeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGeneric_typeContext() *Generic_typeContext {
	var p = new(Generic_typeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_generic_type
	return p
}

func InitEmptyGeneric_typeContext(p *Generic_typeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_generic_type
}

func (*Generic_typeContext) IsGeneric_typeContext() {}

func NewGeneric_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Generic_typeContext {
	var p = new(Generic_typeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_generic_type

	return p
}

func (s *Generic_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Generic_typeContext) Base_type() IBase_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBase_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBase_typeContext)
}

func (s *Generic_typeContext) User_type() IUser_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUser_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUser_typeContext)
}

func (s *Generic_typeContext) Container_type() IContainer_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IContainer_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IContainer_typeContext)
}

func (s *Generic_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Generic_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Generic_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterGeneric_type(s)
	}
}

func (s *Generic_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitGeneric_type(s)
	}
}

func (p *TParser) Generic_type() (localctx IGeneric_typeContext) {
	localctx = NewGeneric_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, TParserRULE_generic_type)
	p.SetState(138)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TParserTYPE_BOOL, TParserTYPE_INT, TParserTYPE_FLOAT, TParserTYPE_STRING:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(135)
			p.Base_type()
		}

	case TParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(136)
			p.User_type()
		}

	case TParserTYPE_MAP, TParserTYPE_LIST:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(137)
			p.Container_type()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IType_annotationsContext is an interface to support dynamic dispatch.
type IType_annotationsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LEFT_PAREN() antlr.TerminalNode
	AllAnnotation() []IAnnotationContext
	Annotation(i int) IAnnotationContext
	RIGHT_PAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsType_annotationsContext differentiates from other interfaces.
	IsType_annotationsContext()
}

type Type_annotationsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyType_annotationsContext() *Type_annotationsContext {
	var p = new(Type_annotationsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_type_annotations
	return p
}

func InitEmptyType_annotationsContext(p *Type_annotationsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_type_annotations
}

func (*Type_annotationsContext) IsType_annotationsContext() {}

func NewType_annotationsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Type_annotationsContext {
	var p = new(Type_annotationsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_type_annotations

	return p
}

func (s *Type_annotationsContext) GetParser() antlr.Parser { return s.parser }

func (s *Type_annotationsContext) LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(TParserLEFT_PAREN, 0)
}

func (s *Type_annotationsContext) AllAnnotation() []IAnnotationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAnnotationContext); ok {
			len++
		}
	}

	tst := make([]IAnnotationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAnnotationContext); ok {
			tst[i] = t.(IAnnotationContext)
			i++
		}
	}

	return tst
}

func (s *Type_annotationsContext) Annotation(i int) IAnnotationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAnnotationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAnnotationContext)
}

func (s *Type_annotationsContext) RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(TParserRIGHT_PAREN, 0)
}

func (s *Type_annotationsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(TParserCOMMA)
}

func (s *Type_annotationsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(TParserCOMMA, i)
}

func (s *Type_annotationsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Type_annotationsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Type_annotationsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterType_annotations(s)
	}
}

func (s *Type_annotationsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitType_annotations(s)
	}
}

func (p *TParser) Type_annotations() (localctx IType_annotationsContext) {
	localctx = NewType_annotationsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, TParserRULE_type_annotations)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(140)
		p.Match(TParserLEFT_PAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(141)
		p.Annotation()
	}
	p.SetState(146)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == TParserCOMMA {
		{
			p.SetState(142)
			p.Match(TParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(143)
			p.Annotation()
		}

		p.SetState(148)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(149)
		p.Match(TParserRIGHT_PAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRpc_defContext is an interface to support dynamic dispatch.
type IRpc_defContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	KW_RPC() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LEFT_PAREN() antlr.TerminalNode
	Rpc_req() IRpc_reqContext
	RIGHT_PAREN() antlr.TerminalNode
	Rpc_resp() IRpc_respContext
	Rpc_annotations() IRpc_annotationsContext

	// IsRpc_defContext differentiates from other interfaces.
	IsRpc_defContext()
}

type Rpc_defContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRpc_defContext() *Rpc_defContext {
	var p = new(Rpc_defContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_rpc_def
	return p
}

func InitEmptyRpc_defContext(p *Rpc_defContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_rpc_def
}

func (*Rpc_defContext) IsRpc_defContext() {}

func NewRpc_defContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Rpc_defContext {
	var p = new(Rpc_defContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_rpc_def

	return p
}

func (s *Rpc_defContext) GetParser() antlr.Parser { return s.parser }

func (s *Rpc_defContext) KW_RPC() antlr.TerminalNode {
	return s.GetToken(TParserKW_RPC, 0)
}

func (s *Rpc_defContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(TParserIDENTIFIER, 0)
}

func (s *Rpc_defContext) LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(TParserLEFT_PAREN, 0)
}

func (s *Rpc_defContext) Rpc_req() IRpc_reqContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRpc_reqContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRpc_reqContext)
}

func (s *Rpc_defContext) RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(TParserRIGHT_PAREN, 0)
}

func (s *Rpc_defContext) Rpc_resp() IRpc_respContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRpc_respContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRpc_respContext)
}

func (s *Rpc_defContext) Rpc_annotations() IRpc_annotationsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRpc_annotationsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRpc_annotationsContext)
}

func (s *Rpc_defContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Rpc_defContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Rpc_defContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterRpc_def(s)
	}
}

func (s *Rpc_defContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitRpc_def(s)
	}
}

func (p *TParser) Rpc_def() (localctx IRpc_defContext) {
	localctx = NewRpc_defContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, TParserRULE_rpc_def)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(151)
		p.Match(TParserKW_RPC)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(152)
		p.Match(TParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(153)
		p.Match(TParserLEFT_PAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(154)
		p.Rpc_req()
	}
	{
		p.SetState(155)
		p.Match(TParserRIGHT_PAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(156)
		p.Rpc_resp()
	}
	{
		p.SetState(157)
		p.Rpc_annotations()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRpc_reqContext is an interface to support dynamic dispatch.
type IRpc_reqContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode

	// IsRpc_reqContext differentiates from other interfaces.
	IsRpc_reqContext()
}

type Rpc_reqContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRpc_reqContext() *Rpc_reqContext {
	var p = new(Rpc_reqContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_rpc_req
	return p
}

func InitEmptyRpc_reqContext(p *Rpc_reqContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_rpc_req
}

func (*Rpc_reqContext) IsRpc_reqContext() {}

func NewRpc_reqContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Rpc_reqContext {
	var p = new(Rpc_reqContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_rpc_req

	return p
}

func (s *Rpc_reqContext) GetParser() antlr.Parser { return s.parser }

func (s *Rpc_reqContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(TParserIDENTIFIER, 0)
}

func (s *Rpc_reqContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Rpc_reqContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Rpc_reqContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterRpc_req(s)
	}
}

func (s *Rpc_reqContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitRpc_req(s)
	}
}

func (p *TParser) Rpc_req() (localctx IRpc_reqContext) {
	localctx = NewRpc_reqContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, TParserRULE_rpc_req)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(159)
		p.Match(TParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRpc_respContext is an interface to support dynamic dispatch.
type IRpc_respContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	TYPE_STREAM() antlr.TerminalNode
	LESS_THAN() antlr.TerminalNode
	User_type() IUser_typeContext
	GREATER_THAN() antlr.TerminalNode

	// IsRpc_respContext differentiates from other interfaces.
	IsRpc_respContext()
}

type Rpc_respContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRpc_respContext() *Rpc_respContext {
	var p = new(Rpc_respContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_rpc_resp
	return p
}

func InitEmptyRpc_respContext(p *Rpc_respContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_rpc_resp
}

func (*Rpc_respContext) IsRpc_respContext() {}

func NewRpc_respContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Rpc_respContext {
	var p = new(Rpc_respContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_rpc_resp

	return p
}

func (s *Rpc_respContext) GetParser() antlr.Parser { return s.parser }

func (s *Rpc_respContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(TParserIDENTIFIER, 0)
}

func (s *Rpc_respContext) TYPE_STREAM() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_STREAM, 0)
}

func (s *Rpc_respContext) LESS_THAN() antlr.TerminalNode {
	return s.GetToken(TParserLESS_THAN, 0)
}

func (s *Rpc_respContext) User_type() IUser_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUser_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUser_typeContext)
}

func (s *Rpc_respContext) GREATER_THAN() antlr.TerminalNode {
	return s.GetToken(TParserGREATER_THAN, 0)
}

func (s *Rpc_respContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Rpc_respContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Rpc_respContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterRpc_resp(s)
	}
}

func (s *Rpc_respContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitRpc_resp(s)
	}
}

func (p *TParser) Rpc_resp() (localctx IRpc_respContext) {
	localctx = NewRpc_respContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, TParserRULE_rpc_resp)
	p.SetState(167)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(161)
			p.Match(TParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case TParserTYPE_STREAM:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(162)
			p.Match(TParserTYPE_STREAM)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(163)
			p.Match(TParserLESS_THAN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(164)
			p.User_type()
		}
		{
			p.SetState(165)
			p.Match(TParserGREATER_THAN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRpc_annotationsContext is an interface to support dynamic dispatch.
type IRpc_annotationsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LEFT_BRACE() antlr.TerminalNode
	RIGHT_BRACE() antlr.TerminalNode
	AllAnnotation() []IAnnotationContext
	Annotation(i int) IAnnotationContext

	// IsRpc_annotationsContext differentiates from other interfaces.
	IsRpc_annotationsContext()
}

type Rpc_annotationsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRpc_annotationsContext() *Rpc_annotationsContext {
	var p = new(Rpc_annotationsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_rpc_annotations
	return p
}

func InitEmptyRpc_annotationsContext(p *Rpc_annotationsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_rpc_annotations
}

func (*Rpc_annotationsContext) IsRpc_annotationsContext() {}

func NewRpc_annotationsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Rpc_annotationsContext {
	var p = new(Rpc_annotationsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_rpc_annotations

	return p
}

func (s *Rpc_annotationsContext) GetParser() antlr.Parser { return s.parser }

func (s *Rpc_annotationsContext) LEFT_BRACE() antlr.TerminalNode {
	return s.GetToken(TParserLEFT_BRACE, 0)
}

func (s *Rpc_annotationsContext) RIGHT_BRACE() antlr.TerminalNode {
	return s.GetToken(TParserRIGHT_BRACE, 0)
}

func (s *Rpc_annotationsContext) AllAnnotation() []IAnnotationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAnnotationContext); ok {
			len++
		}
	}

	tst := make([]IAnnotationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAnnotationContext); ok {
			tst[i] = t.(IAnnotationContext)
			i++
		}
	}

	return tst
}

func (s *Rpc_annotationsContext) Annotation(i int) IAnnotationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAnnotationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAnnotationContext)
}

func (s *Rpc_annotationsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Rpc_annotationsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Rpc_annotationsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterRpc_annotations(s)
	}
}

func (s *Rpc_annotationsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitRpc_annotations(s)
	}
}

func (p *TParser) Rpc_annotations() (localctx IRpc_annotationsContext) {
	localctx = NewRpc_annotationsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, TParserRULE_rpc_annotations)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(169)
		p.Match(TParserLEFT_BRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(173)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == TParserIDENTIFIER {
		{
			p.SetState(170)
			p.Annotation()
		}

		p.SetState(175)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(176)
		p.Match(TParserRIGHT_BRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAnnotationContext is an interface to support dynamic dispatch.
type IAnnotationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	EQUAL() antlr.TerminalNode
	Const_value() IConst_valueContext

	// IsAnnotationContext differentiates from other interfaces.
	IsAnnotationContext()
}

type AnnotationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAnnotationContext() *AnnotationContext {
	var p = new(AnnotationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_annotation
	return p
}

func InitEmptyAnnotationContext(p *AnnotationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_annotation
}

func (*AnnotationContext) IsAnnotationContext() {}

func NewAnnotationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AnnotationContext {
	var p = new(AnnotationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_annotation

	return p
}

func (s *AnnotationContext) GetParser() antlr.Parser { return s.parser }

func (s *AnnotationContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(TParserIDENTIFIER, 0)
}

func (s *AnnotationContext) EQUAL() antlr.TerminalNode {
	return s.GetToken(TParserEQUAL, 0)
}

func (s *AnnotationContext) Const_value() IConst_valueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConst_valueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConst_valueContext)
}

func (s *AnnotationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnnotationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AnnotationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterAnnotation(s)
	}
}

func (s *AnnotationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitAnnotation(s)
	}
}

func (p *TParser) Annotation() (localctx IAnnotationContext) {
	localctx = NewAnnotationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, TParserRULE_annotation)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(178)
		p.Match(TParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(181)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == TParserEQUAL {
		{
			p.SetState(179)
			p.Match(TParserEQUAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(180)
			p.Const_value()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBase_typeContext is an interface to support dynamic dispatch.
type IBase_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPE_BOOL() antlr.TerminalNode
	TYPE_INT() antlr.TerminalNode
	TYPE_FLOAT() antlr.TerminalNode
	TYPE_STRING() antlr.TerminalNode
	QUESTION() antlr.TerminalNode

	// IsBase_typeContext differentiates from other interfaces.
	IsBase_typeContext()
}

type Base_typeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBase_typeContext() *Base_typeContext {
	var p = new(Base_typeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_base_type
	return p
}

func InitEmptyBase_typeContext(p *Base_typeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_base_type
}

func (*Base_typeContext) IsBase_typeContext() {}

func NewBase_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Base_typeContext {
	var p = new(Base_typeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_base_type

	return p
}

func (s *Base_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Base_typeContext) TYPE_BOOL() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_BOOL, 0)
}

func (s *Base_typeContext) TYPE_INT() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_INT, 0)
}

func (s *Base_typeContext) TYPE_FLOAT() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_FLOAT, 0)
}

func (s *Base_typeContext) TYPE_STRING() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_STRING, 0)
}

func (s *Base_typeContext) QUESTION() antlr.TerminalNode {
	return s.GetToken(TParserQUESTION, 0)
}

func (s *Base_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Base_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Base_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterBase_type(s)
	}
}

func (s *Base_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitBase_type(s)
	}
}

func (p *TParser) Base_type() (localctx IBase_typeContext) {
	localctx = NewBase_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, TParserRULE_base_type)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(183)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&3840) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(185)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == TParserQUESTION {
		{
			p.SetState(184)
			p.Match(TParserQUESTION)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUser_typeContext is an interface to support dynamic dispatch.
type IUser_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	QUESTION() antlr.TerminalNode

	// IsUser_typeContext differentiates from other interfaces.
	IsUser_typeContext()
}

type User_typeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUser_typeContext() *User_typeContext {
	var p = new(User_typeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_user_type
	return p
}

func InitEmptyUser_typeContext(p *User_typeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_user_type
}

func (*User_typeContext) IsUser_typeContext() {}

func NewUser_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *User_typeContext {
	var p = new(User_typeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_user_type

	return p
}

func (s *User_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *User_typeContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(TParserIDENTIFIER, 0)
}

func (s *User_typeContext) QUESTION() antlr.TerminalNode {
	return s.GetToken(TParserQUESTION, 0)
}

func (s *User_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *User_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *User_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterUser_type(s)
	}
}

func (s *User_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitUser_type(s)
	}
}

func (p *TParser) User_type() (localctx IUser_typeContext) {
	localctx = NewUser_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, TParserRULE_user_type)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(187)
		p.Match(TParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(189)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == TParserQUESTION {
		{
			p.SetState(188)
			p.Match(TParserQUESTION)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IContainer_typeContext is an interface to support dynamic dispatch.
type IContainer_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Map_type() IMap_typeContext
	List_type() IList_typeContext

	// IsContainer_typeContext differentiates from other interfaces.
	IsContainer_typeContext()
}

type Container_typeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyContainer_typeContext() *Container_typeContext {
	var p = new(Container_typeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_container_type
	return p
}

func InitEmptyContainer_typeContext(p *Container_typeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_container_type
}

func (*Container_typeContext) IsContainer_typeContext() {}

func NewContainer_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Container_typeContext {
	var p = new(Container_typeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_container_type

	return p
}

func (s *Container_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Container_typeContext) Map_type() IMap_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMap_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMap_typeContext)
}

func (s *Container_typeContext) List_type() IList_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IList_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IList_typeContext)
}

func (s *Container_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Container_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Container_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterContainer_type(s)
	}
}

func (s *Container_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitContainer_type(s)
	}
}

func (p *TParser) Container_type() (localctx IContainer_typeContext) {
	localctx = NewContainer_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, TParserRULE_container_type)
	p.SetState(193)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TParserTYPE_MAP:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(191)
			p.Map_type()
		}

	case TParserTYPE_LIST:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(192)
			p.List_type()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMap_typeContext is an interface to support dynamic dispatch.
type IMap_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPE_MAP() antlr.TerminalNode
	LESS_THAN() antlr.TerminalNode
	Key_type() IKey_typeContext
	COMMA() antlr.TerminalNode
	Value_type() IValue_typeContext
	GREATER_THAN() antlr.TerminalNode

	// IsMap_typeContext differentiates from other interfaces.
	IsMap_typeContext()
}

type Map_typeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMap_typeContext() *Map_typeContext {
	var p = new(Map_typeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_map_type
	return p
}

func InitEmptyMap_typeContext(p *Map_typeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_map_type
}

func (*Map_typeContext) IsMap_typeContext() {}

func NewMap_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Map_typeContext {
	var p = new(Map_typeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_map_type

	return p
}

func (s *Map_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Map_typeContext) TYPE_MAP() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_MAP, 0)
}

func (s *Map_typeContext) LESS_THAN() antlr.TerminalNode {
	return s.GetToken(TParserLESS_THAN, 0)
}

func (s *Map_typeContext) Key_type() IKey_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IKey_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IKey_typeContext)
}

func (s *Map_typeContext) COMMA() antlr.TerminalNode {
	return s.GetToken(TParserCOMMA, 0)
}

func (s *Map_typeContext) Value_type() IValue_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValue_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValue_typeContext)
}

func (s *Map_typeContext) GREATER_THAN() antlr.TerminalNode {
	return s.GetToken(TParserGREATER_THAN, 0)
}

func (s *Map_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Map_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Map_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterMap_type(s)
	}
}

func (s *Map_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitMap_type(s)
	}
}

func (p *TParser) Map_type() (localctx IMap_typeContext) {
	localctx = NewMap_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, TParserRULE_map_type)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(195)
		p.Match(TParserTYPE_MAP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(196)
		p.Match(TParserLESS_THAN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(197)
		p.Key_type()
	}
	{
		p.SetState(198)
		p.Match(TParserCOMMA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(199)
		p.Value_type()
	}
	{
		p.SetState(200)
		p.Match(TParserGREATER_THAN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IKey_typeContext is an interface to support dynamic dispatch.
type IKey_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPE_STRING() antlr.TerminalNode
	TYPE_INT() antlr.TerminalNode

	// IsKey_typeContext differentiates from other interfaces.
	IsKey_typeContext()
}

type Key_typeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyKey_typeContext() *Key_typeContext {
	var p = new(Key_typeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_key_type
	return p
}

func InitEmptyKey_typeContext(p *Key_typeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_key_type
}

func (*Key_typeContext) IsKey_typeContext() {}

func NewKey_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Key_typeContext {
	var p = new(Key_typeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_key_type

	return p
}

func (s *Key_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Key_typeContext) TYPE_STRING() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_STRING, 0)
}

func (s *Key_typeContext) TYPE_INT() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_INT, 0)
}

func (s *Key_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Key_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Key_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterKey_type(s)
	}
}

func (s *Key_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitKey_type(s)
	}
}

func (p *TParser) Key_type() (localctx IKey_typeContext) {
	localctx = NewKey_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, TParserRULE_key_type)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(202)
		_la = p.GetTokenStream().LA(1)

		if !(_la == TParserTYPE_INT || _la == TParserTYPE_STRING) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IList_typeContext is an interface to support dynamic dispatch.
type IList_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPE_LIST() antlr.TerminalNode
	LESS_THAN() antlr.TerminalNode
	Value_type() IValue_typeContext
	GREATER_THAN() antlr.TerminalNode

	// IsList_typeContext differentiates from other interfaces.
	IsList_typeContext()
}

type List_typeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyList_typeContext() *List_typeContext {
	var p = new(List_typeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_list_type
	return p
}

func InitEmptyList_typeContext(p *List_typeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_list_type
}

func (*List_typeContext) IsList_typeContext() {}

func NewList_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *List_typeContext {
	var p = new(List_typeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_list_type

	return p
}

func (s *List_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *List_typeContext) TYPE_LIST() antlr.TerminalNode {
	return s.GetToken(TParserTYPE_LIST, 0)
}

func (s *List_typeContext) LESS_THAN() antlr.TerminalNode {
	return s.GetToken(TParserLESS_THAN, 0)
}

func (s *List_typeContext) Value_type() IValue_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValue_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValue_typeContext)
}

func (s *List_typeContext) GREATER_THAN() antlr.TerminalNode {
	return s.GetToken(TParserGREATER_THAN, 0)
}

func (s *List_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *List_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *List_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterList_type(s)
	}
}

func (s *List_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitList_type(s)
	}
}

func (p *TParser) List_type() (localctx IList_typeContext) {
	localctx = NewList_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, TParserRULE_list_type)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(204)
		p.Match(TParserTYPE_LIST)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(205)
		p.Match(TParserLESS_THAN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(206)
		p.Value_type()
	}
	{
		p.SetState(207)
		p.Match(TParserGREATER_THAN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IValue_typeContext is an interface to support dynamic dispatch.
type IValue_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Base_type() IBase_typeContext
	User_type() IUser_typeContext
	Container_type() IContainer_typeContext

	// IsValue_typeContext differentiates from other interfaces.
	IsValue_typeContext()
}

type Value_typeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValue_typeContext() *Value_typeContext {
	var p = new(Value_typeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_value_type
	return p
}

func InitEmptyValue_typeContext(p *Value_typeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_value_type
}

func (*Value_typeContext) IsValue_typeContext() {}

func NewValue_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Value_typeContext {
	var p = new(Value_typeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_value_type

	return p
}

func (s *Value_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Value_typeContext) Base_type() IBase_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBase_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBase_typeContext)
}

func (s *Value_typeContext) User_type() IUser_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUser_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUser_typeContext)
}

func (s *Value_typeContext) Container_type() IContainer_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IContainer_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IContainer_typeContext)
}

func (s *Value_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Value_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Value_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterValue_type(s)
	}
}

func (s *Value_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitValue_type(s)
	}
}

func (p *TParser) Value_type() (localctx IValue_typeContext) {
	localctx = NewValue_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, TParserRULE_value_type)
	p.SetState(212)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case TParserTYPE_BOOL, TParserTYPE_INT, TParserTYPE_FLOAT, TParserTYPE_STRING:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(209)
			p.Base_type()
		}

	case TParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(210)
			p.User_type()
		}

	case TParserTYPE_MAP, TParserTYPE_LIST:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(211)
			p.Container_type()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConst_valueContext is an interface to support dynamic dispatch.
type IConst_valueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	KW_TRUE() antlr.TerminalNode
	KW_FALSE() antlr.TerminalNode
	INTEGER() antlr.TerminalNode
	FLOAT() antlr.TerminalNode
	STRING() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsConst_valueContext differentiates from other interfaces.
	IsConst_valueContext()
}

type Const_valueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConst_valueContext() *Const_valueContext {
	var p = new(Const_valueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_const_value
	return p
}

func InitEmptyConst_valueContext(p *Const_valueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = TParserRULE_const_value
}

func (*Const_valueContext) IsConst_valueContext() {}

func NewConst_valueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Const_valueContext {
	var p = new(Const_valueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = TParserRULE_const_value

	return p
}

func (s *Const_valueContext) GetParser() antlr.Parser { return s.parser }

func (s *Const_valueContext) KW_TRUE() antlr.TerminalNode {
	return s.GetToken(TParserKW_TRUE, 0)
}

func (s *Const_valueContext) KW_FALSE() antlr.TerminalNode {
	return s.GetToken(TParserKW_FALSE, 0)
}

func (s *Const_valueContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(TParserINTEGER, 0)
}

func (s *Const_valueContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(TParserFLOAT, 0)
}

func (s *Const_valueContext) STRING() antlr.TerminalNode {
	return s.GetToken(TParserSTRING, 0)
}

func (s *Const_valueContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(TParserIDENTIFIER, 0)
}

func (s *Const_valueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Const_valueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Const_valueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.EnterConst_value(s)
	}
}

func (s *Const_valueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TParserListener); ok {
		listenerT.ExitConst_value(s)
	}
}

func (p *TParser) Const_value() (localctx IConst_valueContext) {
	localctx = NewConst_valueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, TParserRULE_const_value)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(214)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&503316576) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
