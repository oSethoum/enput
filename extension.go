package enput

import (
	"embed"

	"entgo.io/ent/entc/gen"
)

//go:embed templates
var templates embed.FS

func (e *extension) Hooks() []gen.Hook {
	return e.hooks
}

func NewExtension() *extension {
	e := new(extension)
	e.hooks = append(e.hooks, e.debug, e.generate)
	return e
}
