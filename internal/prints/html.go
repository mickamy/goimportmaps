package prints

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"sort"

	"github.com/mickamy/goimportmaps"
)

//go:embed template.html
var htmlTemplate string

func HTML(w io.Writer, graph goimportmaps.Graph) error {
	var buf bytes.Buffer
	buf.WriteString("graph TD\n")

	keys := make([]string, 0, len(graph))
	for k := range graph {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, from := range keys {
		toList := graph[from]
		sort.Strings(toList)
		for _, to := range toList {
			_, _ = fmt.Fprintf(&buf, "  %s --> %s\n", from, to)
		}
	}

	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(w, struct{ Graph string }{Graph: buf.String()}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
