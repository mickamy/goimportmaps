package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mickamy/goimportmaps/internal/cli/version"
)

var cmd = &cobra.Command{
	Use:   "goimportmaps",
	Short: "Visualize and validate Go package dependencies",
	Long: `goimportmaps is a CLI tool that analyzes and visualizes internal package dependencies in your Go project.
It helps maintain architectural integrity by detecting forbidden imports and generating dependency graphs
in formats like Mermaid, Graphviz, or HTML.`,
}

func init() {
	cmd.AddCommand(version.Cmd)
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
