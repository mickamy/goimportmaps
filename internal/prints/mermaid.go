package prints

import (
	"fmt"
	"io"
	"sort"

	"github.com/mickamy/goimportmaps"
	"github.com/mickamy/goimportmaps/internal/config"
)

func Mermaid(w io.Writer, graph goimportmaps.Graph, modulePath string, violations []config.Violation) {
	_, _ = fmt.Fprintln(w, "```mermaid")
	_, _ = fmt.Fprintln(w, "graph TD")

	violationMap := make(map[string]map[string]bool)
	for _, v := range violations {
		if violationMap[v.From] == nil {
			violationMap[v.From] = make(map[string]bool)
		}
		violationMap[v.From][v.To] = true
	}

	keys := make([]string, 0, len(graph))
	for k := range graph {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, from := range keys {
		toList := graph[from]
		sort.Strings(toList)
		for _, to := range toList {
			from = Shorten(from, modulePath)
			to = Shorten(to, modulePath)
			if violationMap[from][to] {
				_, _ = fmt.Fprintf(w, "  %s --> %s %% ❌ Violation\n", from, to)
			} else {
				_, _ = fmt.Fprintf(w, "  %s --> %s\n", from, to)
			}
		}
	}

	_, _ = fmt.Fprintln(w, "```")

	_, _ = fmt.Fprintln(w, "```")
}
