package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azulrossini/budde-group-challenge/backend/internal/api"
	"github.com/azulrossini/budde-group-challenge/backend/internal/handlers"
	"github.com/azulrossini/budde-group-challenge/backend/internal/models"
	"github.com/azulrossini/budde-group-challenge/backend/internal/service"
)

// fakeRepository lets handler tests run without a database, per docs/SPEC.md §11.
type fakeRepository struct {
	bill models.Bill
}

func (f *fakeRepository) GetBillShares(ctx context.Context, billID int64) (models.BillShares, error) {
	return models.BillShares{Bill: f.bill}, nil
}

func (f *fakeRepository) ReplaceShares(ctx context.Context, billID int64, inputs []models.ShareInput) (models.BillShares, error) {
	panic("ReplaceShares should not be called when validation fails first")
}

func TestReplaceBillShares_InvalidSum_Returns422(t *testing.T) {
	repo := &fakeRepository{bill: models.Bill{ID: 1, Description: "Team dinner", TotalCents: 12050}}
	router := handlers.NewRouter(handlers.New(service.New(repo)), nil)

	body, err := json.Marshal(api.ReplaceSharesRequest{
		Shares: []api.ShareInput{
			{Name: "Alice", PercentageBasisPoints: 4950},
			{Name: "Bob", PercentageBasisPoints: 5000},
		},
	})
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/bills/1/shares", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusUnprocessableEntity, rec.Body.String())
	}

	var got api.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if got.Code != api.VALIDATIONERROR {
		t.Errorf("code = %q, want %q", got.Code, api.VALIDATIONERROR)
	}

	if got.Details == nil {
		t.Fatal("details is nil, want at least one field error")
	}

	foundSumError := false
	for _, d := range *got.Details {
		if d.Field == "shares" {
			foundSumError = true
		}
	}
	if !foundSumError {
		t.Errorf("details = %+v, want a \"shares\" field error for the sum mismatch", *got.Details)
	}
}
