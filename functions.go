package enput

import (
	"fmt"
	"path"
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
	gen.Funcs["comparable"] = comparable
	gen.Funcs["enum_or_edge_filed"] = enum_or_edge_filed
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

func imports(g *gen.Graph, isInput ...bool) []string {
	imps := []string{}

	for _, s := range g.Schemas {
		for _, f := range s.Fields {
			if len(f.Enums) > 0 && len(isInput) > 0 && isInput[0] {
				imps = append(imps, path.Join(g.Package, strings.Split(f.Info.Ident, ".")[0]))
			}
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

func comparable(f *load.Field) bool {
	return has_prefixes(extract_type(f), []string{
		"string",
		"int",
		"uint",
		"float",
		"time.Time",
	})
}

func enum_or_edge_filed(s *load.Schema, f *load.Field) bool {
	for _, e := range s.Edges {
		if e.Field == f.Name {
			return extract_type(f) == "enum"
		}
	}
	return false
}
