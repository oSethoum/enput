package enput

import (
	"fmt"
	"path"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	"entgo.io/ent/schema/field"
	"github.com/iancoleman/strcase"
)

var (
	pascal     = gen.Funcs["pascal"].(func(string) string)
	kebab      = strcase.ToKebab
	snake      = gen.Funcs["snake"].(func(string) string)
	buggyCamel = gen.Funcs["camel"].(func(string) string)
	camel      = func(s string) string { return buggyCamel(snake(s)) }
)

func initFunctions(e *Extension) {
	is_base := func(f *load.Field) bool {
		return f.Name == "id" || f.Name == "created_at" || f.Name == "updated_at" || f.Immutable ||
			f.Name == "ID" || f.Name == "createdAt" || f.Name == "updatedAt" || f.Name == "CreatedAt" || f.Name == "UpdatedAt"
	}

	caseFunc := func(name string) string {
		if e.data.Config.Case == Pascal {
			return pascal(name)
		} else if e.data.Config.Case == Camel {

			return camel(name)
		} else {
			return snake(name)
		}
	}

	fieldTag := func(f *load.Field, omitempty ...bool) string {
		if len(omitempty) > 0 && omitempty[0] {
			f.Tag = fmt.Sprintf(`json:"%s,omitempty"`, caseFunc(f.Name))

			if e.data.Config.FormTag {
				return f.Tag + fmt.Sprintf(` form:"%s,omitempty"`, caseFunc(f.Name))
			}
			return f.Tag

		} else {

			f.Tag = fmt.Sprintf(`json:"%s"`, caseFunc(f.Name))

			if e.data.Config.FormTag {
				return f.Tag + fmt.Sprintf(` form:"%s"`, caseFunc(f.Name))
			}

			return f.Tag
		}
	}

	addIdsTag := func(schema string) string {
		tag := fmt.Sprintf(`json:"%s"`, caseFunc(fmt.Sprintf("add%sIds", pascal(schema))))

		if e.data.Config.FormTag {
			tag += fmt.Sprintf(` form:"%s"`, caseFunc(fmt.Sprintf("add%sIds", pascal(schema))))
		}

		return tag
	}

	removeIdsTag := func(schema string) string {
		tag := fmt.Sprintf(`json:"%s"`, caseFunc(fmt.Sprintf("remove%sIds", pascal(schema))))

		if e.data.Config.FormTag {
			tag += fmt.Sprintf(` form:"%s"`, caseFunc(fmt.Sprintf("remove%sIds", pascal(schema))))
		}

		return tag
	}

	is_edge_field := func(s load.Schema, f load.Field) bool {

		for _, e := range s.Edges {
			if f.Name == e.Field {
				return true
			}

		}

		return false
	}

	clearTag := func(schema string) string {
		tag := fmt.Sprintf(`json:"%s"`, caseFunc(fmt.Sprintf("clear%s", pascal(schema))))

		if e.data.Config.FormTag {
			tag += fmt.Sprintf(` form:"%s"`, caseFunc(fmt.Sprintf("clear%s", pascal(schema))))
		}

		return tag
	}

	uniqueEdgeTag := func(schema string) string {
		return fmt.Sprintf(`json:"%s"`, caseFunc(fmt.Sprintf("%sId", pascal(schema))))
	}

	is_json := func(f *load.Field) bool {
		return f.Info.Type.String() == "json.RawMessage"
	}

	imports_without_json := func(g *gen.Graph, isInput ...bool) []string {
		imps := []string{}
		for _, s := range g.Schemas {
			for _, f := range s.Fields {
				if f.Info.Type.String() == "json.RawMessage" {
					continue
				}
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

		// remove duplication
		result := []string{}

		for _, i := range imps {
			found := false
			for _, r := range result {
				if r == i {
					found = true
				}
			}
			if !found {
				result = append(result, i)
			}
		}

		return result
	}

	imports := func(g *gen.Graph, isInput ...bool) []string {
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

		// remove duplication
		result := []string{}

		for _, i := range imps {
			found := false
			for _, r := range result {
				if r == i {
					found = true
				}
			}
			if !found {
				result = append(result, i)
			}
		}

		return result
	}

	extract_type_info := func(t *field.TypeInfo) string {
		if t.Ident != "" {
			return t.Ident
		}
		return t.Type.String()
	}

	extract_type := func(f *load.Field) string {
		return extract_type_info(f.Info)
	}

	null_field_create := func(f *load.Field) bool {
		return f.Optional || f.Default
	}

	null_field_update := func(f *load.Field) bool {
		return !strings.HasPrefix(extract_type(f), "[]")
	}

	edge_field := func(e *load.Edge) bool {
		return e.Field != ""
	}

	is_comparable := func(f *load.Field) bool {
		return has_prefixes(extract_type(f), []string{
			"string",
			"int",
			"uint",
			"float",
			"time.Time",
		})
	}

	enum_or_edge_filed := func(s *load.Schema, f *load.Field) bool {
		for _, e := range s.Edges {
			if e.Field == f.Name {
				return true
			}
		}

		return extract_type(f) == "enum"

	}

	get_name := func(f *load.Field) string {
		n := camel(f.Name)

		if e.data.Config.Case == Snake {
			n = snake(f.Name)
		}

		if e.data.Config.Case == Pascal {
			n = camel(f.Name)
		}

		if strings.HasSuffix(n, "ID") {
			n = strings.TrimSuffix(n, "ID") + "Id"
		}
		return n
	}

	get_type_info := func(f *field.TypeInfo) string {
		s := extract_type_info(f)
		t := "any"
		slice := false
		if strings.HasPrefix(s, "[]") {
			slice = true
			s = strings.TrimPrefix(s, "[]")
		}
		for k, v := range go_ts {
			if strings.HasPrefix(s, k) {
				t = v
				break
			}
		}

		if slice {
			return t + "[]"
		}
		return t
	}

	get_type := func(f *load.Field) string {
		if len(f.Enums) > 0 {
			enums := []string{}
			for _, v := range f.Enums {
				enums = append(enums, "\""+v.V+"\"")
			}
			return strings.Join(enums, " | ")
		} else {
			s := extract_type(f)

			t := "any"
			slice := false
			if strings.HasPrefix(s, "[]") {
				slice = true
				s = strings.TrimPrefix(s, "[]")
			}
			for k, v := range go_ts {
				if strings.HasPrefix(s, k) {
					t = v
					break
				}
			}

			if slice {
				return t + "[]"
			}
			return t
		}
	}

	is_slice := func(f *load.Field) bool {
		return strings.HasSuffix(get_type(f), "[]")
	}

	id_type := func(s *load.Schema) string {
		for _, f := range s.Fields {
			if strings.ToLower(f.Name) == "id" {
				return get_type(f)
			}
		}
		return "number"
	}

	orderable := func(f *load.Field) bool {

		if len(f.Enums) > 0 {
			return true
		}

		return has_prefixes(extract_type(f), []string{
			"string",
			"int",
			"uint",
			"float",
			"time.Time",
			"bool",
		})
	}

	order_fields := func(s *load.Schema) string {
		fields := []string{}
		for _, f := range s.Fields {
			if orderable(f) {
				fields = append(fields, snake(get_name(f)))
			}
		}
		return "\"" + strings.Join(fields, "\" | \"") + "\""
	}

	select_fields := func(s *load.Schema) string {
		fields := []string{}
		for _, f := range s.Fields {
			fields = append(fields, snake(get_name(f)))
		}
		return "\"" + strings.Join(fields, "\" | \"") + "\""
	}

	aggregate_fields := func(s *load.Schema) string {
		fields := []string{}
		for _, f := range s.Fields {
			if f.Info.Type.String() != "json.RawMessage" {
				fields = append(fields, snake(get_name(f)))

			}
		}
		return "\"" + strings.Join(fields, "\" | \"") + "\""
	}

	group_by_fields := func(s *load.Schema) string {
		fields := []string{}
		for _, f := range s.Fields {
			fields = append(fields, snake(get_name(f)))
		}
		return "\"" + strings.Join(fields, "\" | \"") + "\""
	}

	gen.Funcs["case"] = caseFunc
	gen.Funcs["camel"] = camel
	gen.Funcs["tag"] = fieldTag
	gen.Funcs["imports_without_json"] = imports_without_json
	gen.Funcs["imports"] = imports
	gen.Funcs["is_json"] = is_json
	gen.Funcs["null_field_create"] = null_field_create
	gen.Funcs["null_field_update"] = null_field_update
	gen.Funcs["extract_type"] = extract_type
	gen.Funcs["extract_type_info"] = extract_type_info
	gen.Funcs["edge_field"] = edge_field
	gen.Funcs["is_comparable"] = is_comparable
	gen.Funcs["enum_or_edge_filed"] = enum_or_edge_filed
	gen.Funcs["get_name"] = get_name
	gen.Funcs["get_type"] = get_type
	gen.Funcs["get_type_info"] = get_type_info
	gen.Funcs["is_slice"] = is_slice
	gen.Funcs["id_type"] = id_type
	gen.Funcs["order_fields"] = order_fields
	gen.Funcs["select_fields"] = select_fields
	gen.Funcs["aggregate_fields"] = aggregate_fields
	gen.Funcs["group_by_fields"] = group_by_fields
	gen.Funcs["dir"] = path.Dir
	gen.Funcs["is_base"] = is_base
	gen.Funcs["is_edge_field"] = is_edge_field
	gen.Funcs["kebab"] = kebab
	gen.Funcs["add_ids_tag"] = addIdsTag
	gen.Funcs["remove_ids_tag"] = removeIdsTag
	gen.Funcs["clear_tag"] = clearTag
	gen.Funcs["unique_edge_tag"] = uniqueEdgeTag
}
