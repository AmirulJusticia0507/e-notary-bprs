package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
)

type CollateralRepository struct{ db *sql.DB }

func NewCollateralRepository(db *sql.DB) *CollateralRepository { return &CollateralRepository{db: db} }

func (r *CollateralRepository) Create(ctx context.Context, collateral *domain.Collateral) error {
	const q = `INSERT INTO collaterals (financing_id, type, details, estimated_value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRowContext(ctx, q, collateral.FinancingID, collateral.Type, collateral.Details, collateral.EstimatedValue, time.Now(), time.Now()).Scan(&collateral.ID)
}

func (r *CollateralRepository) FindByFinancingID(ctx context.Context, financingID int64) ([]domain.Collateral, error) {
	const q = `SELECT id, financing_id, type, details, estimated_value, created_at, updated_at
		FROM collaterals WHERE financing_id = $1 AND deleted_at IS NULL ORDER BY id`
	rows, err := r.db.QueryContext(ctx, q, financingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	collaterals := make([]domain.Collateral, 0)
	for rows.Next() {
		var collateral domain.Collateral
		if err := rows.Scan(&collateral.ID, &collateral.FinancingID, &collateral.Type, &collateral.Details, &collateral.EstimatedValue, &collateral.CreatedAt, &collateral.UpdatedAt); err != nil {
			return nil, err
		}
		collaterals = append(collaterals, collateral)
	}
	return collaterals, rows.Err()
}
