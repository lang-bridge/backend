package translation

import (
	"context"

	"golang.org/x/text/language"
	"platform/internal/translations/entity/key"
)

type Value struct {
	KeyID       key.ID
	Translation string
	Language    language.Tag
}

type Repository interface {
	UpsertTranslations(ctx context.Context, translation []Value) error
}
