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
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-spring/gs-http-gen/gen/generator"
)

func TestGen(t *testing.T) {
	root := "testdata"
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(root, e.Name())
		testProject(t, dir)
	}
}

func testProject(t *testing.T, dir string) {
	b, err := os.ReadFile(filepath.Join(dir, "test.json"))
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]struct {
		Server bool
		Client bool
	}
	if err = json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}

	idlDir := filepath.Join(dir, "idl")
	for lang, c := range m {
		outDir := filepath.Join(dir, lang, "proto")
		_ = os.RemoveAll(outDir)
		_ = os.Mkdir(outDir, os.ModePerm)
		config := &generator.Config{
			IDLSrcDir:    idlDir,
			OutputDir:    outDir,
			EnableServer: c.Server,
			EnableClient: c.Client,
			PackageName:  "proto",
			ToolVersion:  "v0.0.1",
		}
		if err = Gen(lang, config); err != nil {
			t.Fatal(err)
		}
	}
}
