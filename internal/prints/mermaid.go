package prints

import (
	"fmt"
	"io"
	"sort"

	"github.com/mickamy/goimportmaps"
)

func Mermaid(w io.Writer, imports goimportmaps.Graph) {
	_, _ = fmt.Fprintln(w, "```mermaid")
	_, _ = fmt.Fprintln(w, "graph TD")

	// To ensure consistent output
	keys := make([]string, 0, len(imports))
	for k := range imports {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, from := range keys {
		toList := imports[from]
		sort.Strings(toList)
		for _, to := range toList {
			_, _ = fmt.Fprintf(w, "  %s --> %s\n", from, to)
		}
	}

	_, _ = fmt.Fprintln(w, "```")
}
