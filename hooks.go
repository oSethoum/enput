package enput

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"entgo.io/ent/entc/gen"
)

func (e *extension) generate(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		s := parseTemplate("input", g)
		err := os.WriteFile(path.Join(g.Target, "enput_input.go"), []byte(s), 0666)
		if err != nil {
			log.Fatalln(err)
		}

		s = parseTemplate("query", g)
		err = os.WriteFile(path.Join(g.Target, "enput_query.go"), []byte(s), 0666)
		if err != nil {
			log.Fatalln(err)
		}
		return next.Generate(g)
	})
}

func (e *extension) debug(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		b, _ := json.Marshal(g.Schemas)
		os.WriteFile("debug.json", b, 0666)
		return next.Generate(g)
	})
}
