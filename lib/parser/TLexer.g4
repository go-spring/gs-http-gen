// --------------------
// Lexer Grammar
// --------------------
lexer grammar TLexer;

// Define additional channels for whitespace and comments.
// This ensures they don’t interfere with parsing, but can still be preserved if needed.
channels {WS_CHAN, SL_COMMENT_CHAN, ML_COMMENT_CHAN}

// --------------------
// Keywords
// --------------------
KW_CONST : 'const';
KW_ENUM  : 'enum';
KW_TYPE  : 'type';
KW_RPC   : 'rpc';
KW_TRUE  : 'true';
KW_FALSE : 'false';

// --------------------
// Basic types
// --------------------
TYPE_ANY    : 'any';
TYPE_BOOL   : 'bool';
TYPE_INT    : 'int';
TYPE_FLOAT  : 'float';
TYPE_STRING : 'string';
TYPE_BINARY : 'binary';
TYPE_STREAM : 'stream';

// --------------------
// Container types
// --------------------
TYPE_MAP  : 'map';
TYPE_LIST : 'list';

// --------------------
// Special symbols
// --------------------
LESS_THAN    : '<';
GREATER_THAN : '>';
LEFT_PAREN   : '(';
RIGHT_PAREN  : ')';
LEFT_BRACE   : '{';
RIGHT_BRACE  : '}';
EQUAL        : '=';
COMMA        : ',';
QUESTION     : '?';

// --------------------
// String literal
// Supports escape sequences (e.g., \" for quote, \\ for backslash)
// --------------------
STRING
    : '"' ( '\\' . | ~["\\] )* '"'
    ;

// --------------------
// Identifier
// Starts with a letter, followed by letters, digits, underscores, or dots
// --------------------
IDENTIFIER
    : LETTER (LETTER | DIGIT | '.' | '_')*
    ;

// --------------------
// Integer literal
// Decimal integer with optional sign (+/-) or hexadecimal integer prefixed with 0x.
// --------------------
INTEGER
    : ('+' | '-')? DIGIT+ | '0x' HEX_DIGIT+
    ;

// --------------------
// Floating-point number
// Supports decimals and scientific notation (e.g., 1.23e+10)
// --------------------
FLOAT
    : ('+' | '-')? ( DIGIT+ ('.' DIGIT+)? | '.' DIGIT+ ) (('E' | 'e') ('+'|'-')? DIGIT+ )?
    ;

// --------------------
// Fragments (used internally, not emitted as tokens)
// --------------------
fragment DIGIT     : '0'..'9';
fragment LETTER    : 'A'..'Z' | 'a'..'z';
fragment HEX_DIGIT : DIGIT | 'A'..'F' | 'a'..'f';

// --------------------
// Whitespace
// Skipped by sending to WS_CHAN
// --------------------
WHITESPACE
    : [ \t\r\n]+ -> channel(WS_CHAN)
    ;

// --------------------
// Single-line comments
// Supports both // and # styles
// --------------------
SINGLE_LINE_COMMENT
    : ('//' | '#') ~[\r\n]* ('\r'? '\n')? -> channel(SL_COMMENT_CHAN)
    ;

// --------------------
// Multi-line comments
// Supports /* ... */ with non-greedy matching
// --------------------
MULTI_LINE_COMMENT
    : '/*' .*? '*/' -> channel(ML_COMMENT_CHAN)
    ;
