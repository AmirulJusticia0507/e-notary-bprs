package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/repository"
)

type PasswordResetRepository struct{ db *sql.DB }

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(ctx context.Context, reset *domain.PasswordReset) error {
	const q = `INSERT INTO password_resets (user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowContext(ctx, q, reset.UserID, reset.TokenHash, reset.ExpiresAt, time.Now()).Scan(&reset.ID)
}

func (r *PasswordResetRepository) FindValidByTokenHash(ctx context.Context, tokenHash string) (*domain.PasswordReset, error) {
	const q = `SELECT id, user_id, token_hash, expires_at, created_at
		FROM password_resets WHERE token_hash = $1 AND expires_at > NOW()`
	reset := &domain.PasswordReset{}
	err := r.db.QueryRowContext(ctx, q, tokenHash).Scan(&reset.ID, &reset.UserID, &reset.TokenHash, &reset.ExpiresAt, &reset.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reset token invalid or expired: %w", repository.ErrNotFound)
		}
		return nil, err
	}
	return reset, nil
}

func (r *PasswordResetRepository) DeleteByID(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM password_resets WHERE id = $1`, id)
	return err
}

func (r *PasswordResetRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM password_resets WHERE expires_at <= NOW()`)
	return err
}

// ListPending untuk relay admin: siapa yang minta reset + kapan kedaluwarsa.
// Token plaintext TIDAK disimpan; admin menerbitkan kode baru per user
// via AdminCreateResetToken yang mengembalikan token sekali saja.
func (r *PasswordResetRepository) ListPending(ctx context.Context) ([]domain.PasswordReset, error) {
	const q = `SELECT pr.id, pr.user_id, u.email, pr.expires_at, pr.created_at
		FROM password_resets pr
		JOIN users u ON u.id = pr.user_id
		WHERE pr.expires_at > NOW() ORDER BY pr.created_at DESC`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	resets := make([]domain.PasswordReset, 0)
	for rows.Next() {
		var reset domain.PasswordReset
		if err := rows.Scan(&reset.ID, &reset.UserID, &reset.Email, &reset.ExpiresAt, &reset.CreatedAt); err != nil {
			return nil, err
		}
		resets = append(resets, reset)
	}
	return resets, rows.Err()
}

func (r *PasswordResetRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM password_resets WHERE user_id = $1`, userID)
	return err
}
