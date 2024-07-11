package fixture

import (
	"context"
	"github.com/stretchr/testify/require"
	"log/slog"
	"platform/internal/repository/postgres/gen"
	"platform/pkg/ctxlog"
	"platform/pkg/db/dbtx"
	"platform/test/ptesting"
	"testing"
)

const createProject = `
INSERT INTO projects (name)
VALUES ($1)
RETURNING id
`

func NewProject(t *testing.T, g *ptesting.Gen, db dbtx.DBTX) gen.Project {
	name := g.NextString(3, 10)
	var id int64
	err := db.QueryRowxContext(ctxlog.WithLogger(context.Background(), slog.Default()), createProject, name).Scan(&id)
	require.NoError(t, err)
	return gen.Project{
		ID:   id,
		Name: name,
	}
}
