package parser

import (
	"fmt"

	"golang.org/x/tools/go/packages"

	"github.com/mickamy/goimportmaps"
)

// ExtractImports loads Go packages and extracts import relationships.
func ExtractImports(pattern string) (goimportmaps.Graph, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedModule,
	}

	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	graph := make(goimportmaps.Graph)

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
