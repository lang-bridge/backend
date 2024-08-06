package postgres

import (
	"context"
	"fmt"

	"platform/internal/pkg/db/dbtx"

	"platform/internal/repository/postgres/gen"
	"platform/internal/translations/entity/translation"
)

type TranslationsRepository struct {
	db dbtx.DBTX
}

const upsertTranslations = `
INSERT INTO translations (key_id, language, value) 
VALUES (:key_id, :language, :value)
ON CONFLICT (key_id, language) DO UPDATE SET value = EXCLUDED.value
`

func (t TranslationsRepository) UpsertTranslations(ctx context.Context, translation []translation.Value) error {
	rows := make([]gen.Translation, len(translation))
	for i, t := range translation {
		rows[i] = gen.Translation{
			KeyID:    int64(t.KeyID),
			Language: t.Language,
			Value:    t.Translation,
		}
	}
	_, err := t.db.NamedExecContext(ctx, upsertTranslations, rows)
	if err != nil {
		return fmt.Errorf("failed to upsert translations: %w", err)
	}
	return nil
}

func NewTranslationsRepository(db dbtx.DBTX) translation.Repository {
	return &TranslationsRepository{db: db}
}
