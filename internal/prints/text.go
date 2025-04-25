package prints

import (
	"fmt"
	"io"

	"github.com/mickamy/goimportmaps"
)

func Text(w io.Writer, graph goimportmaps.Graph, modulePath string) {
	for from, toList := range graph {
		for _, to := range toList {
			from = Shorten(from, modulePath)
			to = Shorten(to, modulePath)
			_, _ = fmt.Fprintf(w, "  %s --> %s\n", from, to)
		}
	}
}
