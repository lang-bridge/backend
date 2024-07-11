package postgres

import (
	"go.uber.org/fx"
	"platform/internal/repository/postgres/gen"
	"platform/pkg/db/dbtx"
)

var Module = fx.Module("repository/postgres",
	fx.Provide(NewDB),
	fx.Provide(func(db dbtx.DBTX) *gen.Queries {
		return gen.New(db)
	}),
	fx.Provide(
		NewTagsRepository,
		NewKeysRepository,
		NewTranslationsRepository,
	),
)
