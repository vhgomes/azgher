CREATE TABLE project_info (
    id             SERIAL PRIMARY KEY,
    project_id     INTEGER NOT NULL UNIQUE REFERENCES projects(id) ON DELETE CASCADE,
    summary        TEXT NOT NULL DEFAULT '',
    github_md_key  TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
