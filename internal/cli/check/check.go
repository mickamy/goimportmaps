package check

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mickamy/goimportmaps/internal/config"
	"github.com/mickamy/goimportmaps/internal/parser"
)

var (
	mode = "forbidden"
)

var Cmd = &cobra.Command{
	Use:   "check [pattern]",
	Short: "Check for forbidden imports defined in .goimportmaps.yaml",
	Long: `Check your Go package dependencies against forbidden import rules.

Rules must be defined in a .goimportmaps.yaml file at the project root.
If any violations are found, they will be printed to stderr and the program will exit with code 1.`,
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

		Run(cfg, mode, args[0])
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&mode, "mode", "m", "forbidden", "check mode (forbidden or allowed)")
}

func Run(cfg *config.Config, mode config.Mode, pattern string) {
	data, err := parser.ExtractImports(pattern)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	violations := cfg.Validate(data, mode)
	if len(violations) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "\n🚨 %d violation(s) found\n\n", len(violations))

		for _, violation := range violations {
			_, _ = fmt.Fprintln(os.Stderr, "🚨 Violation:", violation.Message)
		}
		os.Exit(1)
	}
}
