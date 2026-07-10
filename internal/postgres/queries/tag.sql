-- name: CreateTag :one
INSERT INTO tags (
    category,
    value
) VALUES (
    $1, $2
)
ON CONFLICT (category, value) DO UPDATE SET
    -- Se já existe, retorna o existente (upsert silencioso)
    -- O sqlc precisa do RETURNING; com DO UPDATE SET id = EXCLUDED.id força retorno
    category = EXCLUDED.category
RETURNING id, category, value;

-- name: GetTagByID :one
SELECT id, category, value
FROM tags
WHERE id = $1;

-- name: GetTagByCategoryAndValue :one
SELECT id, category, value
FROM tags
WHERE category = $1 AND value = $2;

-- name: ListTagsByCategory :many
SELECT id, category, value
FROM tags
WHERE category = $1
ORDER BY value
LIMIT $2 OFFSET $3;

-- name: ListTags :many
SELECT id, category, value
FROM tags
ORDER BY category, value
LIMIT $1 OFFSET $2;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE id = $1;
