package golang

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/go-spring/gs-http-gen/gen/generator"
	"github.com/go-spring/gs-http-gen/lib/tidl"
	"github.com/go-spring/gs-http-gen/lib/vidl"
	"github.com/lvan100/errutil"
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

// TypeKind represents kind of a Go field type
type TypeKind int

const (
	TypeKindUnknown TypeKind = iota
	TypeKindBaseType
	TypeKindOptionalBaseType
	TypeKindEnumType
	TypeKindOptionalEnumType
	TypeKindStructType
	TypeKindOptionalStructType
	TypeKindListType
	TypeKindMapType
)

// Type represents a Go struct
type Type struct {
	Name    string
	Fields  []TypeField
	Comment string
	Split   bool
}

// TypeField represents a field in a Go struct
type TypeField struct {
	Type     string
	TypeKind TypeKind
	Name     string
	Tag      string
	Validate *string
	Binding  *Binding
	Comment  string
}

// Binding represents a field binding from path, or query
type Binding struct {
	From string // Source: path/query
	Name string // Field name in the source
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
	Name     string // Method name
	Request  string // Request type name
	Response string // Response type name
	Stream   bool   // Whether this RPC is a streaming RPC
	Path     string // HTTP path
	Method   string // HTTP method (GET, POST, etc.)
	Comment  string // Comment of the RPC
}

type Go struct {
	Files  map[string]tidl.Document
	Meta   *tidl.MetaInfo
	Reqs   map[string]struct{} // request type name
	Consts map[string][]Const
	Enums  map[string][]Enum
	Types  map[string][]Type
	RPCs   []RPC
	Funcs  map[string]ValidateFunc // Collected validation functions
}

func Convert(dir string) (Go, error) {
	files, meta, err := generator.ParseDir(dir)
	if err != nil {
		return Go{}, err
	}

	g := Go{
		Files:  files,
		Meta:   meta,
		Reqs:   make(map[string]struct{}),
		Consts: make(map[string][]Const),
		Enums:  make(map[string][]Enum),
		Types:  make(map[string][]Type),
		Funcs:  make(map[string]ValidateFunc),
	}

	// Collect all RPC definitions
	for _, doc := range files {
		for _, r := range doc.RPCs {
			rpc, err := convertRPC(r)
			if err != nil {
				return g, err
			}
			g.RPCs = append(g.RPCs, rpc)
			g.Reqs[rpc.Request] = struct{}{}
		}
	}
	sort.Slice(g.RPCs, func(i, j int) bool {
		return g.RPCs[i].Name < g.RPCs[j].Name
	})

	for fileName, doc := range files {
		consts, err := convertConsts(g, doc)
		if err != nil {
			return g, errutil.Explain(nil, "convert consts error: %w", err)
		}
		enums, err := convertEnums(g, doc)
		if err != nil {
			return g, errutil.Explain(nil, "convert enums error: %w", err)
		}
		types, err := convertTypes(g, doc)
		if err != nil {
			return g, errutil.Explain(nil, "convert types error: %w", err)
		}
		{
			var temp []Type
			for _, t := range types {
				if _, ok := g.Reqs[t.Name]; ok {
					whole, body := SplitType(t)
					temp = append(temp, whole, body)
				} else {
					temp = append(temp, t)
				}
			}
			types = temp
		}
		g.Consts[fileName] = consts
		g.Enums[fileName] = enums
		g.Types[fileName] = types
	}
	return g, nil
}

// SplitType splits a type into a whole type and a body type.
func SplitType(t Type) (whole Type, body Type) {
	whole.Split = true
	whole.Name = t.Name
	whole.Comment = t.Comment
	body.Name = t.Name + "Body"
	for _, field := range t.Fields {
		if field.Binding != nil {
			whole.Fields = append(whole.Fields, field)
		} else {
			body.Fields = append(body.Fields, field)
		}
	}
	return
}

// convertConsts converts IDL constants to Go constants
func convertConsts(g Go, doc tidl.Document) ([]Const, error) {
	var ret []Const
	for _, c := range doc.Consts {
		t := tidl.BaseType{Name: c.Type}
		typeName, err := getTypeName(g, t, nil)
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
func convertEnums(g Go, doc tidl.Document) ([]Enum, error) {
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
			fieldName, err := getJSONName(f.Name, f.Annotations)
			if err != nil {
				return nil, errutil.Explain(nil, "get json name for type %s field %s error: %w", t.Name, f.Name, err)
			}
			fields = append(fields, EnumField{
				Name:  fieldName,
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

// getJSONName returns the JSON name for a struct field.
func getJSONName(fieldName string, arr []tidl.Annotation) (string, error) {
	if a, ok := tidl.GetAnnotation(arr, "json"); ok {
		if a.Value == nil {
			return "", errutil.Explain(nil, `annotation "json" value is nil`)
		}
		s := strings.TrimSpace(*a.Value)
		if s == "" {
			return "", errutil.Explain(nil, `annotation "json" value is empty`)
		}
		s = strings.Trim(s, "\"") // remove quotes
		s = strings.TrimSpace(strings.SplitN(s, ",", 2)[0])
		if s != "" {
			return s, nil
		}
	}
	return fieldName, nil
}

// convertTypes converts IDL struct types to Go struct types
func convertTypes(g Go, doc tidl.Document) ([]Type, error) {
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
			typ, err = convertRedefinedType(g, t)
		} else {
			typ, err = convertType(g, t)
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, typ)
	}
	return ret, nil
}

// convertRedefinedType instantiates a redefined generic struct type
func convertRedefinedType(g Go, r tidl.Type) (Type, error) {

	t, ok := tidl.GetType(g.Files, r.Redefined.Name)
	if !ok {
		err := errutil.Explain(nil, "type %s not found", r.Redefined.Name)
		return Type{}, errutil.Explain(nil, "convert redefined type %s error: %w", r.Name, err)
	}

	var fields []tidl.TypeField
	for _, f := range t.Fields {
		// Replace generic placeholders with concrete types
		f.FieldType = replaceGenericType(f.FieldType, *t.GenericName, r.Redefined.GenericType)
		fields = append(fields, f)
	}

	return convertType(g, tidl.Type{
		Name:     r.Name,
		Fields:   fields,
		Position: r.Position,
		Comments: r.Comments,
	})
}

// replaceGenericType replaces a generic type in a field with a concrete type
func replaceGenericType(t tidl.TypeDefinition, genericName string, genericType tidl.TypeDefinition) tidl.TypeDefinition {
	switch u := t.(type) {
	case tidl.UserType:
		if u.Name == genericName {
			return genericType
		}
		return u
	case tidl.ListType:
		u.Item = replaceGenericType(u.Item, genericName, genericType)
		return u
	case tidl.MapType:
		u.Value = replaceGenericType(u.Value, genericName, genericType)
		return u
	default:
		return t
	}
}

// convertType converts an IDL struct type to a Go struct type
func convertType(g Go, t tidl.Type) (Type, error) {
	r := Type{
		Name: t.Name,
	}

	// Handle oneof
	if t.OneOf {
		r.Fields = append(r.Fields, TypeField{
			Type:     r.Name + "TypeAsString",
			TypeKind: TypeKindEnumType,
			Name:     "FieldType",
			Tag:      "`json:\"field_type\"`",
		})
	}

	// Handle fields
	for _, f := range t.Fields {

		// Handle embedded types (flatten their fields into the struct)
		if embedType, ok := f.FieldType.(tidl.EmbedType); ok {
			srcType, ok := tidl.GetType(g.Files, embedType.Name)
			if !ok {
				return Type{}, errutil.Explain(nil, "embedded type %s not found for field in type %s", embedType.Name, r.Name)
			}
			retType, err := convertType(g, srcType)
			if err != nil {
				return Type{}, errutil.Explain(nil, "failed to convert embedded type %s in type %s: %w", embedType.Name, r.Name, err)
			}
			// Append embedded type's fields
			r.Fields = append(r.Fields, retType.Fields...)
			continue
		}

		// Convert field name to PascalCase for Go
		fieldName := tidl.ToPascal(f.Name)

		// Determine Go type for the field
		typeName, err := getTypeName(g, f.FieldType, f.Annotations)
		if err != nil {
			return Type{}, errutil.Explain(nil, "get type name for field %s in type %s error: %w", f.Name, r.Name, err)
		}

		// Determine the category of the field (base, enum, struct, list, map)
		typeKind, err := getTypeKind(g, typeName)
		if err != nil {
			return Type{}, errutil.Explain(nil, "get type kind for field %s in type %s error: %w", f.Name, r.Name, err)
		}

		// Parse HTTP binding info from annotations (path, query)
		binding, err := parseBinding(f.Annotations)
		if err != nil {
			return Type{}, errutil.Explain(nil, "parse binding for field %s in type %s error: %w", f.Name, r.Name, err)
		}

		// Generate struct tag for JSON, query/path bindings
		fieldTag, err := genFieldTag(f.Name, typeName, f.Annotations, binding)
		if err != nil {
			return Type{}, errutil.Explain(nil, "generate field tag for field %s in type %s error: %w", f.Name, r.Name, err)
		}

		// Generate validation expressions for the field
		validate, err := genValidate(r.Name, fieldName, typeName, f.Annotations, g.Funcs)
		if err != nil {
			return Type{}, errutil.Explain(nil, "generate validate for field %s in type %s error: %w", f.Name, r.Name, err)
		}

		// Add the field to the struct
		r.Fields = append(r.Fields, TypeField{
			Type:     typeName,
			TypeKind: typeKind,
			Name:     fieldName,
			Tag:      fieldTag,
			Validate: validate,
			Binding:  binding,
			Comment:  formatComment(f.Comments),
		})
	}
	return r, nil
}

// getTypeName returns the Go type name for a given IDL type.
// It also respects the "go.type" annotation, which overrides the default type.
func getTypeName(g Go, t tidl.TypeDefinition, arr []tidl.Annotation) (string, error) {

	// Handle explicit "go.type" annotation
	if a, ok := tidl.GetAnnotation(arr, "go.type"); ok {
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
	case tidl.AnyType:
		return "", errutil.Explain(nil, `any type must have annotation "go.type"`)
	case tidl.BaseType:
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
		if typ.Optional {
			typeName = "*" + typeName
		}
		return typeName, nil
	case tidl.UserType:
		typeName := typ.Name
		// Handle enum_as_string annotation
		if _, ok := tidl.GetAnnotation(arr, "enum_as_string"); ok {
			if _, ok := tidl.GetEnum(g.Files, typ.Name); !ok {
				return "", errutil.Explain(nil, "enum %s not found", typ.Name)
			}
			typeName += "AsString"
		}
		if typ.Optional {
			typeName = "*" + typeName
		}
		return typeName, nil
	case tidl.ListType:
		itemType, err := getTypeName(g, typ.Item, nil)
		if err != nil {
			return "", err
		}
		return "[]" + itemType, nil
	case tidl.MapType:
		keyType := "string"
		if typ.Key == "int" {
			keyType = "int64"
		}
		valueType, err := getTypeName(g, typ.Value, nil)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("map[%s]%s", keyType, valueType), nil
	case tidl.BinaryType:
		return "[]byte", nil // todo (lvan100) handle file
	default:
		return "", errutil.Explain(nil, "unknown type: %s", t.Text())
	}
}

// getTypeKind categorizes a Go type for code generation purposes.
func getTypeKind(g Go, typeName string) (TypeKind, error) {
	typeName, optional := strings.CutPrefix(typeName, "*")

	switch typeName {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "string", "bool":
		if optional {
			return TypeKindOptionalBaseType, nil
		}
		return TypeKindBaseType, nil
	default: // for linter
	}

	switch {
	case strings.HasPrefix(typeName, "[]"):
		if optional {
			return TypeKindUnknown, errutil.Explain(nil, "list type can not be optional")
		}
		return TypeKindListType, nil
	case strings.HasPrefix(typeName, "map["):
		if optional {
			return TypeKindUnknown, errutil.Explain(nil, "map type can not be optional")
		}
		return TypeKindMapType, nil
	default:
		if _, ok := tidl.GetEnum(g.Files, strings.TrimSuffix(typeName, "AsString")); ok {
			if optional {
				return TypeKindOptionalEnumType, nil
			}
			return TypeKindEnumType, nil
		}
		if _, ok := tidl.GetType(g.Files, typeName); ok {
			if optional {
				return TypeKindOptionalStructType, nil
			}
			return TypeKindStructType, nil
		}
		return TypeKindUnknown, errutil.Explain(nil, "unknown type: %s", typeName)
	}
}

// parseBinding parses a field's HTTP binding information from annotations.
// Supported sources: path, query.
func parseBinding(arr []tidl.Annotation) (*Binding, error) {
	a, ok := tidl.GetAnnotation(arr, "path", "query")
	if !ok {
		return nil, nil
	}
	if a.Value == nil {
		return nil, errutil.Explain(nil, "annotation %q value is nil", a.Key)
	}
	val := strings.TrimSpace(strings.Trim(*a.Value, "\""))
	if val == "" {
		return nil, errutil.Explain(nil, "annotation %q value is empty", a.Key)
	}
	return &Binding{From: a.Key, Name: val}, nil
}

// genFieldTag generates the struct tag for a Go struct field.
// It includes JSON tags and optional binding tags (path, query).
func genFieldTag(fieldName, typeName string, arr []tidl.Annotation, binding *Binding) (string, error) {
	var tags []string

	var jsonName string
	var omitZero bool
	omitEmpty := strings.HasPrefix(typeName, "*")

	// Parse "json" annotation
	if a, ok := tidl.GetAnnotation(arr, "json"); ok {
		if a.Value == nil {
			return "", errutil.Explain(nil, `annotation "json" value is nil`)
		}
		s := strings.TrimSpace(*a.Value)
		if s == "" {
			return "", errutil.Explain(nil, `annotation "json" value is empty`)
		}
		s = strings.Trim(s, "\"") // Remove quotes
		for i, v := range strings.Split(s, ",") {
			v = strings.TrimSpace(v)
			if i == 0 {
				if v != "" {
					jsonName = v
				}
				continue
			}
			switch v {
			case "omitempty":
				omitEmpty = true
			case "non-omitempty":
				omitEmpty = false
			case "omitzero":
				omitZero = true
			default: // for linter
			}
		}
	}

	if jsonName == "" {
		jsonName += fieldName
	}
	if omitEmpty {
		jsonName += ",omitempty"
	}
	if omitZero {
		jsonName += ",omitzero"
	}
	tags = append(tags, fmt.Sprintf("json:\"%s\"", jsonName))

	// Generate binding tag
	if binding != nil {
		tags = append(tags, fmt.Sprintf("%s:\"%s\"", binding.From, binding.Name))
	}

	return "`" + strings.Join(tags, " ") + "`", nil
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
func genValidate(receiverType, fieldName, fieldType string, arr []tidl.Annotation, funcs map[string]ValidateFunc) (*string, error) {
	a, ok := tidl.GetAnnotation(arr, "validate")
	if !ok {
		return nil, nil
	}
	if a.Value == nil {
		return nil, errutil.Explain(nil, `annotation "validate" value is nil`)
	}

	// Unquote the validation expression string
	strValue, err := strconv.Unquote(*a.Value)
	if err != nil {
		return nil, errutil.Explain(nil, `annotation "validate" value is not properly quoted`)
	}
	if strValue == "" {
		return nil, errutil.Explain(nil, `annotation "validate" value is empty`)
	}

	// Parse the validation expression
	expr, err := vidl.Parse(strValue)
	if err != nil {
		return nil, errutil.Explain(nil, `failed to parse validate expression %s: %w`, strValue, err)
	}

	optional := strings.HasPrefix(fieldType, "*")
	dollar := "x." + fieldName
	if optional {
		dollar = "*" + dollar
	}

	// Generate the Go expression for validation
	str, err := genValidateExpr(dollar, fieldType, expr, funcs)
	if err != nil {
		return nil, errutil.Explain(nil, `failed to generate validate expression for %s: %w`, strValue, err)
	}

	// Wrap in an if statement returning an error on failure
	str = fmt.Sprintf(`if !(%s) {
		return errutil.Explain(nil,"validate failed on %s.%s")
	}`, str, receiverType, fieldName)

	if optional {
		str = fmt.Sprintf(`if x.%s != nil { %s }`, fieldName, str)
	}
	return &str, nil
}

// genValidateExpr recursively generates Go code for a validation expression
func genValidateExpr(fieldName, fieldType string, expr vidl.Expr, funcs map[string]ValidateFunc) (string, error) {
	switch x := expr.(type) {
	case vidl.BinaryExpr:
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

	case vidl.UnaryExpr:
		str, err := genValidateExpr(fieldName, fieldType, x.Expr, funcs)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s%s", x.Op, str), nil

	case *vidl.FuncCall:
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
			str, err := genValidateExpr(fieldName, fieldType, arg, funcs)
			if err != nil {
				return "", err
			}
			args = append(args, str)
		}
		return fmt.Sprintf("%s(%s)", x.Name, strings.Join(args, ", ")), nil

	case *vidl.InnerExpr:
		str, err := genValidateExpr(fieldName, fieldType, x.Expr, funcs)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s)", str), nil

	case vidl.PrimaryExpr:
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
func convertRPC(r tidl.RPC) (RPC, error) {

	// Retrieve the required "path" annotation
	path, ok := tidl.GetAnnotation(r.Annotations, "path")
	if !ok {
		return RPC{}, errutil.Explain(nil, `annotation "path" not found in rpc %s`, r.Name)
	}
	if path.Value == nil {
		return RPC{}, errutil.Explain(nil, `annotation "path" value is nil in rpc %s`, r.Name)
	}

	// Retrieve the required "method" annotation
	method, ok := tidl.GetAnnotation(r.Annotations, "method")
	if !ok {
		return RPC{}, errutil.Explain(nil, `annotation "method" not found in rpc %s`, r.Name)
	}
	if method.Value == nil {
		return RPC{}, errutil.Explain(nil, `annotation "method" value is nil in rpc %s`, r.Name)
	}

	return RPC{
		Name:     r.Name,
		Request:  r.Request.Name,
		Response: r.Response.UserType.Name,
		Stream:   r.Response.Stream,
		Path:     strings.Trim(*path.Value, `"`),
		Method:   strings.ToUpper(strings.Trim(*method.Value, `"`)),
		Comment:  formatComment(r.Comments),
	}, nil
}

// formatComment converts a tidl.Comments into Go comments.
func formatComment(c tidl.Comments) string {
	var comment string
	for _, s := range c.Above {
		comment += s.Text[0]
	}
	if c.Right != nil {
		if c.Above != nil {
			comment += "\n"
		}
		comment += strings.Join(c.Right.Text, "\n")
	}
	return comment
}
