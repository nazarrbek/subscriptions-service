package dto

import "github.com/nazarrbek/subscriptions-service/internal/models"

// ListResponse wraps a paginated list of subscriptions with metadata.
type ListResponse struct {
	Data   []models.Subscription `json:"data"`
	Total  int                   `json:"total"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
}
