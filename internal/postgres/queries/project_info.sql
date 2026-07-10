-- name: CreateProjectInfo :one
INSERT INTO project_info (
    project_id,
    summary,
    github_md_key
) VALUES (
    $1, $2, $3
)
RETURNING id, project_id, summary, github_md_key, created_at, updated_at;

-- name: UpsertProjectInfo :one
-- Usado para criar ou atualizar a info de um projeto (evita duplicidade por project_id)
INSERT INTO project_info (
    project_id,
    summary,
    github_md_key
) VALUES (
    $1, $2, $3
) ON CONFLICT (project_id) DO UPDATE SET
    summary = EXCLUDED.summary,
    github_md_key = EXCLUDED.github_md_key,
    updated_at = now()
RETURNING id, project_id, summary, github_md_key, created_at, updated_at;

-- name: GetProjectInfoByProjectID :one
SELECT id, project_id, summary, github_md_key, created_at, updated_at
FROM project_info
WHERE project_id = $1;

-- name: GetProjectInfoByID :one
SELECT id, project_id, summary, github_md_key, created_at, updated_at
FROM project_info
WHERE id = $1;

-- name: DeleteProjectInfo :exec
DELETE FROM project_info
WHERE project_id = $1;
