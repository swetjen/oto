package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestLoadTypeOverrides(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "types.yaml")
	err := os.WriteFile(path, []byte(`version: "1"
types:
  overrides:
    - go_type: "string"
      ts_type: "Text"
`), 0o644)
	is.NoErr(err)

	overrides, err := LoadTypeOverrides(path)
	is.NoErr(err)
	is.Equal(len(overrides), 1)
	is.Equal(overrides[0].GoType, "string")
	is.True(overrides[0].TSType != nil)
	is.Equal(*overrides[0].TSType, "Text")
}

func TestTypeOverridesApply(t *testing.T) {
	is := is.New(t)
	tsType := "Text"
	p := New("./testdata/services/pleasantries")
	p.TypeOverrides = []TypeOverride{
		{
			GoType: "string",
			TSType: &tsType,
		},
	}
	def, err := p.Parse()
	is.NoErr(err)
	obj, err := def.Object("GreetRequest")
	is.NoErr(err)
	is.Equal(obj.Fields[0].Type.TSType, "Text")
	is.Equal(obj.Fields[0].Type.JSType, "string")
}
