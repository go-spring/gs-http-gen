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

package main

import (
	"fmt"
	"os"

	"github.com/go-spring/gs-http-gen/gen"
	"github.com/go-spring/gs-http-gen/gen/generator"
	"github.com/spf13/cobra"
)

// ToolVersion defines the current version of gs-http-gen tool.
const ToolVersion = "v0.0.3"

func main() {
	var (
		showVersion  bool
		language     string
		outputDir    string
		goPackage    string
		enableServer bool
		enableClient bool
	)

	root := &cobra.Command{
		Use:   "gs-http-gen",
		Short: "A code generation tool for HTTP services based on IDL files",
		Long: `gs-http-gen is a code generation tool that reads service definitions
from IDL files and generates server and/or client code in Go (default),
PHP, Java, or other supported languages.`,
		SilenceUsage: true,
	}

	root.Flags().BoolVar(&showVersion, "version", false, "Display the version of gs-http-gen tool")
	root.Flags().StringVar(&language, "language", "go", "Target language for code generation (go, php, java, etc.)")
	root.Flags().BoolVar(&enableServer, "server", false, "Generate server-side code")
	root.Flags().BoolVar(&enableClient, "client", false, "Generate client-side code")
	root.Flags().StringVar(&outputDir, "output", ".", "Output directory for generated code (default: current directory)")
	root.Flags().StringVar(&goPackage, "go_package", "proto", "Go package name for generated code")

	root.RunE = func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(root.Short)
			fmt.Println(ToolVersion)
			return nil
		}

		config := &generator.Config{
			IDLSrcDir:    ".",
			OutputDir:    outputDir,
			EnableServer: enableServer,
			EnableClient: enableClient,
			GoPackage:    goPackage,
			ToolVersion:  ToolVersion,
		}
		return gen.Gen(language, config)
	}

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
