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
	"github.com/mickamy/goimportmaps/internal/metrics"
	"github.com/mickamy/goimportmaps/internal/module"
)

//go:embed template.html
var htmlTemplate string

var htmlTemplateWithMetrics = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <title>Go Import Graph with Metrics</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <script src="https://cdn.jsdelivr.net/npm/mermaid@11.6.0/dist/mermaid.min.js"></script>
    <script>mermaid.initialize({ startOnLoad: true });</script>
    <style>
        body {
            font-family: system-ui, sans-serif;
            margin: 0;
            padding: 2rem;
            background: #f9fafb;
            color: #111;
        }
        h1, h2 {
            margin-bottom: 1.5rem;
        }
        h1 {
            font-size: 1.8rem;
        }
        h2 {
            font-size: 1.4rem;
            margin-top: 2rem;
        }
        .mermaid {
            background: #fff;
            padding: 1rem;
            border: 1px solid #ddd;
            border-radius: 8px;
            overflow-x: auto;
            margin-bottom: 2rem;
        }
        .metrics-table {
            background: #fff;
            border: 1px solid #ddd;
            border-radius: 8px;
            overflow: hidden;
        }
        table {
            width: 100%;
            border-collapse: collapse;
        }
        th, td {
            padding: 0.75rem;
            text-align: left;
            border-bottom: 1px solid #eee;
        }
        th {
            background: #f8f9fa;
            font-weight: 600;
        }
        .violation {
            background: #fef2f2;
            color: #dc2626;
        }
        .good {
            background: #f0fdf4;
            color: #16a34a;
        }
        .status {
            text-align: center;
            font-size: 1.2rem;
        }
        .instability-bar {
            width: 100px;
            height: 8px;
            background: #e5e7eb;
            border-radius: 4px;
            overflow: hidden;
        }
        .instability-fill {
            height: 100%;
            transition: width 0.3s ease;
        }
        .instability-low {
            background: #16a34a;
        }
        .instability-medium {
            background: #eab308;
        }
        .instability-high {
            background: #dc2626;
        }
        .violation-reasons {
            font-size: 0.85rem;
            margin-top: 0.25rem;
            color: #dc2626;
        }
        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            margin-bottom: 2rem;
        }
        .summary-card {
            background: #fff;
            padding: 1.5rem;
            border: 1px solid #ddd;
            border-radius: 8px;
            text-align: center;
        }
        .summary-value {
            font-size: 2rem;
            font-weight: 700;
            margin-bottom: 0.5rem;
        }
        .summary-label {
            color: #6b7280;
            font-size: 0.9rem;
        }
        @media (max-width: 768px) {
            body {
                padding: 1rem;
            }
            .mermaid {
                font-size: 0.85rem;
            }
            table {
                font-size: 0.9rem;
            }
            th, td {
                padding: 0.5rem;
            }
        }
    </style>
</head>
<body>
<h1>📦 Go Import Graph with Metrics</h1>

{{ if or .ViolationCount .CouplingViolations }}
<div class="summary">
    {{ if .ViolationCount }}
    <div class="summary-card">
        <div class="summary-value" style="color: #dc2626;">{{ .ViolationCount }}</div>
        <div class="summary-label">Import Violations</div>
    </div>
    {{ end }}
    {{ if .CouplingViolations }}
    <div class="summary-card">
        <div class="summary-value" style="color: #dc2626;">{{ .CouplingViolations }}</div>
        <div class="summary-label">Coupling Violations</div>
    </div>
    {{ end }}
</div>
{{ end }}

<h2>📊 Dependency Graph</h2>
<div class="mermaid">
    {{ .Graph | safe }}
</div>

{{ if .HasMetrics }}
<h2>📈 Coupling Metrics</h2>
<div class="metrics-table">
    <table>
        <thead>
            <tr>
                <th>Package</th>
                <th>Ca</th>
                <th>Ce</th>
                <th>Instability</th>
                <th>Status</th>
            </tr>
        </thead>
        <tbody>
            {{ range .PackageMetrics }}
            <tr class="{{ if .HasViolation }}violation{{ else }}good{{ end }}">
                <td>{{ .Package }}</td>
                <td>{{ .AfferentCoupling }}</td>
                <td>{{ .EfferentCoupling }}</td>
                <td>
                    <div style="display: flex; align-items: center; gap: 0.5rem;">
                        <span>{{ printf "%.2f" .Instability }}</span>
                        <div class="instability-bar">
                            <div class="instability-fill {{ if lt .Instability 0.3 }}instability-low{{ else if lt .Instability 0.7 }}instability-medium{{ else }}instability-high{{ end }}" 
                                 style="width: {{ printf "%.0f" (multiply .Instability 100) }}%"></div>
                        </div>
                    </div>
                    {{ if .ViolationReasons }}
                    <div class="violation-reasons">
                        {{ range .ViolationReasons }}• {{ . }}<br>{{ end }}
                    </div>
                    {{ end }}
                </td>
                <td class="status">{{ if .HasViolation }}🚨{{ else }}✅{{ end }}</td>
            </tr>
            {{ end }}
        </tbody>
    </table>
</div>
{{ end }}

<script>
// Add multiply template function support
</script>
</body>
</html>`

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

type PackageMetricsData struct {
	Package          string
	AfferentCoupling int
	EfferentCoupling int
	Instability      float64
	HasViolation     bool
	ViolationReasons []string
}

type htmlTemplateDataWithMetrics struct {
	Graph              string
	ViolationCount     int
	PackageMetrics     []PackageMetricsData
	HasMetrics         bool
	CouplingViolations int
}

func HTMLWithMetrics(w io.Writer, graph goimportmaps.Graph, modulePath string, violations []config.Violation, analysis *metrics.CouplingAnalysis, maxEfferent, maxAfferent int, maxInstability float64) error {
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
			shortFrom := module.Shorten(from, modulePath)
			shortTo := module.Shorten(to, modulePath)
			_, _ = fmt.Fprintf(&buf, "  %s --> %s\n", shortFrom, shortTo)
		}
	}

	for _, v := range violations {
		violationSet[module.Shorten(v.Source, modulePath)] = true
		violationSet[module.Shorten(v.Import, modulePath)] = true
	}

	// add coupling violation nodes
	couplingViolations := analysis.GetHighCouplingPackages(maxEfferent, maxAfferent, maxInstability)
	for _, pkg := range couplingViolations {
		violationSet[module.Shorten(pkg.Package, modulePath)] = true
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

	// Prepare metrics data
	var packageMetrics []PackageMetricsData
	if analysis != nil {
		packageKeys := make([]string, 0, len(analysis.Packages))
		for pkg := range analysis.Packages {
			packageKeys = append(packageKeys, pkg)
		}
		sort.Strings(packageKeys)

		for _, pkg := range packageKeys {
			metrics := analysis.Packages[pkg]
			shortPkg := module.Shorten(pkg, modulePath)

			hasViolation := metrics.EfferentCoupling > maxEfferent ||
				metrics.AfferentCoupling > maxAfferent ||
				metrics.Instability > maxInstability

			var reasons []string
			if metrics.EfferentCoupling > maxEfferent {
				reasons = append(reasons, fmt.Sprintf("High efferent coupling (%d > %d)", metrics.EfferentCoupling, maxEfferent))
			}
			if metrics.AfferentCoupling > maxAfferent {
				reasons = append(reasons, fmt.Sprintf("High afferent coupling (%d > %d)", metrics.AfferentCoupling, maxAfferent))
			}
			if metrics.Instability > maxInstability {
				reasons = append(reasons, fmt.Sprintf("High instability (%.2f > %.2f)", metrics.Instability, maxInstability))
			}

			packageMetrics = append(packageMetrics, PackageMetricsData{
				Package:          shortPkg,
				AfferentCoupling: metrics.AfferentCoupling,
				EfferentCoupling: metrics.EfferentCoupling,
				Instability:      metrics.Instability,
				HasViolation:     hasViolation,
				ViolationReasons: reasons,
			})
		}
	}

	tmpl, err := template.New("html").Funcs(template.FuncMap{
		"safe":     func(s string) template.HTML { return template.HTML(s) },
		"multiply": func(a, b float64) float64 { return a * b },
	}).Parse(htmlTemplateWithMetrics)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(w, htmlTemplateDataWithMetrics{
		Graph:              buf.String(),
		ViolationCount:     len(violations),
		PackageMetrics:     packageMetrics,
		HasMetrics:         analysis != nil,
		CouplingViolations: len(couplingViolations),
	}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
