package keys

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"platform/internal/translations/entity/project"
	"platform/internal/translations/entity/translation"
	"platform/internal/translations/service/keys"
	"platform/pkg/db/dbtx"
	"platform/test"
	"platform/test/ptesting"
	"platform/test/ptesting/fixture"
)

func TestCreateKey(t *testing.T) {
	test.RunTest(t, func(s chi.Router, db dbtx.DBTX) {
		srv := httptest.NewServer(s)
		t.Cleanup(func() {
			srv.Close()
		})

		ptesting.ForAll(t)(func(t *testing.T, gen *ptesting.Gen) {
			prj := fixture.NewProject(t, gen, db)

			nextKey := gen.NextKey(project.ID(prj.ID))
			translations := ptesting.Array(gen, 5, func(gen *ptesting.Gen) translation.Value {
				return gen.NextTranslation(nextKey.ID)
			})
			nextTag := gen.NextString(3, 10)

			var requestTranslation []keys.Translate
			langs := map[string]any{}
			for _, tr := range translations {
				if _, ok := langs[tr.Language.String()]; ok {
					continue
				}
				langs[tr.Language.String()] = ""
				requestTranslation = append(requestTranslation, keys.Translate{
					Language: tr.Language,
					Value:    tr.Translation,
				})
			}

			data := map[string]interface{}{
				"name":        nextKey.Name,
				"platforms":   nextKey.Platforms,
				"existedTags": nextKey.Tags,
				"newTags":     []string{nextTag},
				"translates":  requestTranslation,
			}
			jsonReq, err := json.Marshal(data)
			require.NoError(t, err)

			req, err := http.NewRequest("POST",
				fmt.Sprintf("%s/api/v1/projects/%d/keys", srv.URL, prj.ID), bytes.NewReader(jsonReq))
			require.NoError(t, err)

			resp, err := srv.Client().Do(req)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)
		})
	})
	sentry.Flush(time.Minute)
}

func TestCreateKey_BadRequest(t *testing.T) {
	test.RunTest(t, func(s chi.Router, db dbtx.DBTX) {
		srv := httptest.NewServer(s)
		t.Cleanup(func() {
			srv.Close()
		})

		ptesting.ForAll(t)(func(t *testing.T, gen *ptesting.Gen) {
			prj := fixture.NewProject(t, gen, db)

			key := gen.NextKey(project.ID(prj.ID))
			tr := gen.NextTranslation(key.ID)

			data := map[string]interface{}{
				"name":        key.Name,
				"platforms":   key.Platforms,
				"existedTags": key.Tags,
				"newTags":     []string{gen.NextString(3, 5)},
				"translates": []keys.Translate{
					{
						Language: tr.Language,
						Value:    tr.Translation,
					},
				},
			}

			nextInt := gen.NextInt(0, 2)
			switch nextInt {
			case 0:
				delete(data, "name")
			case 1:
				delete(data, "platforms")
			case 2:
				delete(data, "translates")
			}

			jsonReq, err := json.Marshal(data)
			require.NoError(t, err)

			req, err := http.NewRequest("POST",
				fmt.Sprintf("%s/api/v1/projects/%d/keys", srv.URL, prj.ID), bytes.NewReader(jsonReq))
			require.NoError(t, err)

			resp, err := srv.Client().Do(req)
			require.NoError(t, err)
			require.Equal(t, 400, resp.StatusCode)
		})
	})
}
