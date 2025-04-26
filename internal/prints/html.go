package prints

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"sort"
	"strings"

	"github.com/mickamy/goimportmaps"
	"github.com/mickamy/goimportmaps/internal/config"
	"github.com/mickamy/goimportmaps/internal/module"
)

//go:embed template.html
var htmlTemplate string

func HTML(w io.Writer, graph goimportmaps.Graph, modulePath string, violations []config.Violation) error {
	var buf bytes.Buffer
	buf.WriteString("graph TD\n")

	violationSet := make(map[string]bool)

	keys := make([]string, 0, len(graph))
	for from := range graph {
		keys = append(keys, from)
	}
	sort.Strings(keys)

	for _, from := range keys {
		toList := graph[from]
		sort.Strings(toList)
		for _, to := range toList {
			from = module.Shorten(from, modulePath)
			to = module.Shorten(to, modulePath)
			_, _ = fmt.Fprintf(&buf, "  %s --> %s\n", from, to)
		}
	}

	for _, v := range violations {
		violationSet[module.Shorten(v.Source, modulePath)] = true
		violationSet[module.Shorten(v.Import, modulePath)] = true
	}

	if len(violationSet) > 0 {
		buf.WriteString("  classDef violation stroke:#f00,stroke-width:4px;\n")

		nodes := make([]string, 0, len(violationSet))
		for n := range violationSet {
			nodes = append(nodes, n)
		}
		sort.Strings(nodes)

		_, _ = fmt.Fprintf(&buf, "  class %s violation;\n", strings.Join(nodes, ","))
	}

	tmpl, err := template.New("html").Funcs(template.FuncMap{
		"safe": func(s string) template.HTML { return template.HTML(s) },
	}).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(w, htmlTemplateData{
		Graph:          buf.String(),
		ViolationCount: len(violations),
	}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

type htmlTemplateData struct {
	Graph          string
	ViolationCount int
}
