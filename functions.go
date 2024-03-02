package enput

import (
	"strings"

	"entgo.io/ent/entc/gen"
	"github.com/iancoleman/strcase"
)

var (
	pascal     = gen.Funcs["pascal"].(func(string) string)
	kebab      = strcase.ToKebab
	snake      = gen.Funcs["snake"].(func(string) string)
	buggyCamel = gen.Funcs["camel"].(func(string) string)
	camel      = func(s string) string { return buggyCamel(snake(s)) }
	strCase    = snake
)

func (e *Extension) initFunctions() {

	if e.data.Config.Case == Camel {
		strCase = camel
	} else if e.data.Config.Case == Pascal {
		strCase = pascal
	}

	is_comparable := func(field Field) bool {
		return has_prefixes(field.Type, []string{"int", "uint", "string", "float", "time.Time"})
	}

	orderable := func(f Field) bool {
		return has_prefixes(f.Type, []string{
			"string",
			"int",
			"uint",
			"float",
			"time.Time",
			"bool",
		})
	}

	order_fields := func(s Schema) string {
		fields := []string{}
		for _, f := range s.Fields {
			if orderable(f) {
				fields = append(fields, snake(f.Name))
			}
		}
		return "\"" + strings.Join(fields, "\" | \"") + "\""
	}

	select_fields := func(s Schema) string {
		fields := []string{}
		for _, f := range s.Fields {
			fields = append(fields, snake(f.Name))
		}
		return "\"" + strings.Join(fields, "\" | \"") + "\""
	}

	// aggregate_fields := func(s Schema) string {
	// 	fields := []string{}
	// 	for _, f := range s.Fields {
	// 		if f.Type != "json.RawMessage" {
	// 			fields = append(fields, snake(f.Name))

	// 		}
	// 	}
	// 	return "\"" + strings.Join(fields, "\" | \"") + "\""
	// }

	// group_by_fields := func(s Schema) string {
	// 	fields := []string{}
	// 	for _, f := range s.Fields {
	// 		fields = append(fields, snake(f.Name))
	// 	}
	// 	return "\"" + strings.Join(fields, "\" | \"") + "\""
	// }

	gen.Funcs["order_fields"] = order_fields
	gen.Funcs["select_fields"] = select_fields
	// gen.Funcs["aggregate_fields"] = aggregate_fields
	// gen.Funcs["group_by_fields"] = group_by_fields
	gen.Funcs["is_comparable"] = is_comparable
	gen.Funcs["case"] = strCase
	gen.Funcs["pascal"] = pascal
	gen.Funcs["kebab"] = kebab
	gen.Funcs["camel"] = camel
}
