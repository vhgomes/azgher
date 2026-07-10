CREATE TABLE projects (
    id               SERIAL PRIMARY KEY,
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name             TEXT NOT NULL,
    description      TEXT NOT NULL DEFAULT '',
    github_repo_link TEXT NOT NULL,
    project_status   TEXT NOT NULL DEFAULT 'pending',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_projects_status CHECK (project_status IN ('pending', 'ready', 'failed'))
);

CREATE INDEX idx_projects_user_id ON projects(user_id);
