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

		if e.data.Config.Debug {
			b, err := json.Marshal(e.data)
			if err == nil {
				writeFile("data.json", string(b))
			}
		}

		if in(Input, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.OutDir, "ent/input.go"), parseTemplate("ent/input", e.data))
		}

		if in(DB, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.OutDir, "db/db.go"), parseTemplate("ent/db", e.data))
		}

		if in(Query, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.OutDir, "ent/query.go"), parseTemplate("ent/query", e.data))
		}

		if in(Queries, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.OutDir, "handlers/queries.go"), parseTemplate("fiber/queries", e.data))
		}

		if in(RoutesQueries, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.OutDir, "routes/queries.go"), parseTemplate("fiber/routes-queries", e.data))
		}

		if in(RoutesMutations, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.OutDir, "routes/mutations.go"), parseTemplate("fiber/routes-mutations", e.data))
		}

		if in(Swagger, e.data.Config.Files) {
			writeFile(path.Join(root, e.data.Config.OutDir, "models/swagger.go"), parseTemplate("fiber/swagger", e.data))
		}

		if in(TsTypes, e.data.Config.Files) {
			for _, v := range e.data.Config.TsClientPath {

				writeFile(path.Join(root, v, "types.ts"), parseTemplate("typescript/types", e.data))
			}
		}

		if in(TsApi, e.data.Config.Files) {
			for _, v := range e.data.Config.TsClientPath {
				writeFile(path.Join(root, v, "api.ts"), parseTemplate("typescript/api", e.data))
			}
		}

		if in(TsRequest, e.data.Config.Files) {
			for _, v := range e.data.Config.TsClientPath {
				writeFile(path.Join(root, v, "request.ts"), parseTemplate("typescript/request", e.data))
			}
		}

		if e.data.Config.WithNestedMutations {
			if in(Services, e.data.Config.Files) {
				writeFile(path.Join(root, e.data.Config.OutDir, "services/mutations.go"), parseTemplate("ent/services", e.data))
			}

			if in(Mutations, e.data.Config.Files) {
				writeFile(path.Join(root, e.data.Config.OutDir, "handlers/mutations.go"), parseTemplate("fiber/services-mutations", e.data))
			}
		} else {
			if in(Mutations, e.data.Config.Files) {
				writeFile(path.Join(root, e.data.Config.OutDir, "handlers/mutations.go"), parseTemplate("fiber/mutations", e.data))
			}
		}

		if in(DartTypes, e.data.Config.Files) {
			for _, v := range e.data.Config.DartClientPath {
				writeFile(path.Join(root, v, "types.dart"), parseTemplate("dart/types", e.data))
			}
		}

		return next.Generate(g)
	})
}
