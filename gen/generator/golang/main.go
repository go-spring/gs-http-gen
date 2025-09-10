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

package golang

import (
	"go/format"
	"os"
	"sort"

	"github.com/go-spring/gs-http-gen/gen/generator"
	"github.com/go-spring/gs-http-gen/lib/tidl"
)

// Context holds all necessary information used during code generation.
type Context struct {
	config *generator.Config        // Generator configuration options
	meta   *tidl.MetaInfo           // Service metadata (name, version, etc.)
	files  map[string]tidl.Document // Parsed TIDL documents keyed by file name
	funcs  map[string]ValidateFunc  // Collected validation functions
}

type Generator struct{}

// Gen is the main entry point for generating code.
func (g *Generator) Gen(config *generator.Config, files map[string]tidl.Document, meta *tidl.MetaInfo) error {
	ctx := Context{
		config: config,
		meta:   meta,
		files:  files,
		funcs:  make(map[string]ValidateFunc),
	}

	var rpcs []tidl.RPC

	// Collect all RPC definitions
	for fileName, doc := range files {
		if err := g.genType(ctx, fileName, doc); err != nil {
			return err
		}
		rpcs = append(rpcs, doc.RPCs...)
	}

	sort.Slice(rpcs, func(i, j int) bool {
		return rpcs[i].Name < rpcs[j].Name
	})

	// Generate server code if enabled in the configuration
	if config.EnableServer {
		if err := g.genValidate(ctx); err != nil {
			return err
		}
		if err := g.genServer(ctx, rpcs); err != nil {
			return err
		}
	}

	// Generate client code if enabled in the configuration
	if config.EnableClient {
		if err := g.genClient(ctx, rpcs); err != nil {
			return err
		}
	}

	return nil
}

// formatFile formats the given Go source code and writes it to a file.
func formatFile(fileName string, b []byte) error {
	b, err := format.Source(b)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, b, os.ModePerm)
}
