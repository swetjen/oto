package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestLoadTypesConfig(t *testing.T) {
	is := is.New(t)

	// Create a temporary types.yaml file
	content := `version: "1"
overrides:
  - go_type: "time.Time"
    js_type: "Date"
    ts_type: "Date"
    swift_type: "Date"
    dart_type: "DateTime"
  - go_type: "github.com/google/uuid.UUID"
    js_type: "string"
    ts_type: "string"
    swift_type: "UUID"
    dart_type: "String"
  - go_type: "int64"
    swift_type: "Int64"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "types.yaml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	is.NoErr(err)

	// Load the config
	cfg, err := LoadTypesConfig(configPath)
	is.NoErr(err)

	is.Equal(cfg.Version, "1")
	is.Equal(len(cfg.Overrides), 3)

	// Check first override
	is.Equal(cfg.Overrides[0].GoType, "time.Time")
	is.Equal(cfg.Overrides[0].JSType, "Date")
	is.Equal(cfg.Overrides[0].TSType, "Date")
	is.Equal(cfg.Overrides[0].SwiftType, "Date")
	is.Equal(cfg.Overrides[0].DartType, "DateTime")

	// Check second override
	is.Equal(cfg.Overrides[1].GoType, "github.com/google/uuid.UUID")
	is.Equal(cfg.Overrides[1].JSType, "string")

	// Check partial override (only swift_type specified)
	is.Equal(cfg.Overrides[2].GoType, "int64")
	is.Equal(cfg.Overrides[2].SwiftType, "Int64")
	is.Equal(cfg.Overrides[2].JSType, "") // Not specified
}

func TestBuildOverrideMap(t *testing.T) {
	is := is.New(t)

	cfg := &TypesConfig{
		Version: "1",
		Overrides: []TypeOverride{
			{
				GoType:    "time.Time",
				JSType:    "Date",
				TSType:    "Date",
				SwiftType: "Date",
				DartType:  "DateTime",
			},
			{
				GoType:    "int64",
				SwiftType: "Int64",
			},
		},
	}

	m := cfg.BuildOverrideMap()

	is.Equal(len(m), 2)

	timeOverride, found := m["time.Time"]
	is.True(found)
	is.Equal(timeOverride.JSType, "Date")
	is.Equal(timeOverride.TSType, "Date")
	is.Equal(timeOverride.SwiftType, "Date")
	is.Equal(timeOverride.DartType, "DateTime")

	int64Override, found := m["int64"]
	is.True(found)
	is.Equal(int64Override.SwiftType, "Int64")
	is.Equal(int64Override.JSType, "") // Not specified, should be empty
}

func TestLoadTypesConfigFileNotFound(t *testing.T) {
	is := is.New(t)

	_, err := LoadTypesConfig("/nonexistent/path/types.yaml")
	is.True(err != nil)
}

func TestLoadTypesConfigInvalidYAML(t *testing.T) {
	is := is.New(t)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")
	err := os.WriteFile(configPath, []byte("not: valid: yaml: ["), 0644)
	is.NoErr(err)

	_, err = LoadTypesConfig(configPath)
	is.True(err != nil)
}

func TestParseWithTypeOverrides(t *testing.T) {
	is := is.New(t)
	patterns := []string{"./testdata/services/pleasantries"}
	parser := New(patterns...)
	parser.Verbose = testing.Verbose()
	parser.ExcludeInterfaces = []string{"Ignorer"}

	// Set up type overrides
	parser.TypeOverrides = map[string]TypeOverride{
		"string": {
			GoType:    "string",
			JSType:    "CustomString",
			TSType:    "CustomString",
			SwiftType: "CustomString",
			DartType:  "CustomString",
		},
		"int": {
			GoType:    "int",
			SwiftType: "Int64", // Only override Swift type
		},
	}

	def, err := parser.Parse()
	is.NoErr(err)

	// Find a string field and verify the override was applied
	welcomeInputObject, err := def.Object("WelcomeRequest")
	is.NoErr(err)

	// The "To" field is a string - should have overridden types
	toField := welcomeInputObject.Fields[0]
	is.Equal(toField.Name, "To")
	is.Equal(toField.Type.JSType, "CustomString")
	is.Equal(toField.Type.TSType, "CustomString")
	is.Equal(toField.Type.SwiftType, "CustomString")
	is.Equal(toField.Type.DartType, "CustomString")

	// The "Times" field is an int - should have only SwiftType overridden
	timesField := welcomeInputObject.Fields[2]
	is.Equal(timesField.Name, "Times")
	is.Equal(timesField.Type.JSType, "number")         // Not overridden, uses default
	is.Equal(timesField.Type.TSType, "number")         // Not overridden, uses default
	is.Equal(timesField.Type.SwiftType, "Int64")       // Overridden
	is.Equal(timesField.Type.DartType, "int")          // Not overridden, uses default
}
