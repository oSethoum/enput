package enput

import (
	"embed"

	"entgo.io/ent/entc/gen"
)

//go:embed templates
var templates embed.FS

func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}

func NewExtension(opts ...option) *Extension {
	e := new(Extension)
	for _, opt := range opts {
		opt(e)
	}

	initFunctions(e)
	e.hooks = append(e.hooks, e.generate)

	return e
}
