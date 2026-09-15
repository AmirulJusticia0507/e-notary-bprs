package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/repository"
)

type FinancingRepository struct{ db *sql.DB }

func NewFinancingRepository(db *sql.DB) *FinancingRepository { return &FinancingRepository{db: db} }

func (r *FinancingRepository) Create(ctx context.Context, app *domain.FinancingApplication) error {
	if app.SyncedAt.IsZero() {
		app.SyncedAt = time.Now()
	}
	const q = `INSERT INTO financing_applications (customer_name, customer_nik, financing_amount, collateral_type, collateral_details, status, synced_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	return r.db.QueryRowContext(ctx, q, app.CustomerName, app.CustomerNIK, app.FinancingAmount, app.CollateralType, app.CollateralDetails, app.Status, app.SyncedAt, time.Now(), time.Now()).Scan(&app.ID)
}

func (r *FinancingRepository) FindAll(ctx context.Context) ([]domain.FinancingApplication, error) {
	const q = `SELECT id, customer_name, customer_nik, financing_amount, collateral_type, collateral_details, status, synced_at, created_at, updated_at
		FROM financing_applications WHERE deleted_at IS NULL ORDER BY id`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	apps := make([]domain.FinancingApplication, 0)
	for rows.Next() {
		var app domain.FinancingApplication
		if err := rows.Scan(&app.ID, &app.CustomerName, &app.CustomerNIK, &app.FinancingAmount, &app.CollateralType, &app.CollateralDetails, &app.Status, &app.SyncedAt, &app.CreatedAt, &app.UpdatedAt); err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, rows.Err()
}

func (r *FinancingRepository) FindByID(ctx context.Context, id int64) (*domain.FinancingApplication, error) {
	const q = `SELECT id, customer_name, customer_nik, financing_amount, collateral_type, collateral_details, status, synced_at, created_at, updated_at
		FROM financing_applications WHERE id = $1 AND deleted_at IS NULL`
	app := &domain.FinancingApplication{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(&app.ID, &app.CustomerName, &app.CustomerNIK, &app.FinancingAmount, &app.CollateralType, &app.CollateralDetails, &app.Status, &app.SyncedAt, &app.CreatedAt, &app.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("financing application not found: %w", repository.ErrNotFound)
		}
		return nil, err
	}
	return app, nil
}

func (r *FinancingRepository) SyncFromCBS(ctx context.Context, apps []domain.FinancingApplication) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	const q = `INSERT INTO financing_applications (customer_name, customer_nik, financing_amount, collateral_type, collateral_details, status, synced_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (customer_nik) DO UPDATE SET customer_name = EXCLUDED.customer_name, financing_amount = EXCLUDED.financing_amount,
		collateral_type = EXCLUDED.collateral_type, collateral_details = EXCLUDED.collateral_details, status = EXCLUDED.status,
		synced_at = EXCLUDED.synced_at, updated_at = EXCLUDED.updated_at`
	now := time.Now()
	for _, app := range apps {
		if _, err := tx.ExecContext(ctx, q, app.CustomerName, app.CustomerNIK, app.FinancingAmount, app.CollateralType, app.CollateralDetails, app.Status, now, now, now); err != nil {
			return err
		}
	}
	return tx.Commit()
}
