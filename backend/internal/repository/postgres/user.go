package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/repository"
)

type UserRepository struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	const q = `INSERT INTO users (full_name, email, password_hash, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRowContext(ctx, q, user.FullName, user.Email, user.PasswordHash, user.Role, user.IsActive, time.Now(), time.Now()).Scan(&user.ID)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `SELECT id, full_name, email, password_hash, role, is_active, failed_login_attempts, locked_until, created_at, updated_at
		FROM users WHERE email = $1 AND deleted_at IS NULL`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, q, email).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", repository.ErrNotFound)
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	const q = `SELECT id, full_name, email, password_hash, role, is_active, failed_login_attempts, locked_until, created_at, updated_at
		FROM users WHERE id = $1 AND deleted_at IS NULL`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", repository.ErrNotFound)
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	const q = `SELECT id, full_name, email, password_hash, role, is_active, failed_login_attempts, locked_until, created_at, updated_at
		FROM users WHERE deleted_at IS NULL ORDER BY id`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	const q = `UPDATE users SET full_name = $1, email = $2, role = $3, is_active = $4, updated_at = $5 WHERE id = $6`
	result, err := r.db.ExecContext(ctx, q, user.FullName, user.Email, user.Role, user.IsActive, time.Now(), user.ID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("user not found: %w", repository.ErrNotFound)
	}
	return nil
}

// UpdatePassword mengganti hash password (dipakai reset/ganti password).
func (r *UserRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	const q = `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	result, err := r.db.ExecContext(ctx, q, passwordHash, time.Now(), id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("user not found: %w", repository.ErrNotFound)
	}
	return nil
}

// RecordFailedLogin mencatat gagal login beruntun dan waktu kunci (NULL bila belum dikunci).
func (r *UserRepository) RecordFailedLogin(ctx context.Context, id int64, attempts int, lockedUntil *time.Time) error {
	const q = `UPDATE users SET failed_login_attempts = $1, locked_until = $2, updated_at = $3 WHERE id = $4`
	result, err := r.db.ExecContext(ctx, q, attempts, lockedUntil, time.Now(), id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("user not found: %w", repository.ErrNotFound)
	}
	return nil
}

// ResetLoginAttempts menghapus hitungan gagal + kunci saat login sukses.
func (r *UserRepository) ResetLoginAttempts(ctx context.Context, id int64) error {
	const q = `UPDATE users SET failed_login_attempts = 0, locked_until = NULL, updated_at = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, q, time.Now(), id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("user not found: %w", repository.ErrNotFound)
	}
	return nil
}
