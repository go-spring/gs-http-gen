/*
 * Copyright 2025 The Go-Spring Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package httpidl

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-spring/gs-http-gen/lib/pathidl"
	"github.com/lvan100/golib/errutil"
	"github.com/lvan100/golib/ordered"
)

// BuiltinFuncs is a set of built-in validation functions
var BuiltinFuncs = map[string]struct{}{
	"len": {},
}

// ValidateFunc represents a validate function.
type ValidateFunc struct {
	Name string
	Type string
}

// RequestMeta represents the metadata of a request type.
type RequestMeta struct {
	OnForm bool
}

// Project represents a collection of IDL files and their associated meta-information.
type Project struct {
	Meta  *MetaInfo
	Files map[string]Document
	Reqs  map[string]RequestMeta
	Funcs map[string]ValidateFunc
}

// ParseDir scans the specified directory for IDL files (*.idl) and a meta.json file.
// It parses each file into a Document structure and validates cross-file type references.
func ParseDir(dir string) (Project, error) {

	p, err := loadProject(dir)
	if err != nil {
		return Project{}, err
	}

	// Validate that all used types are defined
	if err = checkUserTypes(p); err != nil {
		return Project{}, err
	}

	// process errcode
	if err = mergeErrcode(p); err != nil {
		return Project{}, err
	}

	// process types
	if err = processTypes(p); err != nil {
		return Project{}, err
	}

	// process validation
	if err = checkValidate(p); err != nil {
		return Project{}, err
	}

	for _, doc := range p.Files {
		for i := range doc.RPCs {
			rpc := doc.RPCs[i]
			segments, err := pathidl.Parse(rpc.Path)
			if err != nil {
				return Project{}, errutil.Explain(err, `failed to parse path %s`, rpc.Path)
			}
			params := make(map[string]string)
			for _, seg := range segments {
				if seg.Type == pathidl.Static {
					continue
				}
				params[seg.Value] = ""
			}
			srcType, ok := GetType(p.Files, rpc.Request)
			if !ok {
				return Project{}, errutil.Explain(nil, "type %s is used but not defined", rpc.Request)
			}
			for _, f := range srcType.Fields {
				if f.Binding == nil || f.Binding.Source != "path" {
					continue
				}
				if _, ok = params[f.Binding.Field]; !ok {
					err = errutil.Explain(nil, "path parameter %s not found in request type %s", f.Binding.Field, rpc.Request)
					return Project{}, err
				}
				params[f.Binding.Field] = f.Name
			}
			for k, s := range params {
				if s == "" {
					err = errutil.Explain(nil, "path parameter %s not found in request type %s", k, rpc.Request)
					return Project{}, err
				}
			}
			rpc.PathSegments = segments
			rpc.PathParams = params
			doc.RPCs[i] = rpc
		}
	}

	return p, nil
}

// loadProject loads the project from the specified directory.
func loadProject(dir string) (Project, error) {
	var p Project

	entries, err := os.ReadDir(dir)
	if err != nil {
		return Project{}, errutil.Explain(nil, "read dir %s error: %w", dir, err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		fileName := e.Name()

		// Parse meta.json file if found
		if fileName == "meta.json" {
			var b []byte
			fileName = filepath.Join(dir, fileName)
			if b, err = os.ReadFile(fileName); err != nil {
				return Project{}, errutil.Explain(nil, "read file %s error: %w", fileName, err)
			}
			if err = json.Unmarshal(b, &p.Meta); err != nil {
				return Project{}, errutil.Explain(nil, "parse file %s error: %w", fileName, err)
			}
			continue
		}

		// Skip non-IDL files
		if !strings.HasSuffix(fileName, ".idl") {
			continue
		}

		var b []byte
		fileName = filepath.Join(dir, fileName)
		if b, err = os.ReadFile(fileName); err != nil {
			return Project{}, errutil.Explain(nil, "read file %s error: %w", fileName, err)
		}

		doc, validateFuncs, err := ParseIDL(b)
		if err != nil {
			return Project{}, errutil.Explain(nil, "parse file %s error: %w", fileName, err)
		}
		p.Files[e.Name()] = doc

		// ...
		for _, r := range doc.RPCs {
			p.Reqs[r.Request] = RequestMeta{
				OnForm: strings.HasPrefix(r.ContentType, "application/x-www-form-urlencoded"),
			}
		}

		// record validate func
		for name, f := range validateFuncs {
			if v, ok := p.Funcs[name]; !ok {
				p.Funcs[name] = f
			} else if v.Type != f.Type {
				return Project{}, errutil.Explain(nil, "validate func %s is defined multiple times", name)
			}
		}
	}

	if p.Meta == nil {
		return Project{}, errutil.Explain(nil, "no meta file")
	}
	if len(p.Files) == 0 {
		return Project{}, errutil.Explain(nil, "no idl file")
	}
	return p, nil
}

// checkUserTypes checks if all user-defined types are defined.
func checkUserTypes(p Project) error {
	userTypes := map[string]struct{}{}
	definedTypes := make(map[string]struct{})
	for _, doc := range p.Files {
		for k := range doc.EnumTypes {
			definedTypes[k] = struct{}{}
		}
		for k := range doc.TypeTypes {
			definedTypes[k] = struct{}{}
		}
		maps.Copy(userTypes, doc.UserTypes)
	}
	for k := range userTypes {
		if _, ok := definedTypes[k]; !ok {
			return errutil.Explain(nil, "type %s is used but not defined", k)
		}
	}
	return nil
}

// mergeErrcode merges error codes from different files.
func mergeErrcode(p Project) error {
	files := ordered.MapKeys(p.Files)
	for _, file := range files {
		doc := p.Files[file]
		for _, e := range doc.Enums {
			if !e.Extends {
				continue
			}
			t, fileName := FindEnum(p.Files, e.Name)
			if fileName == nil {
				return errutil.Explain(nil, "enum %s is used but not defined", e.Name)
			}
			for _, field := range e.Fields {
				field.ExtendsFrom = fileName
				t.Fields = append(t.Fields, field)
			}
		}
	}
	return nil
}

// processTypes processes the types in the project.
func processTypes(p Project) error {
	for _, doc := range p.Files {
		for i := range doc.Types {
			t := doc.Types[i]
			t.Fields = t.RawFields
			if t.GenericParam != nil { // generic type, need instance
				// do nothing ...
			} else if t.InstType != nil { // generic type instance
				srcType, ok := GetType(p.Files, t.InstType.BaseName)
				if !ok {
					return errutil.Explain(nil, "type %s is used but not defined", t.InstType.BaseName)
				}
				var fields []TypeField
				for _, f := range srcType.Fields {
					f.Type = replaceGenericType(f.Type, *srcType.GenericParam, t.InstType.GenericType)
					fields = append(fields, f)
				}
				t.Fields = fields
			} else if t.Embedded {
				var fields []TypeField
				for _, f := range t.Fields {
					if e, ok := f.Type.(EmbedType); ok {
						srcType, ok := GetType(p.Files, e.Name)
						if !ok {
							return errutil.Explain(nil, "type %s is used but not defined", e.Name)
						}
						fields = append(fields, srcType.Fields...)
					} else {
						fields = append(fields, f)
					}
				}
				t.Fields = fields
			}

			if v, ok := p.Reqs[t.Name]; ok {
				t.Request = true
				t.OnRequest = true
				t.OnForm = v.OnForm
			}
			doc.Types[i] = t // update
		}
	}
	return nil
}

// replaceGenericType replaces a generic type with a concrete type.
func replaceGenericType(t TypeDefinition, genericName string, genericType TypeDefinition) TypeDefinition {
	switch u := t.(type) {
	case UserType:
		if u.Name == genericName {
			return genericType
		}
		return u
	case ListType:
		u.Item = replaceGenericType(u.Item, genericName, genericType)
		return u
	case MapType:
		u.Value = replaceGenericType(u.Value, genericName, genericType)
		return u
	default:
		return t
	}
}

// checkValidate checks if all types have validated expressions.
func checkValidate(p Project) error {
	for _, doc := range p.Files {
		for _, t := range doc.Types {
			if !t.Request {
				continue
			}
			if _, err := checkTypeValidate(p.Files, t); err != nil {
				return errutil.Explain(err, `failed to get type attr of type %s`, t.Name)
			}
		}
	}
	return nil
}

// checkTypeValidate checks if a type has validated expressions.
func checkTypeValidate(files map[string]Document, t Type) (bool, error) {
	t.OnRequest = true
	for i, f := range t.Fields {
		ok, err := checkUserTypeValidate(files, f.Type)
		if err != nil {
			return false, err
		}
		if ok {
			f.ValidateNested = true
			t.Fields[i] = f
		}
		if f.Required || f.ValidateExpr != nil || f.ValidateNested {
			t.Validate = true
		}
	}
	fileName, index := FindType(files, t.Name)
	files[fileName].Types[index] = t // update
	return t.Validate, nil
}

// checkUserTypeValidate checks if a user-defined type has validated expressions.
func checkUserTypeValidate(files map[string]Document, t TypeDefinition) (bool, error) {
	switch x := t.(type) {
	case UserType:
		srcType, ok := GetType(files, x.Name)
		if !ok {
			if _, ok = GetEnum(files, x.Name); !ok {
				return false, errutil.Explain(nil, "type %s is used but not defined", x.Name)
			}
			return false, nil
		}
		return checkTypeValidate(files, srcType)
	case ListType:
		return checkUserTypeValidate(files, x.Item)
	case MapType:
		return checkUserTypeValidate(files, x.Value)
	default: // for linter
	}
	return false, nil
}
