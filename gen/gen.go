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
	files, meta, err := tidl.ParseDir(config.IDLSrcDir)
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
