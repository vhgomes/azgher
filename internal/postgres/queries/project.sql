-- name: CreateProject :one
INSERT INTO projects (
    user_id,
    name,
    description,
    github_repo_link,
    project_status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, user_id, name, description, github_repo_link, project_status, created_at, updated_at;

-- name: GetProjectByID :one
SELECT id, user_id, name, description, github_repo_link, project_status, created_at, updated_at
FROM projects
WHERE id = $1;

-- name: GetProjectByIDAndUser :one
SELECT id, user_id, name, description, github_repo_link, project_status, created_at, updated_at
FROM projects
WHERE id = $1 AND user_id = $2;

-- name: ListProjectsByUser :many
SELECT id, user_id, name, description, github_repo_link, project_status, created_at, updated_at
FROM projects
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListProjects :many
SELECT id, user_id, name, description, github_repo_link, project_status, created_at, updated_at
FROM projects
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateProject :one
UPDATE projects
SET
    name = $2,
    description = $3,
    github_repo_link = $4,
    project_status = $5,
    updated_at = now()
WHERE id = $1
RETURNING id, user_id, name, description, github_repo_link, project_status, created_at, updated_at;

-- name: UpdateProjectStatus :one
UPDATE projects
SET
    project_status = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, user_id, name, description, github_repo_link, project_status, created_at, updated_at;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;
