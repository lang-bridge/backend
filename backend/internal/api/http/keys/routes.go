package keys

import (
	"fmt"
	"golang.org/x/text/language"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"platform/internal/translations/entity/key"
	"platform/internal/translations/entity/project"
	"platform/internal/translations/service/keys"
	"platform/pkg/ctxlog"
	"platform/pkg/httputil"
	"platform/pkg/httputil/httperr"
)

type Router struct {
	svc keys.Service
}

func NewRouter(svc keys.Service) *Router {
	return &Router{svc: svc}
}

func (s *Router) Register(router chi.Router) {
	router.Route("/api/v1/projects/{projectID}/keys", func(router chi.Router) {
		router.Post("/", httputil.WrapError(s.CreateKey))
	})
}

// CreateKey creates key into project
//
//	@Summary	Create key with translates
//	@Tags		keys
//	@Accept		json
//	@Produce	json
//	@Param		projectID	path	int					true	"Project ID"
//	@Param		request		body	CreateKeyRequest	true	"Create key request"
//	@Router		/api/v1/projects/{projectID}/keys [post]
func (s *Router) CreateKey(w http.ResponseWriter, r *http.Request) error {
	projectID, err := getProjectID(r)
	if err != nil {
		return err
	}

	var req CreateKeyRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}
	if err := req.Validate(); err != nil {
		return httperr.BadRequest(fmt.Errorf("invalid request: %w", err))
	}

	var translate = make([]keys.Translate, len(req.Translates))
	for i, tr := range req.Translates {
		translate[i] = keys.Translate{
			Language: tr.Language,
			Value:    tr.Value,
		}
	}

	view, err := s.svc.CreateKey(r.Context(), keys.CreateKeyParam{
		ProjectID:   projectID,
		Name:        req.Name,
		Platforms:   req.Platforms,
		ExistedTags: req.ExistedTags,
		NewTags:     req.NewTags,
		Translate:   translate,
	})
	if err != nil {
		return fmt.Errorf("failed to create key: %w", err)
	}
	resp := CreateKeyResponse{
		ID: view.Key.ID,
	}

	ctxlog.Error(r.Context(), "test sentry error")

	return httputil.RenderJSON(w, http.StatusOK, resp)
}

type CreateKeyRequest struct {
	Name        string         `json:"name"`
	Platforms   []key.Platform `json:"platforms"`
	ExistedTags []key.TagID    `json:"existedTags"`
	NewTags     []string       `json:"newTags"`
	Translates  []Translate    `json:"translates"`
}

type Translate struct {
	Language language.Tag `json:"language" swaggertype:"string"`
	Value    string       `json:"value"`
}

func (c *CreateKeyRequest) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(c.Platforms) == 0 {
		return fmt.Errorf("platforms is required")
	}
	if len(c.Translates) == 0 {
		return fmt.Errorf("translates is required")
	}
	langs := map[string]any{}
	for _, tr := range c.Translates {
		if _, ok := langs[tr.Language.String()]; ok {
			return fmt.Errorf("duplicated language: %s", tr.Language.String())
		}
		langs[tr.Language.String()] = ""
	}
	return nil
}

type CreateKeyResponse struct {
	ID key.ID `json:"key_id"`
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
