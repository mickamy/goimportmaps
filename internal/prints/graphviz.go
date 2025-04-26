package prints

import (
	"fmt"
	"io"
	"sort"

	"github.com/mickamy/goimportmaps"
	"github.com/mickamy/goimportmaps/internal/module"
)

func Graphviz(w io.Writer, graph goimportmaps.Graph, modulePath string) {
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
			from = module.Shorten(from, modulePath)
			to = module.Shorten(to, modulePath)
			_, _ = fmt.Fprintf(w, "  %q -> %q;\n", from, to)
		}
	}

	_, _ = fmt.Fprintln(w, "}")
}
