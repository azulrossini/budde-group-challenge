package service

import (
	"context"

	"github.com/azulrossini/budde-group-challenge/backend/internal/models"
	"github.com/azulrossini/budde-group-challenge/backend/internal/money"
)

type Repository interface {
	GetBillShares(ctx context.Context, billID int64) (models.BillShares, error)
	ReplaceShares(ctx context.Context, billID int64, inputs []models.ShareInput) (models.BillShares, error)
}

type ShareResult struct {
	ID                    int64
	Name                  string
	PercentageBasisPoints int32
	AmountCents           int64
}

type BillSharesResult struct {
	Bill             models.Bill
	Shares           []ShareResult
	TotalBasisPoints int32
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetBillShares(ctx context.Context, billID int64) (BillSharesResult, error) {
	bs, err := s.repo.GetBillShares(ctx, billID)
	if err != nil {
		return BillSharesResult{}, err
	}
	return buildResult(bs), nil
}

func (s *Service) ReplaceShares(ctx context.Context, billID int64, inputs []models.ShareInput) (BillSharesResult, error) {
	if err := ValidateReplaceShares(inputs); err != nil {
		return BillSharesResult{}, err
	}

	bs, err := s.repo.ReplaceShares(ctx, billID, inputs)
	if err != nil {
		return BillSharesResult{}, err
	}
	return buildResult(bs), nil
}

// buildResult computes each share's amount via the largest-remainder split.
// It assumes bs.Shares' percentages sum to 100.00 — guaranteed by
// ValidateReplaceShares being the only way to persist shares.
func buildResult(bs models.BillShares) BillSharesResult {
	basisPoints := make([]int32, len(bs.Shares))
	var total int32
	for i, share := range bs.Shares {
		basisPoints[i] = share.PercentageBasisPoints
		total += share.PercentageBasisPoints
	}

	amounts := money.Split(bs.Bill.TotalCents, basisPoints)

	shares := make([]ShareResult, len(bs.Shares))
	for i, share := range bs.Shares {
		shares[i] = ShareResult{
			ID:                    share.ID,
			Name:                  share.Name,
			PercentageBasisPoints: share.PercentageBasisPoints,
			AmountCents:           amounts[i],
		}
	}

	return BillSharesResult{Bill: bs.Bill, Shares: shares, TotalBasisPoints: total}
}
