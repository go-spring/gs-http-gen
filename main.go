package main

import (
	"fmt"
	"os"

	"github.com/go-spring/gs-gen/gen"
	"github.com/go-spring/gs-gen/gen/generator"
	"github.com/spf13/cobra"
)

const Version = "v0.0.1"

func main() {
	var (
		version  bool
		language string
		genMode  string
	)

	root := &cobra.Command{
		Use:          "gs-gen",
		Short:        "A http code gen tool",
		SilenceUsage: true,
	}

	root.Flags().BoolVar(&version, "version", false, "show version")
	root.Flags().StringVar(&language, "lang", "go", "language, go/php/java")
	root.Flags().StringVar(&genMode, "mode", "server", "server, client, type")

	root.RunE = func(cmd *cobra.Command, args []string) error {
		if version {
			fmt.Println(root.Short)
			fmt.Println(Version)
			return nil
		}
		config := &generator.Config{
			ProjectDir: ".",
			Version:    Version,
		}
		switch genMode {
		case "server":
			config.Server = true
		case "client":
			config.Client = true
		default:
		}
		return gen.Gen(language, config)
	}

	if err := root.Execute(); err != nil {
		os.Exit(-1)
	}
}
