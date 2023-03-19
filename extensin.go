package enput

import "entgo.io/ent/entc/gen"

func NewExtension() *extension {
	e := extension{}
	e.hooks = append(e.hooks, e.generate)
	return &extension{}
}

func (e *extension) Hooks() []gen.Hook {
	return e.hooks
}

func (e *extension) generate(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		// TODO: do your work here
		return next.Generate(g)
	})
}
