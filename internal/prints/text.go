package prints

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/mickamy/goimportmaps"
	"github.com/mickamy/goimportmaps/internal/metrics"
	"github.com/mickamy/goimportmaps/internal/module"
)

func Text(w io.Writer, graph goimportmaps.Graph, modulePath string) {
	for from, toList := range graph {
		for _, to := range toList {
			from = module.Shorten(from, modulePath)
			to = module.Shorten(to, modulePath)
			_, _ = fmt.Fprintf(w, "  %s --> %s\n", from, to)
		}
	}
}

func TextWithMetrics(w io.Writer, graph goimportmaps.Graph, modulePath string, analysis *metrics.CouplingAnalysis, maxEfferent, maxAfferent int, maxInstability float64) {
	// Print dependency graph
	fmt.Fprintf(w, "📊 Dependency Graph:\n")
	for from, toList := range graph {
		for _, to := range toList {
			from = module.Shorten(from, modulePath)
			to = module.Shorten(to, modulePath)
			_, _ = fmt.Fprintf(w, "  %s --> %s\n", from, to)
		}
	}

	// Print coupling metrics
	fmt.Fprintf(w, "\n📊 Coupling Metrics:\n\n")
	fmt.Fprintf(w, "%-50s %4s %4s %6s %s\n", "Package", "Ca", "Ce", "I", "Status")
	fmt.Fprintf(w, "%-50s %4s %4s %6s %s\n", strings.Repeat("-", 50), "----", "----", "------", "------")

	// Sort packages for consistent output
	var packages []string
	for pkg := range analysis.Packages {
		packages = append(packages, pkg)
	}
	sort.Strings(packages)

	for _, pkg := range packages {
		metrics := analysis.Packages[pkg]
		shortPkg := module.Shorten(pkg, modulePath)
		status := getMetricsStatus(metrics, maxEfferent, maxAfferent, maxInstability)
		
		fmt.Fprintf(w, "%-50s %4d %4d %6.2f %s\n", 
			shortPkg, 
			metrics.AfferentCoupling, 
			metrics.EfferentCoupling, 
			metrics.Instability, 
			status)
	}

	// Print violations
	highCoupling := analysis.GetHighCouplingPackages(maxEfferent, maxAfferent, maxInstability)
	if len(highCoupling) > 0 {
		fmt.Fprintf(w, "\n🚨 Coupling Violations:\n")
		for _, pkg := range highCoupling {
			shortPkg := module.Shorten(pkg.Package, modulePath)
			reasons := getViolationReasons(pkg.Coupling, maxEfferent, maxAfferent, maxInstability)
			fmt.Fprintf(w, "- %s: %s\n", shortPkg, reasons)
		}
	}
}

func getMetricsStatus(metrics metrics.CouplingMetrics, maxEfferent, maxAfferent int, maxInstability float64) string {
	if metrics.EfferentCoupling > maxEfferent || 
	   metrics.AfferentCoupling > maxAfferent || 
	   metrics.Instability > maxInstability {
		return "🚨"
	}
	return "✅"
}

func getViolationReasons(metrics metrics.CouplingMetrics, maxEfferent, maxAfferent int, maxInstability float64) string {
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
	
	return strings.Join(reasons, ", ")
}
