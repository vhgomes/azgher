package domain

type Tag struct {
	ID       int    `db:"id" json:"id"`
	Category string `db:"category" json:"category"`
	Value    string `db:"value" json:"value"`
}

type ProjectTag struct {
	ProjectID int `db:"project_id" json:"project_id"`
	TagID     int `db:"tag_id" json:"tag_id"`
}
