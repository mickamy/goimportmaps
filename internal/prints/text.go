package prints

import (
	"fmt"
	"io"

	"github.com/mickamy/goimportmaps"
)

func Text(w io.Writer, graph goimportmaps.Graph) {
	for from, toList := range graph {
		for _, to := range toList {
			_, _ = fmt.Fprintf(w, "  %s --> %s\n", from, to)
		}
	}
}
