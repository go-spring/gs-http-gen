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
		os.RemoveAll(outDir)
		os.Mkdir(outDir, os.ModePerm)
		config := &generator.Config{
			IDLDir:  idlDir,
			OutDir:  outDir,
			Version: "v0.0.0",
			Server:  c.Server,
			Client:  c.Client,
			PkgName: "proto",
		}
		if err = Gen(lang, config); err != nil {
			t.Fatal(err)
		}
	}
}
