package test

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateKey(t *testing.T) {
	RunTest(t, func(s chi.Router) {
		srv := httptest.NewServer(s)
		t.Cleanup(func() {
			srv.Close()
		})

		jsonReq := `{"name":"key123","platforms":["web"],"existedTags":[123],"newTags":["tag2"],"translates":[{"lang":"en","value":"value"}]}`
		req, err := http.NewRequest("POST", srv.URL+"/api/v1/projects/123/keys", strings.NewReader(jsonReq))
		require.NoError(t, err)

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
	})
}
