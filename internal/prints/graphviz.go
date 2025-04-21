package prints

import (
	"fmt"
	"io"
	"sort"

	"github.com/mickamy/goimportmaps"
)

func Graphviz(w io.Writer, graph goimportmaps.Graph) {
	_, _ = fmt.Fprintln(w, "digraph G {")

	keys := make([]string, 0, len(graph))
	for k := range graph {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, from := range keys {
		toList := graph[from]
		sort.Strings(toList)
		for _, to := range toList {
			_, _ = fmt.Fprintf(w, "  %q -> %q;\n", from, to)
		}
	}

	_, _ = fmt.Fprintln(w, "}")
}
