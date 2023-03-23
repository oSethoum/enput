package enput

import (
	"fmt"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
)

var (
	snake  = gen.Funcs["snake"].(func(string) string)
	_camel = gen.Funcs["camel"].(func(string) string)
	camel  = func(s string) string { return _camel(snake(s)) }
)

func init() {
	gen.Funcs["tag"] = tag
	gen.Funcs["imports"] = imports
	gen.Funcs["null_field_create"] = null_field_create
	gen.Funcs["null_field_update"] = null_field_update
	gen.Funcs["extract_type"] = extract_type
	gen.Funcs["edge_field"] = edge_field
}

func tag(f *load.Field) string {
	if f.Tag == "" {
		name := camel(f.Name)
		if strings.HasSuffix(name, "ID") {
			name = strings.TrimSuffix(name, "ID")
			name += "Id"
		}
		return fmt.Sprintf("json:\"%s\"", name)
	}
	return f.Tag
}

func imports(g *gen.Graph) []string {
	imps := []string{}
	for _, s := range g.Schemas {
		for _, f := range s.Fields {
			if f.Info != nil && len(f.Info.PkgPath) != 0 {
				if !in(f.Info.PkgPath, imps) {
					imps = append(imps, f.Info.PkgPath)
				}
			}
		}
	}
	return imps
}

func null_field_create(f *load.Field) bool {
	return f.Optional || f.Default
}

func null_field_update(field *load.Field) bool {
	return !strings.HasPrefix(extract_type(field), "[]")
}

func extract_type(field *load.Field) string {
	if field.Info.Ident != "" {
		return field.Info.Ident
	}
	return field.Info.Type.String()
}

func edge_field(e *load.Edge) bool {
	return e.Field != ""
}