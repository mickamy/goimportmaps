package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev (built from source; try GitHub release for accurate version)"

var Cmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show goimportmaps version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("goimportmaps version %s\n", version)
	},
}
