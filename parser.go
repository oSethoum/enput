package enput

import (
	"fmt"
	"strings"

	"entgo.io/ent/entc/gen"
)

func (ex *Extension) parse(g *gen.Graph) {
	ex.data.Package = g.Package
	schemas := []Schema{}
	imports := []string{}

	for _, n := range g.Nodes {
		schema := Schema{
			Name: n.Name,
		}

		nFields := append([]*gen.Field{n.ID}, n.Fields...)

		for _, f := range nFields {
			field := Field{
				Name:       f.Name,
				Type:       "any",
				Nillable:   f.Nillable,
				Optional:   f.Optional,
				HasDefault: f.Default,
				EdgeField:  f.IsEdgeField(),
				IsJson:     f.IsBool(),
				Comment:    f.Comment(),
				Enum:       f.IsEnum(),
				Sensitive:  f.Sensitive(),
			}

			// Slice
			field.Slice = strings.HasPrefix(field.Type, "[]")

			// Types
			if f.Type != nil && !f.IsJSON() {
				if f.Type.Ident != "" {
					field.Type = f.Type.Ident
				} else {
					field.Type = f.Type.String()
				}
			}

			// Types.Ts
			if field.Enum {
				field.Types.Ts = "string"
			} else {
				t := strings.ReplaceAll(field.Type, " ", "")
				if field.Slice {
					t = strings.ReplaceAll(t, "[]", "")
				}
				t = go_to_ts(t)
				if field.Slice {
					t += "[]"
				}
				field.Types.Ts = t
			}

			// Types.Dart
			if field.Enum {
				field.Types.Dart = "String"
			} else {
				t := strings.ReplaceAll(field.Type, " ", "")
				if field.Slice {
					t = strings.ReplaceAll(t, "[]", "")
				}
				t = go_to_dart(t)
				if field.Slice {
					t += fmt.Sprintf("List<%s>", t)
				}
				field.Types.Dart = t
			}

			// Enum
			if field.Enum {
				for _, enum := range f.Enums {
					field.Enums = append(field.Enums, enum.Value)
				}

				pkg := ex.data.Package + "/" + strings.Split(field.Type, ".")[0]
				if !in(pkg, ex.data.InputImports) {
					ex.data.InputImports = append(ex.data.InputImports, pkg)
				}

				field.Types.Ts = fmt.Sprintf(`"%s"`, strings.Join(field.Enums, `" | "`))
			}

			// Tags
			if f.StructTag != "" {
				field.Tag = f.StructTag
			} else {
				name := strCase(field.Name)
				if ex.data.Config.FormTag {
					field.Tag = fmt.Sprintf(`json:"%s" form:"%s"`, name, name)
				} else {
					field.Tag = fmt.Sprintf(`json:"%s"`, name)
				}
			}
			if !in(f.Type.PkgPath, schema.Imports) && f.Type.PkgPath != "" {
				schema.Imports = append(schema.Imports, f.Type.PkgPath)
			}

			schema.Fields = append(schema.Fields, field)
		}

		for _, e := range n.Edges {
			edge := Edge{
				Name:     e.Name,
				Type:     e.Type.Name,
				Unique:   e.Unique,
				Optional: e.Optional,
				OwnerFK:  e.OwnFK(),

				Comment: e.Comment(),
			}

			// Field
			if e.Field() != nil {
				edge.EdgeField = true
				edge.Field = e.Field().Name
			} else if e.Ref != nil && e.Ref.Field() != nil {
				edge.Field = e.Ref.Field().Name
			}

			// Tags
			if ex.data.Config.FormTag {
				unique := strCase(pascal(e.Name) + "ID")
				edge.Tags.UniqueTag = fmt.Sprintf(`json:"%s" form:"%s"`, unique, unique)
				add := strCase("Add" + pascal(e.Name) + "IDs")
				edge.Tags.AddIdsTag = fmt.Sprintf(`json:"%s" form:"%s"`, add, add)
				remove := strCase("Remove" + pascal(e.Name) + "IDs")
				edge.Tags.RemoveIdsTag = fmt.Sprintf(`json:"%s" form:"%s"`, remove, remove)
				clear := strCase("Clear" + pascal(e.Name) + "IDs")
				edge.Tags.ClearTag = fmt.Sprintf(`json:"%s" form:"%s"`, clear, clear)
			} else {
				unique := strCase(pascal(e.Name) + "ID")
				edge.Tags.UniqueTag = fmt.Sprintf(`json:"%s"`, unique)
				add := strCase("Add" + pascal(e.Name) + "IDs")
				edge.Tags.AddIdsTag = fmt.Sprintf(`json:"%s"`, add)
				remove := strCase("Remove" + pascal(e.Name) + "IDs")
				edge.Tags.RemoveIdsTag = fmt.Sprintf(`json:"%s"`, remove)
				clear := strCase("Clear" + pascal(e.Name))
				edge.Tags.ClearTag = fmt.Sprintf(`json:"%s"`, clear)
			}
			schema.Edges = append(schema.Edges, edge)
		}

		for _, im := range schema.Imports {
			if !in(im, imports) {
				imports = append(imports, im)
			}
		}

		schemas = append(schemas, schema)
	}

	ex.data.Imports = imports
	ex.data.Schemas = schemas
}
