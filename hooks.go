package enput

import (
	"path"

	"entgo.io/ent/entc/gen"
)

func (e *Extension) generate(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		e.data.Graph = g
		root := rootDir()

		// backend
		if in(Input, e.data.Config.Files) {
			writeFile(path.Join(root, "ent/input.go"), parseTemplate("ent/input", e.data))
		}

		if in(Query, e.data.Config.Files) {
			writeFile(path.Join(root, "ent/query.go"), parseTemplate("ent/query", e.data))
		}

		if in(DartTypes, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.DartClientPath, "types.dart"), parseTemplate("dart/types", e.data))
		}

		if in(TsTypes, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.TsClientPath, "types.ts"), parseTemplate("typescript/types", e.data))
		}

		return next.Generate(g)
	})
}
