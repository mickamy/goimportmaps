package check

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mickamy/goimportmaps/internal/config"
	"github.com/mickamy/goimportmaps/internal/parser"
)

var Cmd = &cobra.Command{
	Use:   "check [pattern]",
	Short: "Check for forbidden imports defined in .goimportmaps.yaml",
	Long: `Check your Go package dependencies against forbidden import rules.

Rules must be defined in a .goimportmaps.yaml file at the project root.
If any violations are found, they will be printed to stderr and the program will exit with code 1.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		Run(cfg, args[0])
	},
}

func Run(cfg *config.Config, pattern string) {
	data, err := parser.ExtractImports(pattern)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	violations := cfg.Validate(data)
	if len(violations) > 0 {
		for _, violation := range violations {
			_, _ = fmt.Fprintln(os.Stderr, "🚨 Violation:", violation.Message)
		}
		os.Exit(1)
	}
}
