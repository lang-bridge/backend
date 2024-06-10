package keys

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"platform/internal/translations/entity/key"
	"platform/internal/translations/entity/project"
	"platform/internal/translations/service/keys"
	"platform/pkg/httputil"
	"strconv"
)

type Router struct {
	svc keys.Service
}

func NewRouter(svc keys.Service) *Router {
	return &Router{svc: svc}
}

func (r *Router) Register(router chi.Router) {
	router.Route("/api/v1/projects/{projectID}/keys", func(router chi.Router) {
		router.Post("/", httputil.WrapError(r.CreateKey))
	})
}

// CreateKey creates key into project
// /api/v1/projects/{projectID}/keys
func (s *Router) CreateKey(w http.ResponseWriter, r *http.Request) error {
	projectID, err := getProjectID(r)
	if err != nil {
		return err
	}

	var req CreateKeyRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	view, err := s.svc.CreateKey(r.Context(), keys.CreateKeyParam{
		ProjectID:   projectID,
		Name:        req.Name,
		Platforms:   req.Platforms,
		ExistedTags: req.ExistedTags,
		NewTags:     req.NewTags,
		Translate:   req.Translates,
	})
	if err != nil {
		return fmt.Errorf("failed to create key: %w", err)
	}
	resp := CreateKeyResponse{
		ID: view.Key.ID,
	}

	return httputil.RenderJSON(w, http.StatusOK, resp)
}

type CreateKeyRequest struct {
	Name        string           `json:"name"`
	Platforms   []key.Platform   `json:"platforms"`
	ExistedTags []key.TagID      `json:"existedTags"`
	NewTags     []string         `json:"newTags"`
	Translates  []keys.Translate `json:"translates"`
}

type CreateKeyResponse struct {
	ID key.ID `json:"key_id"`
}

func (c CreateKeyResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func getProjectID(req *http.Request) (project.ID, error) {
	if param := chi.URLParam(req, "projectID"); param != "" {
		i, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse projectID: %w", err)
		}
		return project.ID(i), nil
	}

	return 0, fmt.Errorf("projectID is not presented into the query")
}
