package keys

import (
	"context"
	"github.com/stretchr/testify/require"
	"platform/internal/translations/entity/key"
	"platform/internal/translations/entity/project"
	"platform/pkg/db/dbtx"
	"platform/test"
	"platform/test/ptesting"
	"platform/test/ptesting/generator"
	"strings"
	"testing"
)

func TestCreateKey(t *testing.T) {
	test.RunTest(t, func(rep key.KeysRepository, db dbtx.DBTX) {
		ptesting.ForAll(t)(func(t *testing.T, gen *ptesting.Gen) {
			prj := generator.NewProject(t, gen, db)
			keyName := strings.ToLower(gen.NextString(5, 10))

			createdKey, err := rep.CreateKey(context.Background(), key.CreateKeyParam{
				ProjectID: project.ID(prj.ID),
				Name:      keyName,
				Platforms: []key.Platform{key.PlatformWeb, key.PlatformIOS},
				Tags:      nil,
			})
			require.NoError(t, err)
			require.Equal(t, keyName, createdKey.Name)

			// Try to create the same key again
			_, err = rep.CreateKey(context.Background(), key.CreateKeyParam{
				ProjectID: project.ID(prj.ID),
				Name:      keyName,
				Platforms: []key.Platform{key.PlatformWeb, key.PlatformIOS},
				Tags:      nil,
			})
			require.ErrorIs(t, err, key.ErrKeyAlreadyExists)

			createdKey, err = rep.CreateKey(context.Background(), key.CreateKeyParam{
				ProjectID: project.ID(prj.ID),
				Name:      strings.ToUpper(keyName),
				Platforms: []key.Platform{key.PlatformWeb, key.PlatformIOS},
				Tags:      nil,
			})
			require.NoError(t, err)
			require.Equal(t, strings.ToUpper(keyName), createdKey.Name)
		})
	})
}
