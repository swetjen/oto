package parser

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// TypesConfig holds the configuration for type overrides.
type TypesConfig struct {
	Version   string         `yaml:"version"`
	Overrides []TypeOverride `yaml:"overrides"`
}

// TypeOverride defines a mapping from a Go type to target language types.
type TypeOverride struct {
	GoType    string `yaml:"go_type"`
	JSType    string `yaml:"js_type"`
	TSType    string `yaml:"ts_type"`
	SwiftType string `yaml:"swift_type"`
	DartType  string `yaml:"dart_type"`
}

// LoadTypesConfig loads a TypesConfig from the specified YAML file path.
func LoadTypesConfig(path string) (*TypesConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "read types config file")
	}
	var cfg TypesConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, errors.Wrap(err, "parse types config YAML")
	}
	return &cfg, nil
}

// BuildOverrideMap converts the Overrides slice into a map keyed by GoType
// for efficient lookup.
func (c *TypesConfig) BuildOverrideMap() map[string]TypeOverride {
	m := make(map[string]TypeOverride, len(c.Overrides))
	for _, override := range c.Overrides {
		m[override.GoType] = override
	}
	return m
}
