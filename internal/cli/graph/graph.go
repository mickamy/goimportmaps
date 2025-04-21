package graph

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mickamy/goimportmaps/internal/parser"
)

var Cmd = &cobra.Command{
	Use:   "graph [pattern]",
	Short: "Print package dependency graph",
	Long: `The graph command analyzes your Go packages and prints their internal import relationships.

Use it to inspect how packages depend on each other, or to generate raw dependency data before formatting it as a graph.

This is useful for understanding the structure of your project and preparing for visualization (e.g., Mermaid output).`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pattern := args[0]

		graph, err := parser.ExtractImports(pattern)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		for from, toList := range graph {
			for _, to := range toList {
				fmt.Printf("%s -> %s\n", from, to)
			}
		}
	},
}
