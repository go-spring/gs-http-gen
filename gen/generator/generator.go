package generator

import (
	"fmt"
	"strings"

	"github.com/go-spring/gs-http-gen/lib/parser"
)

type Config struct {
	IDLDir  string
	OutDir  string
	Version string
	Server  bool
	Client  bool
	PkgName string
}

var generators = map[string]Generator{}

type Generator interface {
	Gen(config *Config, files map[string]parser.Document, meta *parser.MetaInfo) error
}

func GetGenerator(language string) (Generator, bool) {
	g, ok := generators[language]
	return g, ok
}

func RegisterGenerator(language string, g Generator) {
	if _, ok := generators[language]; ok {
		panic(fmt.Errorf("duplicate generator for %s", language))
	}
	generators[language] = g
}

func GetType(files map[string]parser.Document, name string) *parser.Type {
	for _, doc := range files {
		for _, t := range doc.Types {
			if t.Name == name {
				return t
			}
		}
	}
	return nil
}

func GetAnnotation(arr []parser.Annotation, name string) (parser.Annotation, bool) {
	for _, a := range arr {
		if a.Key == name {
			return a, true
		}
	}
	return parser.Annotation{}, false
}

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
