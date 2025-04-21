package graph

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mickamy/goimportmaps/internal/graph"
	"github.com/mickamy/goimportmaps/internal/parser"
)

var (
	format = "text"
)

var Cmd = &cobra.Command{
	Use:   "graph [pattern]",
	Short: "Print package dependency graph",
	Long: `The graph command analyzes your Go packages and prints their internal import relationships.

Use it to inspect how packages depend on each other, or to generate raw dependency data before formatting it as a graph.

This is useful for understanding the structure of your project and preparing for visualization (e.g., Mermaid output).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return Run(args[0], format)
	},
}

func init() {
	Cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text or mermaid)")
}

func Run(pattern string, format string) error {
	data, err := parser.ExtractImports(pattern)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch format {
	case "text":
		for from, toList := range data {
			for _, to := range toList {
				fmt.Printf("%s -> %s\n", from, to)
			}
		}
		return nil
	case "mermaid":
		if err = graph.RenderMermaid(os.Stdout, data); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return nil
	default:
		return fmt.Errorf("unsupported format %s", format)
	}
}
