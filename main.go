package main

import (
	"fmt"
	"os"

	"github.com/go-spring/gs-http-gen/gen"
	"github.com/go-spring/gs-http-gen/gen/generator"
	"github.com/spf13/cobra"
)

const Version = "v0.0.1"

func main() {
	var (
		version  bool
		language string
		server   bool
		client   bool
		output   string
		pkgName  string
	)

	root := &cobra.Command{
		Use:          "gs-http-gen",
		Short:        "A http code gen tool",
		SilenceUsage: true,
	}

	root.Flags().BoolVar(&version, "version", false, "show version")
	root.Flags().StringVar(&language, "language", "go", "go/php/java")
	root.Flags().BoolVar(&server, "server", false, "gen server code")
	root.Flags().BoolVar(&client, "client", false, "gen client code")
	root.Flags().StringVar(&output, "output", ".", "output directory")
	root.Flags().StringVar(&pkgName, "package", "proto", "package name")

	root.RunE = func(cmd *cobra.Command, args []string) error {
		if version {
			fmt.Println(root.Short)
			fmt.Println(Version)
			return nil
		}
		config := &generator.Config{
			IDLDir:  ".",
			OutDir:  output,
			Version: Version,
			Server:  server,
			Client:  client,
			PkgName: pkgName,
		}
		return gen.Gen(language, config)
	}

	if err := root.Execute(); err != nil {
		os.Exit(-1)
	}
}
