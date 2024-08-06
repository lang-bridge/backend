package keys

import "go.uber.org/fx"

var Module = fx.Module("keys", fx.Provide(NewService))
