package service

import (
	"testing"

	"github.com/azulrossini/budde-group-challenge/backend/internal/apperr"
	"github.com/azulrossini/budde-group-challenge/backend/internal/models"
)

func TestValidateReplaceShares_Valid(t *testing.T) {
	err := ValidateReplaceShares([]models.ShareInput{
		{Name: "Alice", PercentageBasisPoints: 5000},
		{Name: "Bob", PercentageBasisPoints: 3000},
		{Name: "Carol", PercentageBasisPoints: 2000},
	})
	if err != nil {
		t.Fatalf("ValidateReplaceShares() = %v, want nil", err)
	}
}

func TestValidateReplaceShares_Rejections(t *testing.T) {
	tests := []struct {
		name   string
		inputs []models.ShareInput
	}{
		{
			name: "sum below 100.00",
			inputs: []models.ShareInput{
				{Name: "Alice", PercentageBasisPoints: 4999},
				{Name: "Bob", PercentageBasisPoints: 5000},
			},
		},
		{
			name: "sum above 100.00",
			inputs: []models.ShareInput{
				{Name: "Alice", PercentageBasisPoints: 5001},
				{Name: "Bob", PercentageBasisPoints: 5000},
			},
		},
		{
			name: "zero percentage",
			inputs: []models.ShareInput{
				{Name: "Alice", PercentageBasisPoints: 0},
				{Name: "Bob", PercentageBasisPoints: 10000},
			},
		},
		{
			name: "negative percentage",
			inputs: []models.ShareInput{
				{Name: "Alice", PercentageBasisPoints: -100},
				{Name: "Bob", PercentageBasisPoints: 10100},
			},
		},
		{
			name: "empty name",
			inputs: []models.ShareInput{
				{Name: "   ", PercentageBasisPoints: 10000},
			},
		},
		{
			name: "duplicate name, case-insensitive",
			inputs: []models.ShareInput{
				{Name: "Alice", PercentageBasisPoints: 5000},
				{Name: "alice", PercentageBasisPoints: 5000},
			},
		},
		{
			name:   "empty list",
			inputs: []models.ShareInput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReplaceShares(tt.inputs)
			if err == nil {
				t.Fatal("ValidateReplaceShares() = nil, want a validation error")
			}
			if err.Code != apperr.CodeValidationError {
				t.Errorf("err.Code = %v, want VALIDATION_ERROR", err.Code)
			}
			if len(err.Details) == 0 {
				t.Error("err.Details is empty, want at least one field error")
			}
		})
	}
}

func TestValidateReplaceShares_ReturnsAllErrorsAtOnce(t *testing.T) {
	err := ValidateReplaceShares([]models.ShareInput{
		{Name: "", PercentageBasisPoints: -100},
		{Name: "Alice", PercentageBasisPoints: 5000},
		{Name: "alice", PercentageBasisPoints: 200000},
	})
	if err == nil {
		t.Fatal("ValidateReplaceShares() = nil, want a validation error")
	}

	// Expect: empty name, negative percentage, duplicate name,
	// out-of-range percentage, and a sum mismatch — five distinct errors.
	const wantDetails = 5
	if len(err.Details) != wantDetails {
		t.Errorf("len(err.Details) = %d, want %d: %+v", len(err.Details), wantDetails, err.Details)
	}
}
