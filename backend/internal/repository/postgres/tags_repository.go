package postgres

import (
	"context"
	"fmt"
	"platform/internal/pkg/ctxlog"
	"platform/internal/pkg/db/dbtx"
	"strings"

	"platform/internal/repository/postgres/gen"
	"platform/internal/translations/entity/key"
	"platform/internal/translations/entity/project"
)

type TagsRepository struct {
	db dbtx.DBTX
	q  *gen.Queries
}

const createTagsQuery = `
INSERT INTO key_tags (project_id, value) 
VALUES (:project_id, :value)
ON CONFLICT (project_id, lower(value)) DO NOTHING
RETURNING id, project_id, value
`

type createTag struct {
	ProjectID project.ID `db:"project_id"`
	Value     string     `db:"value"`
}

func (t *TagsRepository) EnsureTags(ctx context.Context, projectID project.ID, tags []string) ([]key.Tag, error) {
	params := make([]createTag, len(tags))
	for i, tag := range tags {
		params[i] = createTag{
			ProjectID: projectID,
			Value:     tag,
		}
	}
	rows, err := t.db.NamedQueryContext(ctx, createTagsQuery, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create tags in db: %w", err)
	}
	defer rows.Close()

	allTags := map[string]string{}
	for _, tag := range tags {
		allTags[strings.ToLower(tag)] = tag
	}

	var result []key.Tag
	for rows.Next() {
		var tag key.Tag
		err := rows.Scan(&tag.ID, &tag.ProjectID, &tag.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		result = append(result, tag)
		delete(allTags, strings.ToLower(tag.Value))
	}

	if len(allTags) > 0 {
		ctxlog.Warn(ctx, "some tags were not created, because they already exist in the database")
		keys := make([]string, 0, len(allTags))
		for k := range allTags {
			keys = append(keys, k)
		}

		selectedTags, err := t.q.SelectTags(ctx, gen.SelectTagsParams{
			ProjectID: int64(projectID),
			Column2:   keys,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to select tags: %w", err)
		}
		for _, row := range selectedTags {
			tag := key.Tag{
				ID:        key.TagID(row.ID),
				ProjectID: projectID,
				Value:     row.Value,
			}
			result = append(result, tag)
			delete(allTags, strings.ToLower(tag.Value))
		}
		if len(allTags) > 0 {
			ctxlog.Error(ctx, "some tags were not found in the database")
		}
	}

	return result, nil
}

func NewTagsRepository(db dbtx.DBTX, q *gen.Queries) key.TagsRepository {
	return &TagsRepository{
		db: db,
		q:  q,
	}
}
