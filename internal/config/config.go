package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/mickamy/goimportmaps"
)

const (
	path = ".goimportmaps.yaml"
)

type Rule struct {
	Source  string   `yaml:"source"`
	Imports []string `yaml:"imports"`

	CompiledSource  *regexp.Regexp   `yaml:"-"`
	CompiledImports []*regexp.Regexp `yaml:"-"`
}

type Config struct {
	Forbidden []Rule `yaml:"forbidden"`
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
			return nil, fmt.Errorf("invalid from regex `%q: %w`", rule.Source, err)
		}
		for _, imprt := range rule.Imports {
			imprtRegexp, err := regexp.Compile(imprt)
			if err != nil {
				return nil, fmt.Errorf("invalid to regex `%q: %w`", rule.Imports, err)
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
func (c *Config) Validate(graph goimportmaps.Graph) []Violation {
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
