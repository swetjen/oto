package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// TypeMappingConfig describes the YAML configuration for type overrides.
type TypeMappingConfig struct {
	Version string             `yaml:"version"`
	Types   TypeMappingSection `yaml:"types"`
}

// TypeMappingSection holds override entries.
type TypeMappingSection struct {
	Overrides []TypeOverride `yaml:"overrides"`
}

// TypeOverride provides per-language overrides for a Go type.
type TypeOverride struct {
	GoType    string  `yaml:"go_type"`
	JSType    *string `yaml:"js_type,omitempty"`
	TSType    *string `yaml:"ts_type,omitempty"`
	SwiftType *string `yaml:"swift_type,omitempty"`
	DartType  *string `yaml:"dart_type,omitempty"`
}

// LoadTypeOverrides loads type overrides from a YAML file.
func LoadTypeOverrides(path string) ([]TypeOverride, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg TypeMappingConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, errors.Wrap(err, "parse type mapping yaml")
	}
	if cfg.Version != "" && cfg.Version != "1" {
		return nil, fmt.Errorf("unsupported type mapping version %q", cfg.Version)
	}
	overrides := make([]TypeOverride, 0, len(cfg.Types.Overrides))
	for i, override := range cfg.Types.Overrides {
		override.GoType = strings.TrimSpace(override.GoType)
		if override.GoType == "" {
			return nil, fmt.Errorf("type override %d missing go_type", i+1)
		}
		if override.JSType == nil && override.TSType == nil && override.SwiftType == nil && override.DartType == nil {
			return nil, fmt.Errorf("type override %d (%s) has no target types", i+1, override.GoType)
		}
		overrides = append(overrides, override)
	}
	return overrides, nil
}

func (o TypeOverride) matches(ftype FieldType) bool {
	if o.GoType == "" {
		return false
	}
	if strings.Contains(o.GoType, "/") {
		return o.GoType == ftype.TypeID
	}
	return o.GoType == ftype.TypeName ||
		o.GoType == ftype.CleanObjectName ||
		o.GoType == ftype.ObjectName
}

func (p *Parser) applyTypeOverrides(ftype *FieldType) {
	for _, override := range p.TypeOverrides {
		if !override.matches(*ftype) {
			continue
		}
		if override.JSType != nil {
			ftype.JSType = *override.JSType
		}
		if override.TSType != nil {
			ftype.TSType = *override.TSType
		}
		if override.SwiftType != nil {
			ftype.SwiftType = *override.SwiftType
		}
		if override.DartType != nil {
			ftype.DartType = *override.DartType
		}
	}
}
