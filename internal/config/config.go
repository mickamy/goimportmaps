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
		return "", fmt.Errorf("invalid mode %s", s)
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

type Config struct {
	Forbidden []Rule `yaml:"forbidden"`
	Allowed   []Rule `yaml:"allowed"`
}

func Load() (*Config, error) {
	return LoadByPath(path)
}

func LoadByPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
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

	return &cfg, nil
}

type Violation struct {
	Source  string
	Import  string
	Message string
}

// Validate checks the import graph against forbidden rules.
// It returns a slice of human-readable violation messages.
func (c *Config) Validate(graph goimportmaps.Graph, mode Mode) []Violation {
	switch mode {
	case ModeForbidden:
		return c.ValidateForbidden(graph)
	case ModeAllowed:
		return c.ValidateAllowed(graph)
	default:
		panic(fmt.Errorf("invalid mode %s", mode))
	}
}

func (c *Config) ValidateForbidden(graph goimportmaps.Graph) []Violation {
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
						Message: fmt.Sprintf("%s imports %s (matched rule: %s → %s)", source, imprt, rule.Source, imprtRegexp.String()),
					})
				}
			}
		}
	}

	return violations
}

func (c *Config) ValidateAllowed(graph goimportmaps.Graph) []Violation {
	var violations []Violation

	for source, imports := range graph {
		allowed := false

		for _, rule := range c.Allowed {
			if !rule.CompiledSource.MatchString(source) {
				continue
			}

			allowStdlib := true
			if rule.Stdlib != nil {
				allowStdlib = *rule.Stdlib
			}

			for _, imprt := range imports {
				if module.IsStdlib(imprt) && allowStdlib {
					continue
				}

				matched := false
				for _, imprtRegexp := range rule.CompiledImports {
					if imprtRegexp.MatchString(imprt) {
						matched = true
						break
					}
				}
				if matched {
					allowed = true
				} else {
					violations = append(violations, Violation{
						Source:  source,
						Import:  imprt,
						Message: fmt.Sprintf("%s imports %s, but no allowed rule matched", source, imprt),
					})
				}
			}
		}

		if !allowed && len(imports) > 0 {
			for _, imprt := range imports {
				violations = append(violations, Violation{
					Source:  source,
					Import:  imprt,
					Message: fmt.Sprintf("%s imports %s, but no allowed rule matched", source, imprt),
				})
			}
		}
	}

	return violations
}
