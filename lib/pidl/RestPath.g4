grammar RestPath;

// ----------------------
// Top-level rule: a path consists of multiple segments
// ----------------------
path
    : SLASH segment ( SLASH segment )*
    ;

// ----------------------
// Path segment types
// ----------------------
segment
    : STATIC_SEGMENT       // Static segment
    | paramSegment         // Colon-style :param or :param* (wildcard)
    | bracedParam          // Curly-brace style {param} or {param...} (wildcard)
    ;

// ----------------------
// Static path segment, e.g., "users", "books"
// ----------------------
STATIC_SEGMENT
    : [a-zA-Z0-9_-]+
    ;

// ----------------------
// Colon-style parameter :param or :param* (wildcard)
// ----------------------
paramSegment
    : COLON name=IDENTIFIER (wildcard=STAR)?
    ;

// ----------------------
// Curly-brace style parameter {param} or {param...} (wildcard)
// ----------------------
bracedParam
    : LBRACE name=IDENTIFIER (wildcard=ELLIPSIS)? RBRACE
    ;

// ----------------------
// Identifier for parameters, letters, digits, or underscore
// ----------------------
IDENTIFIER
    : [a-zA-Z_] [a-zA-Z0-9_]*
    ;

// --- punctuation tokens ---
COLON      : ':' ;
SLASH      : '/' ;
LBRACE     : '{' ;
RBRACE     : '}' ;
STAR       : '*' ;
ELLIPSIS   : '...' ;