package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/repository"
)

type NotaryRepository struct{ db *sql.DB }

func NewNotaryRepository(db *sql.DB) *NotaryRepository { return &NotaryRepository{db: db} }

func (r *NotaryRepository) Create(ctx context.Context, notary *domain.Notary) error {
	const q = `INSERT INTO notaries (full_name, notary_number, wilayah_kerja, phone_number, email, is_available, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return r.db.QueryRowContext(ctx, q,
		notary.FullName, notary.NotaryNumber, notary.WilayahKerja, notary.PhoneNumber, notary.Email, notary.IsAvailable,
		time.Now(), time.Now(),
	).Scan(&notary.ID)
}

func (r *NotaryRepository) FindAll(ctx context.Context) ([]domain.Notary, error) {
	const q = `SELECT id, full_name, notary_number, wilayah_kerja, phone_number, email, is_available, created_at, updated_at
		FROM notaries WHERE deleted_at IS NULL ORDER BY id`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	notaries := make([]domain.Notary, 0)
	for rows.Next() {
		var n domain.Notary
		if err := rows.Scan(&n.ID, &n.FullName, &n.NotaryNumber, &n.WilayahKerja, &n.PhoneNumber, &n.Email, &n.IsAvailable, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notaries = append(notaries, n)
	}
	return notaries, rows.Err()
}

func (r *NotaryRepository) FindByID(ctx context.Context, id int64) (*domain.Notary, error) {
	const q = `SELECT id, full_name, notary_number, wilayah_kerja, phone_number, email, is_available, created_at, updated_at
		FROM notaries WHERE id = $1 AND deleted_at IS NULL`
	n := &domain.Notary{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&n.ID, &n.FullName, &n.NotaryNumber, &n.WilayahKerja, &n.PhoneNumber, &n.Email, &n.IsAvailable, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notary not found: %w", repository.ErrNotFound)
		}
		return nil, err
	}
	return n, nil
}

func (r *NotaryRepository) Update(ctx context.Context, notary *domain.Notary) error {
	const q = `UPDATE notaries SET full_name = $1, notary_number = $2, wilayah_kerja = $3, phone_number = $4, email = $5, is_available = $6, updated_at = $7
		WHERE id = $8 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, q, notary.FullName, notary.NotaryNumber, notary.WilayahKerja, notary.PhoneNumber, notary.Email, notary.IsAvailable, time.Now(), notary.ID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("notary not found: %w", repository.ErrNotFound)
	}
	return nil
}

func (r *NotaryRepository) Delete(ctx context.Context, id int64) error {
	const q = `UPDATE notaries SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, q, time.Now(), id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("notary not found: %w", repository.ErrNotFound)
	}
	return nil
}
