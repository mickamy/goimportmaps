package prints

import (
	"fmt"
	"io"
	"sort"

	"github.com/mickamy/goimportmaps"
	"github.com/mickamy/goimportmaps/internal/config"
	"github.com/mickamy/goimportmaps/internal/module"
)

func Mermaid(w io.Writer, graph goimportmaps.Graph, modulePath string, violations []config.Violation) {
	_, _ = fmt.Fprintln(w, "```mermaid")
	_, _ = fmt.Fprintln(w, "graph TD")

	violationMap := make(map[string]map[string]bool)
	for _, v := range violations {
		if violationMap[v.Source] == nil {
			violationMap[v.Source] = make(map[string]bool)
		}
		violationMap[v.Source][v.Import] = true
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
			from = module.Shorten(from, modulePath)
			to = module.Shorten(to, modulePath)
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
