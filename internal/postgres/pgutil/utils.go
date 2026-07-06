package pgutil

func NilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
