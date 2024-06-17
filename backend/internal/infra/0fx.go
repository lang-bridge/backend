package infra

import (
	"go.uber.org/fx"
	"platform/pkg/db/tx"
)

var Module = fx.Module("infra",
	fx.Provide(NewLogger),
	fx.Provide(NewDB),
	fx.Provide(tx.NewManager),
)
