package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mickamy/goimportmaps/internal/cli/check"
	"github.com/mickamy/goimportmaps/internal/cli/graph"
	"github.com/mickamy/goimportmaps/internal/cli/version"
	"github.com/mickamy/goimportmaps/internal/config"
	"github.com/mickamy/goimportmaps/internal/module"
	"github.com/mickamy/goimportmaps/internal/parser"
	"github.com/mickamy/goimportmaps/internal/prints"
)

var (
	format = "text"
	mode   = "forbidden"
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

		mode, err := config.NewMode(mode)
		if err != nil {
			return err
		}

		format, err := prints.NewFormat(format)
		if err != nil {
			return err
		}

		Run(cfg, mode, format, args[0])

		return nil
	},
}

func init() {
	cmd.AddCommand(check.Cmd)
	cmd.AddCommand(graph.Cmd)
	cmd.AddCommand(version.Cmd)

	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, mermaid, graphviz or html)")
	cmd.Flags().StringVarP(&mode, "mode", "m", "forbidden", "check mode (forbidden or allowed)")
}

func Run(cfg *config.Config, mode config.Mode, format prints.Format, pattern string) {
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

	violations := cfg.Validate(data, mode)

	switch format {
	case prints.FormatGraphviz:
		prints.Graphviz(os.Stdout, data, modulePath)
	case prints.FormatHTML:
		if err := prints.HTML(os.Stdout, data, modulePath, violations); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	case prints.FormatMermaid:
		prints.Mermaid(os.Stdout, data, modulePath, violations)
	case prints.FormatText:
		prints.Text(os.Stdout, data, modulePath)
	}

	if len(violations) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "\n🚨 %d violation(s) found\n\n", len(violations))

		for _, violation := range violations {
			_, _ = fmt.Fprintln(os.Stderr, "🚨 Violation:", violation.Message)
		}
	}
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
