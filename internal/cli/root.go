package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mickamy/goimportmaps/internal/cli/check"
	"github.com/mickamy/goimportmaps/internal/cli/graph"
	"github.com/mickamy/goimportmaps/internal/cli/version"
	"github.com/mickamy/goimportmaps/internal/config"
	"github.com/mickamy/goimportmaps/internal/parser"
	"github.com/mickamy/goimportmaps/internal/prints"
)

var (
	format = "text"
)

var cmd = &cobra.Command{
	Use:   "goimportmaps",
	Short: "Visualize and validate Go package dependencies",
	Long: `goimportmaps is a CLI tool that helps you understand and enforce 
the architecture of your Go projects by analyzing internal package imports.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return Run(cfg, args[0])
	},
}

func init() {
	cmd.AddCommand(check.Cmd)
	cmd.AddCommand(graph.Cmd)
	cmd.AddCommand(version.Cmd)

	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, mermaid or graphviz)")
}

func Run(cfg *config.Config, pattern string) error {
	data, err := parser.ExtractImports(pattern)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch format {
	case "text":
		prints.Text(os.Stdout, data)
	case "mermaid":
		prints.Mermaid(os.Stdout, data)
	case "graphviz":
		prints.Graphviz(os.Stdout, data)
	default:
		return fmt.Errorf("unsupported format %s", format)
	}

	violations := cfg.Validate(data)
	if len(violations) > 0 {
		for _, violation := range violations {
			_, _ = fmt.Fprintln(os.Stderr, violation)
		}
		os.Exit(1)
	}

	return nil
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
