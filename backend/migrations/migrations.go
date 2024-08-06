package migrations

import (
	"embed"
)

//nolint:revive
//go:embed *.sql
var Source embed.FS
