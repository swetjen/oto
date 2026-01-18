package render

import (
	"fmt"
	"strings"

	"github.com/pacedotdev/oto/parser"
)

type typeRequirements struct {
	js    bool
	ts    bool
	swift bool
	dart  bool
}

type typeRef struct {
	Name string
	Type parser.FieldType
}

func requiredTypeMappings(template string) typeRequirements {
	return typeRequirements{
		js:    strings.Contains(template, ".JSType"),
		ts:    strings.Contains(template, ".TSType"),
		swift: strings.Contains(template, ".SwiftType"),
		dart:  strings.Contains(template, ".DartType"),
	}
}

func (r typeRequirements) any() bool {
	return r.js || r.ts || r.swift || r.dart
}

func collectTypeRefs(def parser.Definition) []typeRef {
	var refs []typeRef
	for _, service := range def.Services {
		for _, method := range service.Methods {
			refs = append(refs, typeRef{
				Name: fmt.Sprintf("%s.%s input", service.Name, method.Name),
				Type: method.InputObject,
			})
			refs = append(refs, typeRef{
				Name: fmt.Sprintf("%s.%s output", service.Name, method.Name),
				Type: method.OutputObject,
			})
		}
	}
	for _, obj := range def.Objects {
		for _, field := range obj.Fields {
			refs = append(refs, typeRef{
				Name: fmt.Sprintf("%s.%s", obj.Name, field.Name),
				Type: field.Type,
			})
		}
	}
	return refs
}

func validateTypeMappings(template string, def parser.Definition) error {
	req := requiredTypeMappings(template)
	if !req.any() {
		return nil
	}
	for _, ref := range collectTypeRefs(def) {
		if req.js && ref.Type.JSType == "" {
			return fmt.Errorf("missing js_type mapping for %s (go_type %q)", ref.Name, ref.Type.TypeName)
		}
		if req.ts && ref.Type.TSType == "" {
			return fmt.Errorf("missing ts_type mapping for %s (go_type %q)", ref.Name, ref.Type.TypeName)
		}
		if req.swift && ref.Type.SwiftType == "" {
			return fmt.Errorf("missing swift_type mapping for %s (go_type %q)", ref.Name, ref.Type.TypeName)
		}
		if req.dart && ref.Type.DartType == "" {
			return fmt.Errorf("missing dart_type mapping for %s (go_type %q)", ref.Name, ref.Type.TypeName)
		}
	}
	return nil
}
