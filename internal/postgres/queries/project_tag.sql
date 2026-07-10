-- name: AddTagToProject :exec
INSERT INTO project_tags (project_id, tag_id)
VALUES ($1, $2)
ON CONFLICT (project_id, tag_id) DO NOTHING;

-- name: RemoveTagFromProject :exec
DELETE FROM project_tags
WHERE project_id = $1 AND tag_id = $2;

-- name: RemoveAllTagsFromProject :exec
DELETE FROM project_tags
WHERE project_id = $1;

-- name: ListTagsForProject :many
-- Retorna as tags completas (sem informações do projeto)
SELECT t.id, t.category, t.value
FROM tags t
JOIN project_tags pt ON t.id = pt.tag_id
WHERE pt.project_id = $1
ORDER BY t.category, t.value;

-- name: ListProjectsForTag :many
-- Retorna os projetos que possuem uma determinada tag
SELECT p.id, p.user_id, p.name, p.description, p.github_repo_link, p.project_status, p.created_at, p.updated_at
FROM projects p
JOIN project_tags pt ON p.id = pt.project_id
WHERE pt.tag_id = $1
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountTagsForProject :one
-- Útil para verificar se um projeto já tem tags
SELECT COUNT(*) AS total
FROM project_tags
WHERE project_id = $1;
