package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/repository"
)

type LegalDocumentRepository struct{ db *sql.DB }

func NewLegalDocumentRepository(db *sql.DB) *LegalDocumentRepository {
	return &LegalDocumentRepository{db: db}
}

func (r *LegalDocumentRepository) Create(ctx context.Context, document *domain.LegalDocument) error {
	const q = `INSERT INTO legal_documents (order_id, file_url, sha256_hash, e_meterai_sn, e_sign_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRowContext(ctx, q, document.OrderID, document.FileURL, document.SHA256Hash, document.EMeteraiSN, document.ESignStatus, time.Now(), time.Now()).Scan(&document.ID)
}

func (r *LegalDocumentRepository) FindByOrderID(ctx context.Context, orderID int64) ([]domain.LegalDocument, error) {
	const q = `SELECT id, order_id, file_url, sha256_hash, e_meterai_sn, e_sign_status, created_at, updated_at
		FROM legal_documents WHERE order_id = $1 AND deleted_at IS NULL ORDER BY id`
	rows, err := r.db.QueryContext(ctx, q, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	documents := make([]domain.LegalDocument, 0)
	for rows.Next() {
		var document domain.LegalDocument
		if err := rows.Scan(&document.ID, &document.OrderID, &document.FileURL, &document.SHA256Hash, &document.EMeteraiSN, &document.ESignStatus, &document.CreatedAt, &document.UpdatedAt); err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	return documents, rows.Err()
}

func (r *LegalDocumentRepository) FindByID(ctx context.Context, id int64) (*domain.LegalDocument, error) {
	const q = `SELECT id, order_id, file_url, sha256_hash, e_meterai_sn, e_sign_status, created_at, updated_at
		FROM legal_documents WHERE id = $1 AND deleted_at IS NULL`
	document := &domain.LegalDocument{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(&document.ID, &document.OrderID, &document.FileURL, &document.SHA256Hash, &document.EMeteraiSN, &document.ESignStatus, &document.CreatedAt, &document.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("legal document not found: %w", repository.ErrNotFound)
		}
		return nil, err
	}
	return document, nil
}

func (r *LegalDocumentRepository) UpdateESignStatus(ctx context.Context, docID int64, status string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE legal_documents SET e_sign_status = $1, updated_at = $2 WHERE id = $3`, status, time.Now(), docID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("legal document not found: %w", repository.ErrNotFound)
	}
	return nil
}
