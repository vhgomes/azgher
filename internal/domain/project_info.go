package domain

type ProjectInfo struct {
	ID          int      `db:"id" json:"id"`
	ProjectID   int      `db:"project_id" json:"project_id"`
	Summary     string   `db:"summary" json:"summary"`
	GithubMDKey string   `db:"github_md_key" json:"-"`
	Tags        []string `db:"-" json:"tags"`
}
