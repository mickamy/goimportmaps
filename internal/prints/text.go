package prints

import (
	"fmt"
	"io"

	"github.com/mickamy/goimportmaps"
	"github.com/mickamy/goimportmaps/internal/module"
)

func Text(w io.Writer, graph goimportmaps.Graph, modulePath string) {
	for from, toList := range graph {
		for _, to := range toList {
			from = module.Shorten(from, modulePath)
			to = module.Shorten(to, modulePath)
			_, _ = fmt.Fprintf(w, "  %s --> %s\n", from, to)
		}
	}
}
