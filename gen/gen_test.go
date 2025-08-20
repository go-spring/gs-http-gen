package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/go-spring/gs-gen/gen/generator"
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
	}
	if err = json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}

	idlDir := filepath.Join(dir, "idl")
	for lang, c := range m {
		projectDir := filepath.Join(dir, lang)

		if c.Server {
			dstDir := filepath.Join(projectDir, "server")
			copyFiles(t, idlDir, dstDir)
			config := &generator.Config{
				ProjectDir: dstDir,
				Version:    "v0.0.0",
				Server:     true,
			}
			if err = Gen(lang, config); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func copyFiles(t *testing.T, srcDir, dstDir string) {
	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		t.Fatal(err)
	}
	dstDir, err = filepath.Abs(dstDir)
	if err != nil {
		t.Fatal(err)
	}
	srcDir = filepath.Join(srcDir, "*")

	script := fmt.Sprintf("cp %s %s", srcDir, dstDir)
	cmd := exec.Command("/bin/bash", "-c", script)
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %s", err, b)
	}
}
