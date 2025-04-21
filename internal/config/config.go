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
	From string `yaml:"from"`
	To   string `yaml:"to"`

	CompiledFrom *regexp.Regexp `yaml:"-"`
	CompiledTo   *regexp.Regexp `yaml:"-"`
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

		if rule.CompiledFrom, err = regexp.Compile(rule.From); err != nil {
			return nil, fmt.Errorf("invalid from regex `%q: %w`", rule.From, err)
		}
		if rule.CompiledTo, err = regexp.Compile(rule.To); err != nil {
			return nil, fmt.Errorf("invalid to regex `%q: %w`", rule.To, err)
		}
	}

	return &cfg, nil
}

type Violation struct {
	From    string
	To      string
	Message string
}

// Validate checks the import graph against forbidden rules.
// It returns a slice of human-readable violation messages.
func (c *Config) Validate(graph goimportmaps.Graph) []Violation {
	var violations []Violation

	for _, rule := range c.Forbidden {
		for from, imports := range graph {
			if !rule.CompiledFrom.MatchString(from) {
				continue
			}

			for _, to := range imports {
				if !rule.CompiledTo.MatchString(to) {
					continue
				}
				violations = append(violations, Violation{
					From:    from,
					To:      to,
					Message: fmt.Sprintf("%s imports %s (matched rule: %s → %s)", from, to, rule.From, rule.To),
				})
			}
		}
	}

	return violations
}
