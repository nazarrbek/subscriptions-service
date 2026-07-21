package dto

import "testing"

func TestCreateSubscriptionRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateSubscriptionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: CreateSubscriptionRequest{
				ServiceName: "Netflix",
				Price:       999,
				UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				StartDate:   "01-2025",
			},
			wantErr: false,
		},
		{
			name: "valid with zero price",
			req: CreateSubscriptionRequest{
				ServiceName: "Free Tier",
				Price:       0,
				UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				StartDate:   "01-2025",
			},
			wantErr: false,
		},
		{
			name: "empty service_name",
			req: CreateSubscriptionRequest{
				ServiceName: "",
				Price:       999,
				UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				StartDate:   "01-2025",
			},
			wantErr: true,
			errMsg:  "service_name is required",
		},
		{
			name: "negative price",
			req: CreateSubscriptionRequest{
				ServiceName: "Netflix",
				Price:       -100,
				UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				StartDate:   "01-2025",
			},
			wantErr: true,
			errMsg:  "price must be non-negative",
		},
		{
			name: "empty user_id",
			req: CreateSubscriptionRequest{
				ServiceName: "Netflix",
				Price:       999,
				UserID:      "",
				StartDate:   "01-2025",
			},
			wantErr: true,
			errMsg:  "user_id is required",
		},
		{
			name: "empty start_date",
			req: CreateSubscriptionRequest{
				ServiceName: "Netflix",
				Price:       999,
				UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				StartDate:   "",
			},
			wantErr: true,
			errMsg:  "start_date is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errMsg)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestUpdateSubscriptionRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     UpdateSubscriptionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: UpdateSubscriptionRequest{
				ServiceName: "Netflix",
				Price:       999,
				StartDate:   "01-2025",
			},
			wantErr: false,
		},
		{
			name: "empty service_name",
			req: UpdateSubscriptionRequest{
				ServiceName: "",
				Price:       999,
				StartDate:   "01-2025",
			},
			wantErr: true,
			errMsg:  "service_name is required",
		},
		{
			name: "negative price",
			req: UpdateSubscriptionRequest{
				ServiceName: "Netflix",
				Price:       -1,
				StartDate:   "01-2025",
			},
			wantErr: true,
			errMsg:  "price must be non-negative",
		},
		{
			name: "empty start_date",
			req: UpdateSubscriptionRequest{
				ServiceName: "Netflix",
				Price:       999,
				StartDate:   "",
			},
			wantErr: true,
			errMsg:  "start_date is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errMsg)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
