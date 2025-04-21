package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/mickamy/goimportmaps"
)

const (
	path = ".goimportmaps.yaml"
)

type Rule struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type Config struct {
	Forbidden []Rule `yaml:"forbidden"`
}

func Load() (Config, error) {
	return LoadByPath(path)
}

func LoadByPath(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("invalid config format: %w", err)
	}

	return cfg, nil
}

// Validate checks the import graph against forbidden rules.
// It returns a slice of human-readable violation messages.
func (c Config) Validate(graph goimportmaps.Graph) []string {
	var violations []string

	for _, rule := range c.Forbidden {
		imports, ok := graph[rule.From]
		if !ok {
			continue
		}
		for _, to := range imports {
			if to == rule.To {
				violations = append(violations,
					fmt.Sprintf("🚨 Violation: %s imports %s", rule.From, rule.To),
				)
			}
		}
	}

	return violations
}
