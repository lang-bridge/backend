package translations

import (
	"context"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"platform/internal/translations/entity/key"
	"platform/internal/translations/entity/project"
	"platform/internal/translations/entity/translation"
	"platform/pkg/db/dbtx"
	"platform/test"
	"platform/test/ptesting"
	"platform/test/ptesting/fixture"
	"testing"
)

func TestUpsertTranslations(t *testing.T) {
	test.RunTest(t, func(rep translation.Repository, keysRepo key.KeysRepository, db dbtx.DBTX) {
		ptesting.ForAll(t)(func(t *testing.T, gen *ptesting.Gen) {
			prj := fixture.NewProject(t, gen, db)

			nextKey := gen.NextKey(project.ID(prj.ID))
			nextKey, err := keysRepo.CreateKey(context.Background(), key.CreateKeyParam{
				ProjectID: project.ID(prj.ID),
				Name:      nextKey.Name,
				Platforms: nextKey.Platforms,
				Tags:      nextKey.Tags,
			})
			require.NoError(t, err)

			langs := ptesting.Elems(gen, 1, 5,
				language.English,
				language.Russian,
				language.French,
				language.German,
				language.Spanish,
				language.Italian,
				language.Chinese,
			)
			translations := make([]translation.Value, len(langs))
			for i, lang := range langs {
				translations[i] = gen.NextTranslation(nextKey.ID)
				translations[i].Language = lang
			}

			err = rep.UpsertTranslations(context.Background(), translations)
			require.NoError(t, err)

			err = rep.UpsertTranslations(context.Background(), translations)
			require.NoError(t, err)
		})
	})
}
