package keys

import (
	"context"
	"golang.org/x/text/language"
	"platform/internal/translations/entity/key"
	"platform/internal/translations/entity/project"
	"platform/internal/translations/entity/translation"
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
