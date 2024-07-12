package keys

import (
	"context"

	"golang.org/x/text/language"
	"platform/internal/translations/entity/key"
	"platform/internal/translations/entity/project"
	"platform/internal/translations/entity/translation"
	"platform/pkg/db/tx"
)

type Service interface {
	CreateKey(ctx context.Context, params CreateKeyParam) (KeyView, error)
}

type CreateKeyParam struct {
	ProjectID   project.ID
	Name        string
	Platforms   []key.Platform
	ExistedTags []key.TagID
	NewTags     []string
	Translate   []Translate
}

type Translate struct {
	Language language.Tag `json:"language"`
	Value    string       `json:"value"`
}

type KeyView struct {
	Key          key.Key
	Translations []translation.Value
}

type svc struct {
	keysRepo  key.KeysRepository
	tagsRepo  key.TagsRepository
	transRepo translation.Repository
	txManager tx.Manager
}

func NewService(
	keysRepo key.KeysRepository,
	tagsRepo key.TagsRepository,
	transRepo translation.Repository,
	txManager tx.Manager,
) Service {
	return &svc{
		keysRepo:  keysRepo,
		tagsRepo:  tagsRepo,
		transRepo: transRepo,
		txManager: txManager,
	}
}
