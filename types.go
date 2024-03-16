package enput

import (
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type Extension struct {
	entc.DefaultExtension
	hooks []gen.Hook
	data  data
}

type Driver = string
type Case = string
type File = uint

const (
	Input File = iota * 2
	DB
	Query
	DartTypes
	TsTypes
	TsApi
	Swagger
	Queries
	Mutations
	RoutesQueries
	RoutesMutations
	Services
)

const (
	Pascal Case = "pascal"
	Camel  Case = "camel"
	Snake  Case = "snake"
)

type option = func(*Extension)

type data struct {
	Config       *Config
	Package      string
	Schemas      []Schema
	Imports      []string
	InputImports []string
}

type Config struct {
	Case                Case
	IDType              string
	TsClientPath        string
	DartClientPath      string
	OutDir              string
	Debug               bool
	Package             string
	FormTag             bool
	WithSwaggerRename   bool
	WithSwagger         bool
	WithNestedMutations bool
	WithHooks           bool
	WithInterceptors    bool
	WithSecurity        bool
	Files               []File
	IgnoreSchemas       IgnoreSchemas
}

type IgnoreSchemas struct {
	Query  []string
	Create []string
	Update []string
	Delete []string
}

type comparable interface {
	~string | ~int | ~float32 | ~uint
}

func go_to_ts(k string) string {
	if strings.HasPrefix(k, "time.Time") {
		return "string"
	}
	if strings.HasPrefix(k, "bool") {
		return "boolean"
	}
	if strings.HasPrefix(k, "int") {
		return "number"
	}
	if strings.HasPrefix(k, "uint") {
		return "number"
	}
	if strings.HasPrefix(k, "float") {
		return "number"
	}
	if strings.HasPrefix(k, "string") {
		return "string"
	}
	return "any"
}

func go_to_dart(k string) string {
	if strings.HasPrefix(k, "time.Time") {
		return "String"
	}
	if strings.HasPrefix(k, "bool") {
		return "bool"
	}
	if strings.HasPrefix(k, "int") {
		return "int"
	}
	if strings.HasPrefix(k, "uint") {
		return "int"
	}
	if strings.HasPrefix(k, "float") {
		return "double"
	}
	if strings.HasPrefix(k, "string") {
		return "String"
	}
	return "Map<String, dynamic>"
}

type Schema struct {
	Name    string
	Imports []string
	Fields  []Field
	Edges   []Edge
}

type Field struct {
	Name       string
	RawType    any
	Type       string
	Tag        string
	Optional   bool
	Nillable   bool
	Slice      bool
	EdgeField  bool
	Enum       bool
	Enums      []string
	Types      FieldTypes
	HasDefault bool
	IsJson     bool
	Comment    string
	Sensitive  bool
}

type FieldTypes struct {
	Ts   string
	Dart string
}

type Edge struct {
	Name      string
	Type      string
	Unique    bool
	Optional  bool
	Field     string
	OwnerFK   bool
	EdgeField bool
	Tags      EdgeTags
	Through   *Edge
	Comment   string
}

type EdgeTags struct {
	UniqueTag    string
	AddIdsTag    string
	RemoveIdsTag string
	ClearTag     string
}
