package postgres

import (
	"go.uber.org/fx"
	"platform/internal/pkg/db/dbtx"
	"platform/internal/repository/postgres/gen"
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
