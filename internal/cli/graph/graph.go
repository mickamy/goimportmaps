package graph

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mickamy/goimportmaps/internal/config"
	"github.com/mickamy/goimportmaps/internal/module"
	"github.com/mickamy/goimportmaps/internal/parser"
	"github.com/mickamy/goimportmaps/internal/prints"
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
	Cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, mermaid, graphviz or html)")
}

func Run(pattern string, format string) error {
	data, err := parser.ExtractImports(pattern)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	modulePath, err := module.Path()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	switch format {
	case "text":
		prints.Text(os.Stdout, data, modulePath)
	case "mermaid":
		prints.Mermaid(os.Stdout, data, modulePath, []config.Violation{})
	case "graphviz":
		prints.Graphviz(os.Stdout, data, modulePath)
	case "html":
		if err := prints.HTML(os.Stdout, data, modulePath, []config.Violation{}); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	default:
		return fmt.Errorf("unsupported format %s", format)
	}

	return nil
}
