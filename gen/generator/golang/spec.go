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
	TypeKindBool    = TypeKind(0)
	TypeKindInt     = TypeKind(1)
	TypeKindUint    = TypeKind(2)
	TypeKindFloat   = TypeKind(3)
	TypeKindString  = TypeKind(4)
	TypeKindStruct  = TypeKind(5)
	TypeKindEnum    = TypeKind(6)
	TypeKindList    = TypeKind(7)
	TypeKindMap     = TypeKind(8)
	TypeKindPointer = TypeKind(9)
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
	Name    string
	Fields  []TypeField
	Comment string

	Request     bool
	OnRequest   bool
	RequestBody bool
}

// TypeField represents a field in a Go struct
type TypeField struct {
	Name     string
	Type     string // for field
	Required bool
	Comment  string

	TypeKind  []TypeKind
	ValueType string // for getter/setter

	JSONTag httpidl.JSONTag
	FormTag httpidl.FormTag
	Binding *httpidl.Binding

	FieldTag       string
	ValidateExpr   *string
	ValidateNested *string
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

// RPC represents a single remote procedure call with HTTP metadata.
type RPC struct {
	Name     string // Method name
	Request  string // Request type name
	Response string // Response type name
	Stream   bool   // Whether this RPC is a streaming RPC
	Comment  string // Comment of the RPC

	Path        string // HTTP path
	Method      string // HTTP method (GET, POST, etc.)
	ContentType string // HTTP Content-Type

	FormatPath   string            // Formatted HTTP path
	PathParams   map[string]string // HTTP path parameters
	PathSegments []pathidl.Segment // HTTP path segments

	ConnTimeout  int // Connection timeout, ms
	ReadTimeout  int // Read timeout, ms
	WriteTimeout int // Write timeout, ms
}

type ReqIndex struct {
	File  string
	Index int
}

type GoSpec struct {
	Meta   *httpidl.MetaInfo
	Files  map[string]httpidl.Document
	Consts map[string][]Const
	Enums  map[string][]Enum
	Types  map[string][]Type
	RPCs   []RPC
	Funcs  map[string]ValidateFunc // Collected validation functions
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

// Convert converts an IDL project to Go code.
func Convert(dir string) (GoSpec, error) {
	project, err := httpidl.ParseDir(dir)
	if err != nil {
		return GoSpec{}, err
	}

	spec := GoSpec{
		Meta:   project.Meta,
		Files:  project.Files,
		Consts: make(map[string][]Const),
		Enums:  make(map[string][]Enum),
		Types:  make(map[string][]Type),
		Funcs:  make(map[string]ValidateFunc),
	}

	// Collect all RPC definitions
	for _, doc := range project.Files {
		for _, r := range doc.RPCs {
			rpc := RPC{
				Name:         r.Name,
				Request:      r.Request.Name,
				Response:     r.Response.Name,
				Stream:       r.Stream,
				Comment:      formatComment(r.Comments),
				Path:         r.Path,
				FormatPath:   r.Path, // 假设是普通路径
				PathParams:   r.PathParams,
				PathSegments: r.PathSegments,
				Method:       r.Method,
				ContentType:  r.ContentType,
				ConnTimeout:  r.ConnTimeout,
				ReadTimeout:  r.ReadTimeout,
				WriteTimeout: r.WriteTimeout,
			}
			spec.RPCs = append(spec.RPCs, rpc)
		}
	}
	sort.Slice(spec.RPCs, func(i, j int) bool {
		return spec.RPCs[i].Name < spec.RPCs[j].Name
	})

	for fileName, doc := range project.Files {
		consts, err := convertConsts(spec, doc)
		if err != nil {
			return GoSpec{}, errutil.Explain(nil, "convert consts error: %w", err)
		}
		enums, err := convertEnums(spec, doc)
		if err != nil {
			return GoSpec{}, errutil.Explain(nil, "convert enums error: %w", err)
		}
		types, err := convertTypes(spec, doc)
		if err != nil {
			return GoSpec{}, errutil.Explain(nil, "convert types error: %w", err)
		}
		{
			var temp []Type
			for _, t := range types {
				if t.Request {
					req, body := SplitRequestType(t)
					temp = append(temp, req, body)
				} else {
					temp = append(temp, t)
				}
			}
			types = temp
		}
		spec.Consts[fileName] = consts
		spec.Enums[fileName] = enums
		spec.Types[fileName] = types
	}

	for i, rpc := range spec.RPCs {
		for k, s := range rpc.PathParams {
			rpc.PathParams[k] = httpidl.ToPascal(s)
		}
		var formatPath strings.Builder
		for _, seg := range rpc.PathSegments {
			formatPath.WriteString("/")
			if seg.Type == pathidl.Static {
				formatPath.WriteString(seg.Value)
				continue
			}
			formatPath.WriteString("%s")
		}
		rpc.FormatPath = formatPath.String()
		spec.RPCs[i] = rpc
	}
	return spec, nil
}

// SplitRequestType splits a type into a whole type and a body type.
func SplitRequestType(t Type) (req Type, body Type) {
	req.Request = true
	req.OnRequest = true
	req.Name = t.Name
	req.Comment = t.Comment

	body.OnRequest = true
	body.RequestBody = true
	body.Name = t.Name + "Body"

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
func convertConsts(spec GoSpec, doc httpidl.Document) ([]Const, error) {
	var ret []Const
	for _, c := range doc.Consts {
		typeName, err := getBaseTypeName(c.Type.Name)
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
func convertEnums(spec GoSpec, doc httpidl.Document) ([]Enum, error) {
	var ret []Enum
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
	return ret, nil
}

// convertTypes converts IDL struct types to Go struct types
func convertTypes(spec GoSpec, doc httpidl.Document) ([]Type, error) {
	var ret []Type
	for _, t := range doc.Types {
		// Skip generic types (they need instantiation)
		if t.GenericParam != nil {
			continue
		}
		typ, err := convertType(spec, t)
		if err != nil {
			return nil, err
		}
		ret = append(ret, typ)
	}
	return ret, nil
}

// convertType converts an IDL struct type to a Go struct type
func convertType(spec GoSpec, t httpidl.Type) (Type, error) {
	r := Type{
		Name:      t.Name,
		Request:   t.Request,
		OnRequest: t.OnRequest,
		//Comment: formatComment(t.Comments),
	}
	for _, f := range t.Fields {
		fieldName := httpidl.ToPascal(f.Name)

		// Get the type name
		typeName, err := getTypeName(spec, f)
		if err != nil {
			return Type{}, errutil.Explain(nil, "get type name for field %s in type %s error: %w", f.Name, r.Name, err)
		}

		// Determine the category of the field (base, enum, struct, list, map)
		typeKind, valueType, err := getTypeKind(spec, typeName)
		if err != nil {
			return Type{}, errutil.Explain(nil, "get type kind for field %s in type %s error: %w", f.Name, r.Name, err)
		}

		// Generate validation expressions for the field
		var validateExpr *string
		if t.OnRequest && f.ValidateExpr != nil {
			var s string
			s, err = genValidateExpr(r.Name, fieldName, typeName, f.ValidateExpr, spec.Funcs)
			if err != nil {
				return Type{}, errutil.Explain(nil, "generate validate for field %s in type %s error: %w", f.Name, r.Name, err)
			}
			if s != "" {
				validateExpr = &s
			}
		}

		// Generate validation expressions for nested fields
		var ValidateNested *string
		if t.OnRequest && f.ValidateNested {
			s := genValidateNested(r.Name, fieldName, "x."+fieldName, typeKind, 0)
			if s != "" {
				ValidateNested = &s
			}
		}

		// Add the field to the struct
		field := TypeField{
			Name:           fieldName,
			Type:           typeName,
			Required:       f.Required,
			Comment:        formatComment(f.Comments),
			TypeKind:       typeKind,
			ValueType:      valueType,
			JSONTag:        f.JSONTag,
			FormTag:        f.FormTag,
			Binding:        f.Binding,
			ValidateExpr:   validateExpr,
			ValidateNested: ValidateNested,
		}
		field.FieldTag = genFieldTag(field)
		r.Fields = append(r.Fields, field)
	}
	return r, nil
}

// getTypeName returns the Go type name for a given IDL type
func getTypeName(spec GoSpec, f httpidl.TypeField) (string, error) {
	if a, ok := httpidl.GetAnnotation(f.Annotations, "go.type"); ok {
		if a.Value == nil {
			return "", errutil.Explain(nil, `annotation "go.type" must have a value`)
		}
		s := strings.Trim(strings.TrimSpace(*a.Value), "\"")
		if s == "" {
			return "", errutil.Explain(nil, `annotation "go.type" must not be empty`)
		}
		return s, nil
	}

	switch typ := f.Type.(type) {
	case httpidl.AnyType:
		return "", errutil.Explain(nil, `any type must have annotation "go.type"`)
	case httpidl.BinaryType:
		return "[]byte", nil // todo (lvan100) file
	case httpidl.BaseType:
		s, err := getBaseTypeName(typ.Name)
		if err != nil {
			return "", err
		}
		return "*" + s, nil
	case httpidl.UserType:
		typeName := typ.Name
		if f.EnumAsString {
			typeName += "AsString"
		}
		return "*" + typeName, nil
	default:
		s, err := getTypeName0(spec, typ)
		if err != nil {
			return "", err
		}
		return s, nil
	}
}

// getBaseTypeName returns the Go type name for a given IDL base type.
func getBaseTypeName(typeName string) (string, error) {
	switch typeName {
	case "string":
		return "string", nil
	case "int":
		return "int64", nil
	case "float":
		return "float64", nil
	case "bool":
		return "bool", nil
	default:
		return "", errutil.Explain(nil, "unknown base type: %s", typeName)
	}
}

// getTypeName0 returns the Go type name for a given IDL type.
func getTypeName0(spec GoSpec, t httpidl.TypeDefinition) (string, error) {
	switch typ := t.(type) {
	case httpidl.BaseType:
		return getBaseTypeName(typ.Name)
	case httpidl.UserType:
		if _, ok := httpidl.GetEnum(spec.Files, typ.Name); ok {
			return typ.Name, nil
		}
		return "*" + typ.Name, nil
	case httpidl.ListType:
		itemType, err := getTypeName0(spec, typ.Item)
		if err != nil {
			return "", err
		}
		return "[]" + itemType, nil
	case httpidl.MapType:
		keyType := "string"
		if typ.Key == "int" {
			keyType = "int64"
		}
		valueType, err := getTypeName0(spec, typ.Value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("map[%s]%s", keyType, valueType), nil
	default:
		return "", errutil.Explain(nil, "unknown type: %s", t.Text())
	}
}

// getTypeKind categorizes a Go type for code generation purposes.
func getTypeKind(spec GoSpec, typeName string) ([]TypeKind, string, error) {
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
		itemType, _, err := getTypeKind(spec, typeName[2:])
		if err != nil {
			return nil, "", err
		}
		return append([]TypeKind{TypeKindList}, itemType...), typeName, nil
	case strings.HasPrefix(typeName, "map["):
		if optional {
			return nil, "", errutil.Explain(nil, "map type can not be optional")
		}
		itemInex := strings.Index(typeName, "]")
		itemType, _, err := getTypeKind(spec, typeName[itemInex+1:])
		if err != nil {
			return nil, "", err
		}
		return append([]TypeKind{TypeKindMap}, itemType...), typeName, nil
	default:
		if _, ok := httpidl.GetEnum(spec.Files, strings.TrimSuffix(typeName, "AsString")); ok {
			if optional {
				return []TypeKind{TypeKindPointer, TypeKindEnum}, typeName, nil
			}
			return []TypeKind{TypeKindEnum}, typeName, nil
		}
		if _, ok := httpidl.GetType(spec.Files, typeName); ok {
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

	if f.Required {
		tags = append(tags, `validate:"required"`)
	}

	return "`" + strings.Join(tags, " ") + "`"
}

// genValidateNested generates the nested validation code for a Go struct field.
func genValidateNested(receiverType, fieldName string, itemName string, typeKind []TypeKind, depth int) string {
	childName := fmt.Sprintf("v%d", depth)
	switch typeKind[0] {
	case TypeKindList, TypeKindMap:
		str := genValidateNested(receiverType, fieldName, childName, typeKind[1:], depth+1)
		if str == "" {
			return ""
		}
		str = fmt.Sprintf(`for _, %s := range %s {
				%s
			}`, childName, itemName, str)
		return str
	case TypeKindPointer:
		if typeKind[1] != TypeKindStruct {
			return ""
		}
		str := fmt.Sprintf(`if %s != nil {
				if validateErr := %s.Validate(); validateErr != nil {
					err = errutil.Stack(err, "validate failed on \"%s.%s\": %%w", validateErr)
				}
			}`, itemName, itemName, receiverType, fieldName)
		return str
	default:
		return ""
	}
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

// genValidateExpr generates the Go code for a validation expression
func genValidateExpr(receiverType, fieldName, fieldType string,
	expr validate.Expr, funcs map[string]ValidateFunc) (string, error) {

	// 对于结构体而言，只应当验证字段非空，其内部字段的验证应当由自己完成
	fieldType, optional := strings.CutPrefix(fieldType, "*")
	dollar := "x." + fieldName
	if optional {
		dollar = "*" + dollar
	}

	// Generate the Go expression for validation
	str, err := compileValidateExpr(dollar, fieldType, expr, funcs)
	if err != nil {
		return "", errutil.Explain(err, `failed to generate validate expression for %s.%s`, receiverType, fieldName)
	}

	// Wrap in an if statement returning an error on failure
	str = fmt.Sprintf(`if !(%s) {
		err = errutil.Stack(err,"validate failed on \"%s.%s\"")
	}`, str, receiverType, fieldName)

	if optional {
		str = fmt.Sprintf(`if x.%s != nil { %s }`, fieldName, str)
	}
	return str, nil
}

// compileValidateExpr recursively generates Go code for a validation expression
func compileValidateExpr(fieldName, fieldType string, expr validate.Expr, funcs map[string]ValidateFunc) (string, error) {
	switch x := expr.(type) {
	case validate.BinaryExpr:
		if x.Left == nil {
			return "", nil
		}
		left, err := compileValidateExpr(fieldName, fieldType, x.Left, funcs)
		if err != nil {
			return "", err
		}
		if x.Right == nil {
			return left, nil
		}
		right, err := compileValidateExpr(fieldName, fieldType, x.Right, funcs)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s %s %s", left, x.Op, right), nil

	case validate.UnaryExpr:
		str, err := compileValidateExpr(fieldName, fieldType, x.Expr, funcs)
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
						FieldType: fieldType,
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
			str, err := compileValidateExpr(fieldName, fieldType, arg, funcs)
			if err != nil {
				return "", err
			}
			args = append(args, str)
		}
		return fmt.Sprintf("%s(%s)", x.Name, strings.Join(args, ", ")), nil

	case *validate.InnerExpr:
		str, err := compileValidateExpr(fieldName, fieldType, x.Expr, funcs)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s)", str), nil

	case validate.PrimaryExpr:
		if x.Inner != nil {
			return compileValidateExpr(fieldName, fieldType, x.Inner, funcs)
		}
		if x.Call != nil {
			return compileValidateExpr(fieldName, fieldType, x.Call, funcs)
		}
		if x.Value == "$" {
			return fieldName, nil
		}
		return x.Value, nil

	default:
		return "", errutil.Explain(nil, "unknown expression type: %s", x.Text())
	}
}
