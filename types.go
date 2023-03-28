package enput

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type extension struct {
	entc.DefaultExtension
	hooks []gen.Hook
}

type comparable interface{ ~string | ~int | ~float32 }
