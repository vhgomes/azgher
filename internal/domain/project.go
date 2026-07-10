package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProjectStatus string

const (
	ProjectStatusPending ProjectStatus = "pending"
	ProjectStatusReady   ProjectStatus = "ready"
	ProjectStatusFailed  ProjectStatus = "failed"
)

type Project struct {
	ID             int           `db:"id" json:"id"`
	UserID         uuid.UUID     `db:"user_id" json:"user_id"`
	Name           string        `db:"name" json:"name"`
	Description    string        `db:"description" json:"description"`
	GithubRepoLink string        `db:"github_repo_link" json:"github_repo_link"`
	Status         ProjectStatus `db:"project_status" json:"project_status"`
	CreatedAt      time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `db:"updated_at" json:"updated_at"`
}
