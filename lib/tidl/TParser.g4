// --------------------
// Parser Grammar
// --------------------
parser grammar TParser;

options { tokenVocab = TLexer; }

// --------------------
// Document root
// --------------------
document
    : definition* EOF
    ;

// --------------------
// Definition types: const, enum, type, rpc
// --------------------
definition
    : const_def | enum_def | type_def | oneof_def | rpc_def
    ;

// --------------------
// Constant definition
// Example: const string a = "1"
// --------------------
const_def
    : KW_CONST const_type IDENTIFIER EQUAL const_value
    ;

// Allowed constant types: bool, int, float, or string
const_type
    : TYPE_BOOL | TYPE_INT | TYPE_FLOAT | TYPE_STRING
    ;

// --------------------
// Enum definition
// Example: enum A { A = 1 }
// --------------------
enum_def
    : KW_ENUM IDENTIFIER LEFT_BRACE enum_field* RIGHT_BRACE
    ;

// Enum field
enum_field
    : IDENTIFIER EQUAL INTEGER
    ;

// --------------------
// Type definition
// Example:
// type A<T> {
//   B?
//   string? field = "1" ( go.type="string" )
// }
// type Alias Map<string,User>
// --------------------
type_def
    : KW_TYPE IDENTIFIER (LESS_THAN IDENTIFIER GREATER_THAN)? LEFT_BRACE type_field* RIGHT_BRACE
    | KW_TYPE IDENTIFIER IDENTIFIER LESS_THAN generic_type GREATER_THAN
    ;

// A type field can be either an embedded type or a named typed field
type_field
    : common_type_field | embed_type_field
    ;

// Embedded field: user-defined type (optionally nullable with '?')
embed_type_field
    : '@'user_type
    ;

// Common field: type + name + optional default value + optional annotations
common_type_field
    : common_field_type IDENTIFIER (EQUAL const_value)? type_annotations?
    ;

// Field type options
common_field_type
    : TYPE_ANY
    | base_type
    | user_type
    | container_type
    | TYPE_BINARY
    ;

// Generic type
generic_type
    : base_type | user_type | container_type
    ;

// --------------------
// Field annotations
// Example: ( go.type="string", db.index=true )
// --------------------
type_annotations
    : LEFT_PAREN annotation (COMMA annotation)* RIGHT_PAREN
    ;

// --------------------
// OneOf definition
// Example:
// oneof Value {
//     A? a
//     B? b
// }
// --------------------
oneof_def
    : KW_ONEOF IDENTIFIER LEFT_BRACE oneof_field* RIGHT_BRACE
    ;

// OneOf fields must be normal named fields
oneof_field
    : common_type_field
    ;

// --------------------
// RPC definition
// Example:
// rpc GetUser (ReqType) RespType { method="GET" }
// --------------------
rpc_def
    : KW_RPC IDENTIFIER LEFT_PAREN rpc_req RIGHT_PAREN rpc_resp rpc_annotations
    ;

// RPC request type: always an identifier
rpc_req
    : IDENTIFIER
    ;

// RPC response type:
// Either an identifier, a generic form (Type<T>), or a stream<T>
rpc_resp
    : IDENTIFIER
    | TYPE_STREAM LESS_THAN user_type GREATER_THAN
    ;

// RPC annotations (inside { ... })
rpc_annotations
    : LEFT_BRACE annotation* RIGHT_BRACE
    ;

// Annotation for type or RPC
// Example: method="GET"
annotation
    : IDENTIFIER (EQUAL const_value)?
    ;

// --------------------
// Base types
// Primitive base types with optional nullable modifier '?'
// --------------------
base_type
    : (TYPE_BOOL | TYPE_INT | TYPE_FLOAT | TYPE_STRING) QUESTION?
    ;

// User-defined type (identifier, optionally nullable with '?')
user_type
    : IDENTIFIER QUESTION?
    ;

// --------------------
// Container types
// map<K,V> or list<T>
// --------------------
container_type
    : map_type | list_type
    ;

// Map type
// Example: map<string,int>
map_type
   : TYPE_MAP LESS_THAN key_type COMMA value_type GREATER_THAN
   ;

// Map keys: only string or int
key_type
    : TYPE_STRING | TYPE_INT
    ;

// List type
// Example: list<User>
list_type
   : TYPE_LIST LESS_THAN value_type GREATER_THAN
   ;

// Allowed list/map value types
value_type
    : base_type | user_type | container_type
    ;

// --------------------
// Constant values
// Can be literals (true, false, numbers, strings)
// Or identifiers (e.g., enum constants)
// --------------------
const_value
    : KW_TRUE | KW_FALSE | INTEGER | FLOAT | STRING | IDENTIFIER
    ;
