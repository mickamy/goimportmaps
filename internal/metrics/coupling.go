package metrics

import (
	"github.com/mickamy/goimportmaps"
)

// CouplingMetrics represents coupling metrics for a package
type CouplingMetrics struct {
	AfferentCoupling int     // Ca - Number of packages that depend on this package
	EfferentCoupling int     // Ce - Number of packages this package depends on
	Instability      float64 // I = Ce / (Ca + Ce)
}

// PackageMetrics holds metrics for a single package
type PackageMetrics struct {
	Package  string
	Coupling CouplingMetrics
}

// CouplingAnalysis holds the complete coupling analysis results
type CouplingAnalysis struct {
	Packages map[string]CouplingMetrics
}

// CalculateCoupling analyzes the dependency graph and calculates coupling metrics
func CalculateCoupling(graph goimportmaps.Graph) *CouplingAnalysis {
	analysis := &CouplingAnalysis{
		Packages: make(map[string]CouplingMetrics),
	}

	// Initialize all packages
	allPackages := getAllPackages(graph)
	for _, pkg := range allPackages {
		analysis.Packages[pkg] = CouplingMetrics{}
	}

	// Calculate efferent coupling (Ce) - direct from graph
	for pkg, imports := range graph {
		metrics := analysis.Packages[pkg]
		metrics.EfferentCoupling = len(imports)
		analysis.Packages[pkg] = metrics
	}

	// Calculate afferent coupling (Ca) - reverse lookup
	for _, imports := range graph {
		for _, importedPkg := range imports {
			if metrics, exists := analysis.Packages[importedPkg]; exists {
				metrics.AfferentCoupling++
				analysis.Packages[importedPkg] = metrics
			}
		}
	}

	// Calculate instability (I)
	for pkg, metrics := range analysis.Packages {
		total := metrics.AfferentCoupling + metrics.EfferentCoupling
		if total > 0 {
			metrics.Instability = float64(metrics.EfferentCoupling) / float64(total)
		} else {
			metrics.Instability = 0.0
		}
		analysis.Packages[pkg] = metrics
	}

	return analysis
}

// getAllPackages extracts all unique packages from the graph
func getAllPackages(graph goimportmaps.Graph) []string {
	packageSet := make(map[string]bool)

	// add all source packages
	for pkg := range graph {
		packageSet[pkg] = true
	}

	// add all imported packages
	for _, imports := range graph {
		for _, importedPkg := range imports {
			packageSet[importedPkg] = true
		}
	}

	packages := make([]string, 0, len(packageSet))
	for pkg := range packageSet {
		packages = append(packages, pkg)
	}

	return packages
}

// GetHighCouplingPackages returns packages with coupling above thresholds
func (a *CouplingAnalysis) GetHighCouplingPackages(maxEfferent, maxAfferent int, maxInstability float64) []PackageMetrics {
	var highCoupling []PackageMetrics

	for pkg, metrics := range a.Packages {
		if metrics.EfferentCoupling > maxEfferent ||
			metrics.AfferentCoupling > maxAfferent ||
			metrics.Instability > maxInstability {
			highCoupling = append(highCoupling, PackageMetrics{
				Package:  pkg,
				Coupling: metrics,
			})
		}
	}

	return highCoupling
}
