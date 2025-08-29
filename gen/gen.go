package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-spring/gs-http-gen/gen/generator"
	"github.com/go-spring/gs-http-gen/gen/generator/golang"
	"github.com/go-spring/gs-http-gen/lib/parser"
)

func init() {
	generator.RegisterGenerator("go", &golang.Generator{})
}

func Gen(language string, config *generator.Config) error {
	g, ok := generator.GetGenerator(language)
	if !ok {
		return fmt.Errorf("unsupported language: %s", language)
	}
	files, meta, err := parse(config.IDLDir)
	if err != nil {
		return err
	}
	if meta == nil {
		return fmt.Errorf("no meta file")
	}
	if len(files) == 0 {
		return fmt.Errorf("no idl file")
	}
	if err = parser.Verify(files, meta); err != nil {
		return err
	}
	return g.Gen(config, files, meta)
}

func parse(dir string) (files map[string]parser.Document, meta *parser.MetaInfo, err error) {
	files = make(map[string]parser.Document)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fileName := e.Name()

		if fileName == "meta.json" {
			if meta, err = parseMeta(dir, fileName); err != nil {
				return nil, nil, err
			}
			continue
		}

		if !strings.HasSuffix(fileName, ".idl") {
			continue
		}

		var doc parser.Document
		doc, err = parseFile(dir, fileName)
		if err != nil {
			return nil, nil, err
		}
		files[fileName] = doc
	}
	return
}

func parseMeta(dir string, fileName string) (*parser.MetaInfo, error) {
	fileName = filepath.Join(dir, fileName)
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return parser.ParseMeta(string(b))
}

func parseFile(dir string, fileName string) (parser.Document, error) {
	fileName = filepath.Join(dir, fileName)
	b, err := os.ReadFile(fileName)
	if err != nil {
		return parser.Document{}, err
	}
	return parser.Parse(string(b))
}
