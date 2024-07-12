package keys

import (
	"context"
	"fmt"

	"platform/internal/translations/entity/key"
	"platform/internal/translations/entity/translation"
)

func (s *svc) CreateKey(ctx context.Context, params CreateKeyParam) (KeyView, error) {
	var resp KeyView
	err := s.txManager.Execute(ctx, func(ctx context.Context) error {
		tags := make([]key.TagID, 0, len(params.ExistedTags)+len(params.NewTags))
		copy(tags, params.ExistedTags)

		if len(params.NewTags) > 0 {
			newTags, err := s.tagsRepo.EnsureTags(ctx, params.ProjectID, params.NewTags)
			if err != nil {
				return fmt.Errorf("failed to create tags: %w", err)
			}
			for _, tag := range newTags {
				tags = append(tags, tag.ID)
			}
		}

		newKey, err := s.keysRepo.CreateKey(ctx, key.CreateKeyParam{
			ProjectID: params.ProjectID,
			Name:      params.Name,
			Platforms: params.Platforms,
			Tags:      tags,
		})
		if err != nil {
			return fmt.Errorf("failed to create key: %w", err)
		}

		translations := make([]translation.Value, len(params.Translate))
		for i, t := range params.Translate {
			translations[i] = translation.Value{
				KeyID:       newKey.ID,
				Language:    t.Language,
				Translation: t.Value,
			}
		}

		err = s.transRepo.UpsertTranslations(ctx, translations)
		if err != nil {
			return fmt.Errorf("failed to upsert translations: %w", err)
		}
		resp.Key = newKey
		resp.Translations = translations
		return nil
	})
	if err != nil {
		return resp, fmt.Errorf("failed to create key: %w", err)
	}
	return resp, nil
}
