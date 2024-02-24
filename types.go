package entify

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
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
	Query
	DartTypes
	TsTypes
)

const (
	Pascal Case = "pascal"
	Camel  Case = "camel"
	Snake  Case = "snake"
)

type option = func(*Extension)

type data struct {
	*gen.Graph
	Config        *Config
	CurrentSchema *load.Schema
}

type Config struct {
	Case           Case
	IDType         string
	TsClientPath   string
	DartClientPath string
	Files          []File
	Package        string
	IgnoreSchemas  []string
	FormTag        bool
}

type comparable interface {
	~string | ~int | ~float32 | ~uint
}

var go_ts = map[string]string{
	"time.Time": "string",
	"bool":      "boolean",
	"int":       "number",
	"uint":      "number",
	"float":     "number",
	"enum":      "string",
	"string":    "string",
	"any":       "any",
	"other":     "any",
	"json":      "any",
}
