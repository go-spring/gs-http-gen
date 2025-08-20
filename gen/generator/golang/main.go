package golang

import (
	"sort"

	"github.com/go-spring/gs-gen/gen/generator"
	"github.com/go-spring/gs-gen/lib/parser"
)

type Context struct {
	config *generator.Config
	meta   *parser.MetaInfo
	files  map[string]parser.Document
}

type Generator struct{}

func (g *Generator) Gen(config *generator.Config, files map[string]parser.Document, meta *parser.MetaInfo) error {
	ctx := Context{
		config: config,
		meta:   meta,
		files:  files,
	}

	if err := g.genMeta(meta); err != nil {
		return err
	}

	var rpcs []*parser.RPC
	for fileName, doc := range files {
		if err := g.genType(ctx, fileName, doc); err != nil {
			return err
		}
		for _, c := range doc.RPCs {
			rpcs = append(rpcs, c)
		}
	}

	sort.Slice(rpcs, func(i, j int) bool {
		return rpcs[i].Name < rpcs[j].Name
	})

	if config.Server {
		if err := g.genServer(ctx, rpcs); err != nil {
			return err
		}
	} else {
		if err := g.genClient(ctx, rpcs); err != nil {
			return err
		}
	}
	return nil
}
