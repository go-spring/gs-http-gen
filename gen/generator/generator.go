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
	"strings"

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

// ToPascal converts a snake_case string to PascalCase.
// For example: "hello_world" becomes "HelloWorld".
func ToPascal(s string) string {
	var sb strings.Builder
	parts := strings.Split(s, "_")
	for _, part := range parts {
		if part == "" {
			continue
		}
		c := part[0]
		if 'a' <= c && c <= 'z' {
			c = c - 'a' + 'A'
		}
		sb.WriteByte(c)
		if len(part) > 1 {
			sb.WriteString(part[1:])
		}
	}
	return sb.String()
}
