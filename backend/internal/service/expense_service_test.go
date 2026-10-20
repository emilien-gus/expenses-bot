package service

import (
	"strings"
	"testing"

	"github.com/ai-bot-expenses/backend/internal/models"
)

func TestExpenseService_CreateValidation(t *testing.T) {
	svc := NewExpenseService(nil)

	tests := []struct {
		name        string
		req         models.CreateExpenseRequest
		expectedErr string
	}{
		{
			name: "Invalid category rejection",
			req: models.CreateExpenseRequest{
				UserID:   123,
				Amount:   350,
				Category: "неизвестная_категория",
				SpentAt:  "2026-10-19",
				Comment:  "кофе",
			},
			expectedErr: "invalid category",
		},
		{
			name: "Invalid date format rejection",
			req: models.CreateExpenseRequest{
				UserID:   123,
				Amount:   350,
				Category: "еда",
				SpentAt:  "19-10-2026", // should be YYYY-MM-DD
				Comment:  "кофе",
			},
			expectedErr: "invalid date format",
		},
		{
			name: "Non-date string rejection",
			req: models.CreateExpenseRequest{
				UserID:   123,
				Amount:   350,
				Category: "еда",
				SpentAt:  "вчера",
				Comment:  "кофе",
			},
			expectedErr: "invalid date format",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Create(tc.req)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.expectedErr)
			}
			if !strings.Contains(err.Error(), tc.expectedErr) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.expectedErr)
			}
		})
	}
}
