package parser

import (
	"fmt"

	"golang.org/x/tools/go/packages"
)

// ImportGraph maps package -> list of imported packages
type ImportGraph map[string][]string

// ExtractImports loads Go packages and extracts import relationships.
func ExtractImports(pattern string) (ImportGraph, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedModule,
	}

	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	graph := make(ImportGraph)

	for _, pkg := range pkgs {
		if pkg.PkgPath == "" {
			continue // skip unnamed packages
		}

		for _, imp := range pkg.Imports {
			if imp.PkgPath == "" {
				continue
			}
			graph[pkg.PkgPath] = append(graph[pkg.PkgPath], imp.PkgPath)
		}
	}

	return graph, nil
}
