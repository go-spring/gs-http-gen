package golang

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-spring/gs-http-gen/lib/httpidl"
	"github.com/go-spring/gs-http-gen/lib/pathidl"
	"github.com/go-spring/gs-http-gen/lib/validate"
	"github.com/lvan100/golib/errutil"
)

// TypeKind represents kind of a Go field type
type TypeKind int

const (
	TypeKindBool TypeKind = iota
	TypeKindInt
	TypeKindUint
	TypeKindFloat
	TypeKindString
	TypeKindStruct
	TypeKindEnum
	TypeKindList
	TypeKindMap
	TypeKindPointer
)

// Const represents a Go constant
type Const struct {
	Type    string
	Name    string
	Value   string
	Comment string
}

// Enum represents a Go enum
type Enum struct {
	Name    string
	Fields  []EnumField
	Comment string
}

// EnumField represents a single field in a Go enum
type EnumField struct {
	Name    string
	Value   int64
	Comment string
}

// Type represents a Go struct
type Type struct {
	Name        string
	Fields      []TypeField
	Comment     string
	Request     bool
	RequestBody bool
}

// TypeField represents a field in a Go struct
type TypeField struct {
	Name      string
	Type      string // for field
	TypeKind  []TypeKind
	ValueType string // for getter/setter
	FieldTag  string
	Required  bool
	JSONTag   httpidl.JSONTag
	FormTag   httpidl.FormTag
	Binding   *httpidl.Binding
	Validate  *string
	Comment   string
}

// IsPointer returns true if the field is a pointer
func (x TypeField) IsPointer() bool {
	return x.TypeKind[0] == TypeKindPointer
}

// BindingCount returns the number of fields in the struct that have binding info
func (t *Type) BindingCount() int {
	var count int
	for _, f := range t.Fields {
		if f.Binding != nil {
			count++
		}
	}
	return count
}

// QueryCount returns the number of fields in the struct that have query binding info
func (t *Type) QueryCount() int {
	var count int
	for _, f := range t.Fields {
		if f.Binding != nil && f.Binding.From == "query" {
			count++
		}
	}
	return count
}

// ValidateCount returns the number of fields in the struct that have validation expressions
func (t *Type) ValidateCount() int {
	var count int
	for _, f := range t.Fields {
		if f.Validate != nil {
			count++
		}
	}
	return count
}

// RPC represents a single remote procedure call with HTTP metadata.
type RPC struct {
	Name        string   // Method name
	Request     string   // Request type name
	Response    string   // Response type name
	Stream      bool     // Whether this RPC is a streaming RPC
	Path        string   // HTTP path
	FormatPath  string   // Formatted HTTP path
	PathParams  []string // HTTP path parameters
	Method      string   // HTTP method (GET, POST, etc.)
	ContentType string   // HTTP Content-Type
	Comment     string   // Comment of the RPC
}

type ReqIndex struct {
	File  string
	Index int
}

type GoCode struct {
	Files  map[string]httpidl.Document
	Meta   *httpidl.MetaInfo
	Reqs   map[string]ReqIndex // request type name
	Consts map[string][]Const
	Enums  map[string][]Enum
	Types  map[string][]Type
	RPCs   []RPC
	Funcs  map[string]ValidateFunc // Collected validation functions
}

func Convert(dir string) (GoCode, error) {
	project, err := httpidl.ParseDir(dir)
	if err != nil {
		return GoCode{}, err
	}

	code := GoCode{
		Files:  project.Files,
		Meta:   project.Meta,
		Reqs:   make(map[string]ReqIndex),
		Consts: make(map[string][]Const),
		Enums:  make(map[string][]Enum),
		Types:  make(map[string][]Type),
		Funcs:  make(map[string]ValidateFunc),
	}

	// Collect all RPC definitions
	for _, doc := range project.Files {
		for _, r := range doc.RPCs {
			rpc, err := convertRPC(r)
			if err != nil {
				return code, err
			}
			code.RPCs = append(code.RPCs, rpc)
			code.Reqs[rpc.Request] = ReqIndex{}
		}
	}
	sort.Slice(code.RPCs, func(i, j int) bool {
		return code.RPCs[i].Name < code.RPCs[j].Name
	})

	for fileName, doc := range project.Files {
		consts, err := convertConsts(code, doc)
		if err != nil {
			return code, errutil.Explain(nil, "convert consts error: %w", err)
		}
		enums, err := convertEnums(code, doc)
		if err != nil {
			return code, errutil.Explain(nil, "convert enums error: %w", err)
		}
		types, err := convertTypes(code, doc)
		if err != nil {
			return code, errutil.Explain(nil, "convert types error: %w", err)
		}
		{
			var temp []Type
			for _, t := range types {
				if _, ok := code.Reqs[t.Name]; ok {
					req, body := SplitRequestType(t)
					code.Reqs[t.Name] = ReqIndex{File: fileName, Index: len(temp)}
					temp = append(temp, req, body)
				} else {
					temp = append(temp, t)
				}
			}
			types = temp
		}
		code.Consts[fileName] = consts
		code.Enums[fileName] = enums
		code.Types[fileName] = types
	}
	for rpcIndex, rpc := range code.RPCs {
		if len(rpc.PathParams) == 0 {
			continue
		}
		params := make(map[string]string)
		for _, p := range rpc.PathParams {
			params[p] = ""
		}
		reqIndex := code.Reqs[rpc.Request]
		t := code.Types[reqIndex.File][reqIndex.Index]
		for _, f := range t.Fields {
			if f.Binding == nil || f.Binding.From != "path" {
				continue
			}
			if _, ok := params[f.Binding.Name]; !ok {
				err = errutil.Explain(nil, "path parameter %s not found in request type %s", f.Binding.Name, rpc.Request)
				return GoCode{}, err
			}
			params[f.Binding.Name] = f.Name
		}
		var paramNames []string
		for _, p := range rpc.PathParams {
			if params[p] == "" {
				err = errutil.Explain(nil, "path parameter %s not found in request type %s", p, rpc.Request)
				return GoCode{}, err
			}
			paramNames = append(paramNames, params[p])
		}
		rpc.PathParams = paramNames
		code.RPCs[rpcIndex] = rpc
	}
	return code, nil
}

// SplitRequestType splits a type into a whole type and a body type.
func SplitRequestType(t Type) (req Type, body Type) {
	req.Request = true
	req.Name = t.Name
	req.Comment = t.Comment

	body.Name = t.Name + "Body"
	body.RequestBody = true

	for _, field := range t.Fields {
		if field.Binding != nil {
			req.Fields = append(req.Fields, field)
		} else {
			body.Fields = append(body.Fields, field)
		}
	}
	return
}

// convertConsts converts IDL constants to Go constants
func convertConsts(code GoCode, doc httpidl.Document) ([]Const, error) {
	var ret []Const
	for _, c := range doc.Consts {
		t := httpidl.BaseType{Name: c.Type}
		typeName, err := getTypeName(code, t, nil, false)
		if err != nil {
			return nil, err
		}
		ret = append(ret, Const{
			Name:    c.Name,
			Type:    typeName,
			Value:   c.Value,
			Comment: formatComment(c.Comments),
		})
	}
	return ret, nil
}

// convertEnums converts IDL enums to Go enums
func convertEnums(code GoCode, doc httpidl.Document) ([]Enum, error) {
	var ret []Enum

	// Convert standard enums
	for _, e := range doc.Enums {
		var fields []EnumField
		for _, f := range e.Fields {
			fields = append(fields, EnumField{
				Name:    f.Name,
				Value:   f.Value,
				Comment: formatComment(f.Comments),
			})
		}
		ret = append(ret, Enum{
			Name:    e.Name,
			Fields:  fields,
			Comment: formatComment(e.Comments),
		})
	}

	// Convert enums from oneof types
	for _, t := range doc.Types {
		if !t.OneOf { // skip non-oneof types
			continue
		}
		var fields []EnumField
		for i, f := range t.Fields {
			fields = append(fields, EnumField{
				Name:  f.Name,
				Value: int64(i),
			})
		}
		ret = append(ret, Enum{
			Name:   t.Name + "Type",
			Fields: fields,
		})
	}
	return ret, nil
}

// convertTypes converts IDL struct types to Go struct types
func convertTypes(code GoCode, doc httpidl.Document) ([]Type, error) {
	var ret []Type
	for _, t := range doc.Types {
		// Skip generic types (they need instantiation via Redefined)
		if t.GenericName != nil {
			continue
		}
		var (
			typ Type
			err error
		)
		if t.Redefined != nil {
			typ, err = convertRedefinedType(code, t)
		} else {
			typ, err = convertType(code, t)
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, typ)
	}
	return ret, nil
}

// convertRedefinedType instantiates a redefined generic struct type
func convertRedefinedType(code GoCode, r httpidl.Type) (Type, error) {

	t, ok := httpidl.GetType(code.Files, r.Redefined.Name)
	if !ok {
		err := errutil.Explain(nil, "type %s not found", r.Redefined.Name)
		return Type{}, errutil.Explain(nil, "convert redefined type %s error: %w", r.Name, err)
	}

	var fields []httpidl.TypeField
	for _, f := range t.Fields {
		// Replace generic placeholders with concrete types
		f.Type = replaceGenericType(f.Type, *t.GenericName, r.Redefined.GenericType)
		fields = append(fields, f)
	}

	return convertType(code, httpidl.Type{
		Name:     r.Name,
		Fields:   fields,
		Position: r.Position,
		Comments: r.Comments,
	})
}

// replaceGenericType replaces a generic type in a field with a concrete type
func replaceGenericType(t httpidl.TypeDefinition, genericName string, genericType httpidl.TypeDefinition) httpidl.TypeDefinition {
	switch u := t.(type) {
	case httpidl.UserType:
		if u.Name == genericName {
			return genericType
		}
		return u
	case httpidl.ListType:
		u.Item = replaceGenericType(u.Item, genericName, genericType)
		return u
	case httpidl.MapType:
		u.Value = replaceGenericType(u.Value, genericName, genericType)
		return u
	default:
		return t
	}
}

// convertType converts an IDL struct type to a Go struct type
func convertType(code GoCode, t httpidl.Type) (Type, error) {
	r := Type{
		Name: t.Name,
	}

	// Handle oneof
	if t.OneOf {
		r.Fields = append(r.Fields, TypeField{
			Type:      "*" + r.Name + "TypeAsString",
			ValueType: r.Name + "TypeAsString",
			TypeKind:  []TypeKind{TypeKindPointer, TypeKindEnum},
			Name:      "FieldType",
			FieldTag:  "`json:\"field_type\"`",
		})
	}

	// Handle fields
	for _, f := range t.Fields {

		// Handle embedded types (flatten their fields into the struct)
		if embedType, ok := f.Type.(httpidl.EmbedType); ok {
			srcType, ok := httpidl.GetType(code.Files, embedType.Name)
			if !ok {
				return Type{}, errutil.Explain(nil, "embedded type %s not found for field in type %s", embedType.Name, r.Name)
			}
			retType, err := convertType(code, srcType)
			if err != nil {
				return Type{}, errutil.Explain(nil, "failed to convert embedded type %s in type %s: %w", embedType.Name, r.Name, err)
			}
			// Append embedded type's fields
			r.Fields = append(r.Fields, retType.Fields...)
			continue
		}

		// Convert field name to PascalCase for Go
		fieldName := httpidl.ToPascal(f.Name)

		// Determine Go type for the field
		typeName, err := getTypeName(code, f.Type, f.Annotations, true)
		if err != nil {
			return Type{}, errutil.Explain(nil, "get type name for field %s in type %s error: %w", f.Name, r.Name, err)
		}

		// Determine the category of the field (base, enum, struct, list, map)
		typeKind, valueType, err := getTypeKind(code, typeName)
		if err != nil {
			return Type{}, errutil.Explain(nil, "get type kind for field %s in type %s error: %w", f.Name, r.Name, err)
		}

		// Generate validation expressions for the field
		var validateExpr *string
		if f.Validate != nil {
			var s string
			s, err = genValidate(r.Name, fieldName, typeName, f.Validate, code.Funcs)
			if err != nil {
				return Type{}, errutil.Explain(nil, "generate validate for field %s in type %s error: %w", f.Name, r.Name, err)
			}
			validateExpr = &s
		}

		// Add the field to the struct
		field := TypeField{
			Type:      typeName,
			ValueType: valueType,
			TypeKind:  typeKind,
			Name:      fieldName,
			JSONTag:   f.JSONTag,
			FormTag:   f.FormTag,
			Binding:   f.Binding,
			Required:  f.Required,
			Validate:  validateExpr,
			Comment:   formatComment(f.Comments),
		}
		field.FieldTag = genFieldTag(field)
		r.Fields = append(r.Fields, field)
	}
	return r, nil
}

// getTypeName returns the Go type name for a given IDL type.
// It also respects the "go.type" annotation, which overrides the default type.
func getTypeName(code GoCode, t httpidl.TypeDefinition, arr []httpidl.Annotation, forceOptional bool) (string, error) {

	// Handle explicit "go.type" annotation
	if a, ok := httpidl.GetAnnotation(arr, "go.type"); ok {
		if a.Value == nil {
			return "", errutil.Explain(nil, `annotation "go.type" must have a value`)
		}
		s := strings.Trim(strings.TrimSpace(*a.Value), "\"")
		if s == "" {
			return "", errutil.Explain(nil, `annotation "go.type" must not be empty`)
		}
		return s, nil
	}

	switch typ := t.(type) {
	case httpidl.AnyType:
		return "", errutil.Explain(nil, `any type must have annotation "go.type"`)
	case httpidl.BaseType:
		var typeName string
		switch typ.Name {
		case "string":
			typeName = "string"
		case "int":
			typeName = "int64"
		case "float":
			typeName = "float64"
		case "bool":
			typeName = "bool"
		default:
			return "", errutil.Explain(nil, "unknown base type: %s", typ.Name)
		}
		if forceOptional {
			typeName = "*" + typeName
		}
		return typeName, nil
	case httpidl.UserType:
		typeName := typ.Name
		// Handle enum_as_string annotation
		if _, ok := httpidl.GetAnnotation(arr, "enum_as_string"); ok {
			if _, ok := httpidl.GetEnum(code.Files, typ.Name); !ok {
				return "", errutil.Explain(nil, "enum %s not found", typ.Name)
			}
			typeName += "AsString"
		}
		if forceOptional {
			typeName = "*" + typeName
		}
		return typeName, nil
	case httpidl.ListType:
		itemType, err := getTypeName(code, typ.Item, nil, false)
		if err != nil {
			return "", err
		}
		return "[]" + itemType, nil
	case httpidl.MapType:
		keyType := "string"
		if typ.Key == "int" {
			keyType = "int64"
		}
		valueType, err := getTypeName(code, typ.Value, nil, false)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("map[%s]%s", keyType, valueType), nil
	case httpidl.BinaryType:
		return "[]byte", nil // todo (lvan100) handle file
	default:
		return "", errutil.Explain(nil, "unknown type: %s", t.Text())
	}
}

// getTypeKind categorizes a Go type for code generation purposes.
func getTypeKind(code GoCode, typeName string) ([]TypeKind, string, error) {
	typeName, optional := strings.CutPrefix(typeName, "*")

	switch typeName {
	case "bool":
		if optional {
			return []TypeKind{TypeKindPointer, TypeKindBool}, typeName, nil
		}
		return []TypeKind{TypeKindBool}, typeName, nil
	case "int", "int8", "int16", "int32", "int64":
		if optional {
			return []TypeKind{TypeKindPointer, TypeKindInt}, typeName, nil
		}
		return []TypeKind{TypeKindInt}, typeName, nil
	case "uint", "uint8", "uint16", "uint32", "uint64":
		if optional {
			return []TypeKind{TypeKindPointer, TypeKindUint}, typeName, nil
		}
		return []TypeKind{TypeKindUint}, typeName, nil
	case "float32", "float64":
		if optional {
			return []TypeKind{TypeKindPointer, TypeKindFloat}, typeName, nil
		}
		return []TypeKind{TypeKindFloat}, typeName, nil
	case "string":
		if optional {
			return []TypeKind{TypeKindPointer, TypeKindString}, typeName, nil
		}
		return []TypeKind{TypeKindString}, typeName, nil
	default: // for linter
	}

	switch {
	case strings.HasPrefix(typeName, "[]"):
		if optional {
			return nil, "", errutil.Explain(nil, "list type can not be optional")
		}
		itemType, _, err := getTypeKind(code, typeName[2:])
		if err != nil {
			return nil, "", err
		}
		return append([]TypeKind{TypeKindList}, itemType...), typeName, nil
	case strings.HasPrefix(typeName, "map["):
		if optional {
			return nil, "", errutil.Explain(nil, "map type can not be optional")
		}
		return []TypeKind{TypeKindMap}, typeName, nil
	default:
		if _, ok := httpidl.GetEnum(code.Files, strings.TrimSuffix(typeName, "AsString")); ok {
			if optional {
				return []TypeKind{TypeKindPointer, TypeKindEnum}, typeName, nil
			}
			return []TypeKind{TypeKindEnum}, typeName, nil
		}
		if _, ok := httpidl.GetType(code.Files, typeName); ok {
			if optional {
				return []TypeKind{TypeKindPointer, TypeKindStruct}, typeName, nil
			}
			return []TypeKind{TypeKindStruct}, typeName, nil
		}
		return nil, "", errutil.Explain(nil, "unknown type: %s", typeName)
	}
}

// genFieldTag generates the struct tag for a Go struct field.
// It includes JSON tags and optional binding tags (path, query).
func genFieldTag(f TypeField) string {
	var tags []string

	// JSON tag
	{
		var sb strings.Builder
		sb.WriteString(`json:"`)
		sb.WriteString(f.JSONTag.Name)
		if f.JSONTag.OmitEmpty {
			sb.WriteString(",omitempty")
		}
		if f.JSONTag.OmitZero {
			sb.WriteString(",omitzero")
		}
		sb.WriteString(`"`)
		tags = append(tags, sb.String())
	}

	if f.Binding == nil {
		s := fmt.Sprintf(`form:"%s"`, f.FormTag.Name)
		tags = append(tags, s)
	} else {
		s := fmt.Sprintf(`%s:"%s"`, f.Binding.From, f.Binding.Name)
		tags = append(tags, s)
	}

	return "`" + strings.Join(tags, " ") + "`"
}

// ValidateFunc represents a custom validation function
type ValidateFunc struct {
	Name      string
	FieldType string
}

// builtinFuncs is a set of built-in validation functions
var builtinFuncs = map[string]struct{}{
	"len": {},
}

// genValidate generates validation code for a struct field based on its "validate" annotation.
// Returns a Go code snippet that checks the field and returns an error if validation fails.
func genValidate(receiverType, fieldName, fieldType string, expr validate.Expr, funcs map[string]ValidateFunc) (string, error) {

	optional := strings.HasPrefix(fieldType, "*")
	dollar := "x." + fieldName
	if optional {
		dollar = "*" + dollar
	}

	// Generate the Go expression for validation
	str, err := genValidateExpr(dollar, fieldType, expr, funcs)
	if err != nil {
		return "", errutil.Explain(err, `failed to generate validate expression for %s.%s`, receiverType, fieldName)
	}

	// Wrap in an if statement returning an error on failure
	str = fmt.Sprintf(`if !(%s) {
		return errutil.Explain(nil,"validate failed on %s.%s")
	}`, str, receiverType, fieldName)

	if optional {
		str = fmt.Sprintf(`if x.%s != nil { %s }`, fieldName, str)
	}
	return str, nil
}

// genValidateExpr recursively generates Go code for a validation expression
func genValidateExpr(fieldName, fieldType string, expr validate.Expr, funcs map[string]ValidateFunc) (string, error) {
	switch x := expr.(type) {
	case validate.BinaryExpr:
		if x.Left == nil {
			return "", nil
		}
		left, err := genValidateExpr(fieldName, fieldType, x.Left, funcs)
		if err != nil {
			return "", err
		}
		if x.Right == nil {
			return left, nil
		}
		right, err := genValidateExpr(fieldName, fieldType, x.Right, funcs)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s %s %s", left, x.Op, right), nil

	case validate.UnaryExpr:
		str, err := genValidateExpr(fieldName, fieldType, x.Expr, funcs)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s%s", x.Op, str), nil

	case *validate.FuncCall:
		if len(x.Args) == 0 {
			return x.Name + "()", nil
		}

		// Register or validate custom functions
		if _, ok := builtinFuncs[x.Name]; !ok {
			if len(x.Args) != 1 {
				return "", errutil.Explain(nil, "func %s only accepts 1 argument of type %s", x.Name, fieldType)
			}
			if !strings.HasPrefix(x.Name, "OneOf") {
				if f, ok := funcs[x.Name]; ok {
					if f.FieldType != fieldType {
						return "", errutil.Explain(nil, "func %s only accepts type %s", x.Name, f.FieldType)
					}
				} else {
					funcs[x.Name] = ValidateFunc{
						Name:      x.Name,
						FieldType: strings.TrimPrefix(fieldType, "*"),
					}
				}
			}
		} else {
			// Validate built-in functions
			switch x.Name {
			case "len":
				if len(x.Args) != 1 {
					return "", errutil.Explain(nil, "func len only accepts 1 argument")
				}
			default: // for linter
			}
		}

		var args []string
		for _, arg := range x.Args {
			str, err := genValidateExpr(fieldName, fieldType, arg, funcs)
			if err != nil {
				return "", err
			}
			args = append(args, str)
		}
		return fmt.Sprintf("%s(%s)", x.Name, strings.Join(args, ", ")), nil

	case *validate.InnerExpr:
		str, err := genValidateExpr(fieldName, fieldType, x.Expr, funcs)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s)", str), nil

	case validate.PrimaryExpr:
		if x.Inner != nil {
			return genValidateExpr(fieldName, fieldType, x.Inner, funcs)
		}
		if x.Call != nil {
			return genValidateExpr(fieldName, fieldType, x.Call, funcs)
		}
		if x.Value == "$" {
			return fieldName, nil
		}
		return x.Value, nil

	default:
		return "", errutil.Explain(nil, "unknown expression type: %s", x.Text())
	}
}

// convertRPC converts a TIDL RPC to a RPC.
func convertRPC(r httpidl.RPC) (RPC, error) {

	// Parse the path (source)
	var pathParams []string
	segments, err := pathidl.Parse(r.Path)
	if err != nil {
		return RPC{}, errutil.Explain(err, `failed to parse path %s`, r.Path)
	}
	var formatPath strings.Builder
	for _, seg := range segments {
		formatPath.WriteString("/")
		if seg.Type == pathidl.Static {
			formatPath.WriteString(seg.Value)
			continue
		}
		formatPath.WriteString("%s")
		pathParams = append(pathParams, seg.Value)
	}

	return RPC{
		Name:        r.Name,
		Request:     r.Request.Name,
		Response:    r.Response.UserType.Name,
		Stream:      r.Response.Stream,
		Path:        r.Path,
		FormatPath:  formatPath.String(),
		PathParams:  pathParams,
		Method:      r.Method,
		ContentType: r.ContentType,
		Comment:     formatComment(r.Comments),
	}, nil
}

// formatComment converts a tidl.Comments into Go comments.
func formatComment(c httpidl.Comments) string {
	var lines []string
	for _, s := range c.Above {
		lines = append(lines, s.Text...)
	}
	if c.Right != nil {
		lines = append(lines, c.Right.Text...)
	}
	return strings.Join(lines, "\n")
}
