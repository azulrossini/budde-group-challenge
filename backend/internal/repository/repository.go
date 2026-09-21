package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/azulrossini/budde-group-challenge/backend/internal/apperr"
	"github.com/azulrossini/budde-group-challenge/backend/internal/models"
)

const msgBillNotFound = "bill not found"

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetBillShares(ctx context.Context, billID int64) (models.BillShares, error) {
	bill, err := scanBill(r.pool.QueryRow(ctx,
		`SELECT id, description, total_cents FROM bills WHERE id = $1`,
		billID,
	))
	if err != nil {
		return models.BillShares{}, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, bill_id, person_name, percentage_bp FROM shares WHERE bill_id = $1 ORDER BY id`,
		billID,
	)
	if err != nil {
		return models.BillShares{}, apperr.Internal(err)
	}
	defer rows.Close()

	shares, err := scanShares(rows)
	if err != nil {
		return models.BillShares{}, err
	}

	return models.BillShares{Bill: bill, Shares: shares}, nil
}

func (r *Repository) ReplaceShares(ctx context.Context, billID int64, inputs []models.ShareInput) (models.BillShares, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return models.BillShares{}, apperr.Internal(err)
	}
	defer tx.Rollback(ctx)

	bill, err := scanBill(tx.QueryRow(ctx,
		`SELECT id, description, total_cents FROM bills WHERE id = $1 FOR UPDATE`,
		billID,
	))
	if err != nil {
		return models.BillShares{}, err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM shares WHERE bill_id = $1`, billID); err != nil {
		return models.BillShares{}, apperr.Internal(err)
	}

	shares, err := insertShares(ctx, tx, billID, inputs)
	if err != nil {
		return models.BillShares{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.BillShares{}, apperr.Internal(err)
	}

	return models.BillShares{Bill: bill, Shares: shares}, nil
}

func scanBill(row pgx.Row) (models.Bill, error) {
	var bill models.Bill
	err := row.Scan(&bill.ID, &bill.Description, &bill.TotalCents)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Bill{}, apperr.NotFound(msgBillNotFound)
		}
		return models.Bill{}, apperr.Internal(err)
	}
	return bill, nil
}

func scanShares(rows pgx.Rows) ([]models.Share, error) {
	var shares []models.Share
	for rows.Next() {
		var share models.Share
		if err := rows.Scan(&share.ID, &share.BillID, &share.Name, &share.PercentageBasisPoints); err != nil {
			return nil, apperr.Internal(err)
		}
		shares = append(shares, share)
	}
	if err := rows.Err(); err != nil {
		return nil, apperr.Internal(err)
	}
	return shares, nil
}

func insertShares(ctx context.Context, tx pgx.Tx, billID int64, inputs []models.ShareInput) ([]models.Share, error) {
	batch := &pgx.Batch{}
	for _, input := range inputs {
		batch.Queue(
			`INSERT INTO shares (bill_id, person_name, percentage_bp) VALUES ($1, $2, $3)
			 RETURNING id, bill_id, person_name, percentage_bp`,
			billID, input.Name, input.PercentageBasisPoints,
		)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	shares := make([]models.Share, len(inputs))
	for i := range inputs {
		if err := results.QueryRow().Scan(&shares[i].ID, &shares[i].BillID, &shares[i].Name, &shares[i].PercentageBasisPoints); err != nil {
			return nil, apperr.Internal(err)
		}
	}

	return shares, nil
}
