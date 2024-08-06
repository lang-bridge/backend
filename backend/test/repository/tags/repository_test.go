package tags

import (
	"context"
	"strings"
	"testing"

	"platform/internal/pkg/db/dbtx"

	"github.com/stretchr/testify/require"
	"platform/internal/translations/entity/key"
	"platform/internal/translations/entity/project"
	"platform/test"
	"platform/test/ptesting"
	"platform/test/ptesting/fixture"
)

func TestRepeatableInsert(t *testing.T) {
	test.RunTest(t, func(rep key.TagsRepository, db dbtx.DBTX) {
		ptesting.ForAll(t)(func(t *testing.T, gen *ptesting.Gen) {
			prj := fixture.NewProject(t, gen, db)
			newTags := ptesting.Array(gen, 3, func(gen *ptesting.Gen) string {
				return gen.NextString(3, 10)
			})

			tags, err := rep.EnsureTags(context.Background(), project.ID(prj.ID), newTags)
			require.NoError(t, err)
			require.Len(t, tags, 3)

			var upsertedTags []string
			for _, tag := range newTags {
				upsertedTags = append(upsertedTags, strings.ToLower(tag))
			}
			upsertedTags = append(upsertedTags, gen.NextString(5, 6))

			tags, err = rep.EnsureTags(context.Background(), 1, upsertedTags)
			require.NoError(t, err)
			require.Len(t, tags, 4)
		})
	})
}
