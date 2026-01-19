package render

import "github.com/pacedotdev/oto/parser"

func authFieldFor(def parser.Definition, method parser.Method) *parser.Field {
	objectName := method.InputObject.CleanObjectName
	for i := range def.Objects {
		if def.Objects[i].Name != objectName {
			continue
		}
		for j := range def.Objects[i].Fields {
			field := &def.Objects[i].Fields[j]
			if field.Auth != nil {
				return field
			}
		}
	}
	return nil
}

func authSpecFor(def parser.Definition, method parser.Method) *parser.AuthSpec {
	field := authFieldFor(def, method)
	if field == nil {
		return nil
	}
	return field.Auth
}

func authSpecs(def parser.Definition) []parser.AuthSpec {
	seen := make(map[string]struct{})
	var specs []parser.AuthSpec
	for i := range def.Objects {
		for j := range def.Objects[i].Fields {
			field := def.Objects[i].Fields[j]
			if field.Auth == nil {
				continue
			}
			key := field.Auth.Scheme + "|" + field.Auth.In + "|" + field.Auth.Name + "|" + field.Auth.Prefix
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			specs = append(specs, *field.Auth)
		}
	}
	return specs
}

func nonAuthFields(object parser.Object) []parser.Field {
	fields := make([]parser.Field, 0, len(object.Fields))
	for _, field := range object.Fields {
		if field.Auth == nil {
			fields = append(fields, field)
		}
	}
	return fields
}

func authInQuery(def parser.Definition) bool {
	specs := authSpecs(def)
	for _, spec := range specs {
		if spec.In == "query" {
			return true
		}
	}
	return false
}
