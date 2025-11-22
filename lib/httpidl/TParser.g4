// --------------------
// Parser Grammar
// --------------------
parser grammar TParser;

options { tokenVocab = TLexer; }

// --------------------
// Document root
// A document consists of zero or more definitions separated by terminators
// and ends with EOF.
// --------------------
document
    : ((definition terminator) | terminator)* EOF
    ;

// --------------------
// Top-level definitions: const, enum, type, oneof, rpc
// --------------------
definition
    : const_def | enum_def | type_def | oneof_def | rpc_def | sse_def
    ;

// --------------------
// Constant definition
// Example: const string a = "1"
// --------------------
const_def
    : KW_CONST base_type IDENTIFIER EQUAL const_value
    ;

// --------------------
// Enum definition
// Example:
// enum A {
//   RED = 1
//   GREEN = 2
// }
// --------------------
enum_def
    : KW_ENUM IDENTIFIER LEFT_BRACE terminator? (enum_field terminator)* terminator? RIGHT_BRACE
    ;

// Enum field: name = integer
enum_field
    : IDENTIFIER EQUAL INTEGER
    ;

// --------------------
// Type definition
// Example 1:
// type A<T> {
//   B
//   string field = "1" (go.type="string")
// }
// Example 2:
// type Alias Map<string,User>
// --------------------
type_def
    // Structured type with optional generic parameter
    : KW_TYPE IDENTIFIER (LESS_THAN IDENTIFIER GREATER_THAN)? LEFT_BRACE terminator? (type_field terminator)* terminator? RIGHT_BRACE
    // Type alias to a generic container
    | KW_TYPE IDENTIFIER IDENTIFIER LESS_THAN value_type GREATER_THAN
    ;

// A type field can be either an embedded type or a named typed field
type_field
    : embed_type_field | common_type_field
    ;

// Embedded field: user-defined type
embed_type_field
    : user_type
    ;

// Common field: type + name + optional annotations
common_type_field
    : KW_REQUIRED? common_field_type IDENTIFIER type_annotations?
    ;

// Field type options
common_field_type
    : TYPE_ANY
    | base_type
    | user_type
    | container_type
    | TYPE_BINARY
    ;

// --------------------
// Field annotations
// Example: (go.type="string", db.index=true)
// Example: (
//     go.type="string"
//     db.index=true
// )
// --------------------
type_annotations
    : LEFT_PAREN terminator? annotation ((COMMA|terminator) annotation)* terminator? RIGHT_PAREN
    ;

// --------------------
// OneOf definition
// Example:
// oneof Value {
//   A
//   B
// }
// --------------------
oneof_def
    : KW_ONEOF IDENTIFIER LEFT_BRACE terminator? (user_type terminator)* terminator? RIGHT_BRACE
    ;

// --------------------
// RPC definition
// Example:
// rpc GetUser (ReqType) RespType { method="GET" }
// --------------------
rpc_def
    : KW_RPC IDENTIFIER LEFT_PAREN rpc_req RIGHT_PAREN rpc_resp rpc_annotations
    ;

// RPC request type: a user-defined type
rpc_req
    : user_type
    ;

// RPC response type: a user-defined type
rpc_resp
    : user_type
    ;

// RPC annotations (inside { ... })
rpc_annotations
    : LEFT_BRACE terminator? (annotation terminator)* terminator? RIGHT_BRACE
    ;

// --------------------
// SSE definition
// Example:
// sse GetUser (ReqType) RespType { method="GET" }
// --------------------
sse_def
    : KW_SSE IDENTIFIER LEFT_PAREN sse_req RIGHT_PAREN sse_resp sse_annotations
    ;

// SSE request type: a user-defined type
sse_req
    : user_type
    ;

// SSE response type: a user-defined type
sse_resp
    : user_type
    ;

// SSE annotations (inside { ... })
sse_annotations
    : LEFT_BRACE terminator? (annotation terminator)* terminator? RIGHT_BRACE
    ;

// Annotation key-value pair
// Example: method="GET"
annotation
    : IDENTIFIER (EQUAL const_value)?
    ;

// --------------------
// Base types
// Primitive base types
// --------------------
base_type
    : TYPE_BOOL | TYPE_INT | TYPE_FLOAT | TYPE_STRING
    ;

// User-defined type
user_type
    : IDENTIFIER
    ;

// --------------------
// Container types: map<K,V> or list<T>
// --------------------
container_type
    : map_type | list_type
    ;

// Map type: map<string,int>
map_type
   : TYPE_MAP LESS_THAN key_type COMMA value_type GREATER_THAN
   ;

// Map keys: only string or int
key_type
    : TYPE_STRING | TYPE_INT
    ;

// List type: list<User>
list_type
   : TYPE_LIST LESS_THAN value_type GREATER_THAN
   ;

// Allowed list/map value types
value_type
    : base_type | user_type | container_type
    ;

// --------------------
// Constant values: literals or identifiers (e.g., enum members)
// --------------------
const_value
    : KW_TRUE | KW_FALSE | INTEGER | FLOAT | STRING | IDENTIFIER
    ;

// --------------------
// Terminator
// Terminator is used to separate statements or fields.
// It allows either one or more newlines.
// --------------------
terminator
    : (NEWLINE)+
    ;
