package dto

import "fmt"

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

// Validate checks that all required fields are present and values are sane.
func (r *CreateSubscriptionRequest) Validate() error {
	if r.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}
	if r.Price < 0 {
		return fmt.Errorf("price must be non-negative")
	}
	if r.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if r.StartDate == "" {
		return fmt.Errorf("start_date is required")
	}
	return nil
}
