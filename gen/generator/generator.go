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

package generator

import (
	"fmt"

	"github.com/go-spring/gs-http-gen/lib/tidl"
)

// Config holds the configuration options for the code generator.
type Config struct {
	IDLSrcDir    string // Directory containing source IDL files
	OutputDir    string // Directory where generated code will be written
	EnableServer bool   // Whether to generate server code
	EnableClient bool   // Whether to generate client code
	PackageName  string // Go package name for generated code
	ToolVersion  string // Version of the code generation tool
}

var generators = map[string]Generator{}

// Generator defines the interface that any language-specific generator
// must implement. The Gen method generates code based on the given
// configuration, documents, and metadata.
type Generator interface {
	Gen(config *Config, files map[string]tidl.Document, meta *tidl.MetaInfo) error
}

// GetGenerator retrieves a registered generator for a given language.
func GetGenerator(language string) (Generator, bool) {
	g, ok := generators[language]
	return g, ok
}

// RegisterGenerator registers a new generator for the given language.
func RegisterGenerator(language string, g Generator) {
	if _, ok := generators[language]; ok {
		panic(fmt.Errorf("duplicate generator for %s", language))
	}
	generators[language] = g
}

// GetEnum searches all documents for an enum type with the given name.
func GetEnum(files map[string]tidl.Document, name string) (tidl.Enum, bool) {
	for _, doc := range files {
		for _, e := range doc.Enums {
			if CapitalizeASCII(e.Name) == name {
				return e, true
			}
		}
	}
	return tidl.Enum{}, false
}

// GetType searches all documents for a type with the given name.
func GetType(files map[string]tidl.Document, name string) (tidl.Type, bool) {
	for _, doc := range files {
		for _, t := range doc.Types {
			if CapitalizeASCII(t.Name) == name {
				return t, true
			}
		}
	}
	return tidl.Type{}, false
}

// CapitalizeASCII capitalizes the first ASCII letter of a string.
func CapitalizeASCII(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-'a'+'A') + s[1:]
	}
	return s
}
