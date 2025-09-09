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

package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-spring/gs-http-gen/gen/generator"
	"github.com/go-spring/gs-http-gen/gen/generator/golang"
	"github.com/go-spring/gs-http-gen/lib/tidl"
)

func init() {
	generator.RegisterGenerator("go", &golang.Generator{})
}

// Gen is the entry point for generating code for the given language.
func Gen(language string, config *generator.Config) error {
	g, ok := generator.GetGenerator(language)
	if !ok {
		return fmt.Errorf("unsupported language: %s", language)
	}
	files, meta, err := parse(config.IDLSrcDir)
	if err != nil {
		return err
	}
	if meta == nil {
		return fmt.Errorf("no meta file")
	}
	if len(files) == 0 {
		return fmt.Errorf("no idl file")
	}
	return g.Gen(config, files, meta)
}

// parse scans the given directory for IDL files and the meta.json file.
func parse(dir string) (files map[string]tidl.Document, meta *tidl.MetaInfo, err error) {
	files = make(map[string]tidl.Document)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fileName := e.Name()

		// Parse meta.json file if found
		if fileName == "meta.json" {
			if meta, err = parseMeta(dir, fileName); err != nil {
				return nil, nil, err
			}
			continue
		}

		// Skip non-IDL files
		if !strings.HasSuffix(fileName, ".idl") {
			continue
		}

		var doc tidl.Document
		doc, err = parseFile(dir, fileName)
		if err != nil {
			return nil, nil, err
		}
		files[fileName] = doc
	}
	return
}

// parseMeta reads and parses the meta.json file to extract service metadata.
func parseMeta(dir string, fileName string) (*tidl.MetaInfo, error) {
	fileName = filepath.Join(dir, fileName)
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return tidl.ParseMeta(string(b))
}

// parseFile reads and parses a single IDL file into a document.
func parseFile(dir string, fileName string) (tidl.Document, error) {
	fileName = filepath.Join(dir, fileName)
	b, err := os.ReadFile(fileName)
	if err != nil {
		return tidl.Document{}, err
	}
	return tidl.Parse(string(b))
}
