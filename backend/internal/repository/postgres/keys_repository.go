package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"platform/internal/pkg/db/dbtx"

	"platform/internal/repository/postgres/gen"
	"platform/internal/translations/entity/key"
	"platform/internal/translations/entity/project"
)

type KeysRepository struct {
	db dbtx.DBTX
	q  *gen.Queries
}

func (k KeysRepository) CreateKey(ctx context.Context, params key.CreateKeyParam) (key.Key, error) {
	tags := make([]int64, len(params.Tags))
	for i, tagID := range params.Tags {
		tags[i] = int64(tagID)
	}
	platforms := make([]gen.Platform, len(params.Platforms))
	for i, platform := range params.Platforms {
		platforms[i] = gen.Platform(platform)
	}

	createdKey, err := k.q.CreateKey(ctx, gen.CreateKeyParams{
		ProjectID: int64(params.ProjectID),
		Name:      params.Name,
		Platforms: platforms,
		Tags:      tags,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return key.Key{}, key.ErrKeyAlreadyExists
		}
		return key.Key{}, fmt.Errorf("failed to create key: %w", err)
	}
	return mapRowKey(createdKey), nil
}

func mapRowKey(row gen.Key) key.Key {
	platforms := make([]key.Platform, len(row.Platforms))
	for i, platform := range row.Platforms {
		platforms[i] = key.Platform(platform)
	}
	tags := make([]key.TagID, len(row.Tags))
	for i, tag := range row.Tags {
		tags[i] = key.TagID(tag)
	}
	return key.Key{
		ID:        key.ID(row.ID),
		ProjectID: project.ID(row.ProjectID),
		Name:      row.Name,
		Platforms: platforms,
		Tags:      tags,
	}
}

func NewKeysRepository(db dbtx.DBTX, q *gen.Queries) key.KeysRepository {
	return &KeysRepository{db: db, q: q}
}
