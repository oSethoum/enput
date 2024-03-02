package enput

import (
	"encoding/json"
	"path"

	"entgo.io/ent/entc/gen"
)

func (e *Extension) generate(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		root := rootDir()
		e.parse(g)

		b, err := json.Marshal(e.data)
		if err == nil {
			writeFile("data.json", string(b))
		}

		// backend
		if in(Input, e.data.Config.Files) {
			writeFile(path.Join(root, "ent/input.go"), parseTemplate("ent/input", e.data))
		}

		if in(DB, e.data.Config.Files) {
			writeFile(path.Join(root, "db/db.go"), parseTemplate("ent/db", e.data))
		}

		if in(Query, e.data.Config.Files) {
			writeFile(path.Join(root, "ent/query.go"), parseTemplate("ent/query", e.data))
		}

		if in(Queries, e.data.Config.Files) {
			writeFile(path.Join(root, "handlers/queries.go"), parseTemplate("fiber/queries", e.data))
		}

		if in(RoutesQueries, e.data.Config.Files) {
			writeFile(path.Join(root, "routes/queries.go"), parseTemplate("fiber/routes-queries", e.data))
		}

		if in(RoutesMutations, e.data.Config.Files) {
			writeFile(path.Join(root, "routes/mutations.go"), parseTemplate("fiber/routes-mutations", e.data))
		}

		if in(Swagger, e.data.Config.Files) {
			writeFile(path.Join(root, "models/swagger.go"), parseTemplate("fiber/swagger", e.data))
		}

		if in(TsTypes, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.TsClientPath, "types.ts"), parseTemplate("typescript/types", e.data))
		}

		if in(TsApi, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.TsClientPath, "api.ts"), parseTemplate("typescript/api", e.data))
		}

		if e.data.Config.WithNestedMutations {
			if in(Services, e.data.Config.Files) {
				writeFile(path.Join(root, "services/mutations.go"), parseTemplate("ent/services", e.data))
			}

			if in(Mutations, e.data.Config.Files) {
				writeFile(path.Join(root, "handlers/mutations.go"), parseTemplate("fiber/services-mutations", e.data))
			}
		} else {
			if in(Mutations, e.data.Config.Files) {
				writeFile(path.Join(root, "handlers/mutations.go"), parseTemplate("fiber/mutations", e.data))
			}
		}

		if in(DartTypes, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.DartClientPath, "types.dart"), parseTemplate("dart/types", e.data))
		}

		return next.Generate(g)
	})
}
