grammar RestPath;

// ----------------------
// Top-level rule: a path consists of multiple segments
// ----------------------
path
    : '/' segment ( '/' segment )*
    ;

// ----------------------
// Path segment types
// ----------------------
segment
    : STATIC_SEGMENT       // Static segment
    | PARAM_SEGMENT        // Colon-style :param or :param* (wildcard)
    | BRACED_PARAM         // Curly-brace style {param} or {param...} (wildcard)
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
PARAM_SEGMENT
    : ':' IDENTIFIER ('*')?
    ;

// ----------------------
// Curly-brace style parameter {param} or {param...} (wildcard)
// ----------------------
BRACED_PARAM
    : '{' IDENTIFIER ('...' )? '}'
    ;

// ----------------------
// Identifier for parameters, letters, digits, or underscore
// ----------------------
IDENTIFIER
    : [a-zA-Z_] [a-zA-Z0-9_]*
    ;
