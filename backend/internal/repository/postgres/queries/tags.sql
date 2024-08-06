-- name: SelectTags :many
SELECT *
FROM key_tags
WHERE project_id = $1
  AND LOWER(value) = ANY ($2::varchar[]);