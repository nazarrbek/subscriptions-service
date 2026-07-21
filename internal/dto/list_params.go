package dto

// ListParams holds pagination parameters for list queries.
type ListParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
