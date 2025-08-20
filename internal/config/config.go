package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/mickamy/goimportmaps"
	"github.com/mickamy/goimportmaps/internal/module"
)

type Mode string

const (
	ModeForbidden Mode = "forbidden"
	ModeAllowed   Mode = "allowed"
)

func NewMode(s string) (Mode, error) {
	switch m := Mode(s); m {
	case ModeForbidden, ModeAllowed:
		return m, nil
	default:
		return "", fmt.Errorf("invalid mode: %s", s)
	}
}

const (
	path = ".goimportmaps.yaml"
)

type Rule struct {
	Source  string   `yaml:"source"`
	Imports []string `yaml:"imports"`
	Stdlib  *bool    `yaml:"stdlib,omitempty"`

	CompiledSource  *regexp.Regexp   `yaml:"-"`
	CompiledImports []*regexp.Regexp `yaml:"-"`
}

type CouplingThresholds struct {
	MaxEfferent     int     `yaml:"max_efferent"`
	MaxAfferent     int     `yaml:"max_afferent"`
	MaxInstability  float64 `yaml:"max_instability"`
	WarnEfferent    int     `yaml:"warn_efferent"`
	WarnAfferent    int     `yaml:"warn_afferent"`
	WarnInstability float64 `yaml:"warn_instability"`
}

type Metrics struct {
	Coupling CouplingThresholds `yaml:"coupling"`
	Enabled  bool               `yaml:"enabled"`
}

type Config struct {
	Forbidden []Rule  `yaml:"forbidden"`
	Allowed   []Rule  `yaml:"allowed"`
	Metrics   Metrics `yaml:"metrics"`
}

func Load() (*Config, error) {
	return LoadByPath(path)
}

func LoadByPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// return default config if file doesn't exist
		if os.IsNotExist(err) {
			return getDefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config format: %w", err)
	}

	for i := range cfg.Forbidden {
		rule := &cfg.Forbidden[i]

		if rule.CompiledSource, err = regexp.Compile(rule.Source); err != nil {
			return nil, fmt.Errorf("invalid source regex `%q: %w`", rule.Source, err)
		}
		for _, imprt := range rule.Imports {
			imprtRegexp, err := regexp.Compile(imprt)
			if err != nil {
				return nil, fmt.Errorf("invalid import pattern `%s`: %w", imprt, err)
			}
			rule.CompiledImports = append(rule.CompiledImports, imprtRegexp)
		}
	}

	for i := range cfg.Allowed {
		rule := &cfg.Allowed[i]

		if rule.CompiledSource, err = regexp.Compile(rule.Source); err != nil {
			return nil, fmt.Errorf("invalid source regex `%q: %w`", rule.Source, err)
		}
		for _, imprt := range rule.Imports {
			imprtRegexp, err := regexp.Compile(imprt)
			if err != nil {
				return nil, fmt.Errorf("invalid import pattern `%s`: %w", imprt, err)
			}
			rule.CompiledImports = append(rule.CompiledImports, imprtRegexp)
		}
	}

	// set default values for metrics if not specified
	if !cfg.Metrics.Enabled {
		cfg.Metrics.Enabled = true
	}
	if cfg.Metrics.Coupling.MaxEfferent == 0 {
		cfg.Metrics.Coupling.MaxEfferent = 10
	}
	if cfg.Metrics.Coupling.MaxAfferent == 0 {
		cfg.Metrics.Coupling.MaxAfferent = 15
	}
	if cfg.Metrics.Coupling.MaxInstability == 0 {
		cfg.Metrics.Coupling.MaxInstability = 0.8
	}
	if cfg.Metrics.Coupling.WarnEfferent == 0 {
		cfg.Metrics.Coupling.WarnEfferent = 7
	}
	if cfg.Metrics.Coupling.WarnAfferent == 0 {
		cfg.Metrics.Coupling.WarnAfferent = 10
	}
	if cfg.Metrics.Coupling.WarnInstability == 0 {
		cfg.Metrics.Coupling.WarnInstability = 0.6
	}

	return &cfg, nil
}

func getDefaultConfig() *Config {
	return &Config{
		Metrics: Metrics{
			Enabled: true,
			Coupling: CouplingThresholds{
				MaxEfferent:     10,
				MaxAfferent:     15,
				MaxInstability:  0.8,
				WarnEfferent:    7,
				WarnAfferent:    10,
				WarnInstability: 0.6,
			},
		},
	}
}

type Violation struct {
	Source  string
	Import  string
	Message string
}

// Validate checks the import graph against forbidden rules.
// It returns a slice of human-readable violation messages.
func (c *Config) Validate(graph goimportmaps.Graph, mode Mode, modulePath string) []Violation {
	switch mode {
	case ModeForbidden:
		return c.ValidateForbidden(graph, modulePath)
	case ModeAllowed:
		return c.ValidateAllowed(graph, modulePath)
	default:
		panic(fmt.Errorf("invalid mode %s", mode))
	}
}

func (c *Config) ValidateForbidden(graph goimportmaps.Graph, modulePath string) []Violation {
	var violations []Violation

	for _, rule := range c.Forbidden {
		for source, imports := range graph {
			if !rule.CompiledSource.MatchString(source) {
				continue
			}

			for _, imprt := range imports {
				for _, imprtRegexp := range rule.CompiledImports {
					if !imprtRegexp.MatchString(imprt) {
						continue
					}
					violations = append(violations, Violation{
						Source:  source,
						Import:  imprtRegexp.String(),
						Message: fmt.Sprintf("%s imports %s (matched rule: %s → %s)", module.Shorten(source, modulePath), module.Shorten(imprt, modulePath), rule.Source, imprtRegexp.String()),
					})
				}
			}
		}
	}

	return violations
}

func (c *Config) ValidateAllowed(graph goimportmaps.Graph, modulePath string) []Violation {
	var violations []Violation

	for source, imports := range graph {
		for _, imprt := range imports {
			matched := false

			for _, rule := range c.Allowed {
				if !rule.CompiledSource.MatchString(source) {
					continue
				}

				allowStdlib := true
				if rule.Stdlib != nil {
					allowStdlib = *rule.Stdlib
				}

				if module.IsStdlib(imprt) && allowStdlib {
					matched = true
					break
				}

				for _, imprtRegexp := range rule.CompiledImports {
					if imprtRegexp.MatchString(imprt) {
						matched = true
						break
					}
				}
				if matched {
					break
				}
			}

			if !matched {
				violations = append(violations, Violation{
					Source:  source,
					Import:  imprt,
					Message: fmt.Sprintf("%s imports %s, but no allowed rule matched", module.Shorten(source, modulePath), module.Shorten(imprt, modulePath)),
				})
			}
		}
	}

	return violations
}
