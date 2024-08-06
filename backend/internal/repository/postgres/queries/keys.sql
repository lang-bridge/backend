-- name: CreateKey :one
INSERT INTO keys (project_id, name, platforms, tags)
VALUES ($1, $2, $3, $4)
ON CONFLICT (project_id, name) DO NOTHING
RETURNING *;