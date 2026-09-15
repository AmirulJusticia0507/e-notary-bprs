package postgres

import (
	"context"
	"database/sql"

	"github.com/e-notary-bprs/backend/internal/domain"
)

type LegalOrderLogRepository struct{ db *sql.DB }

func NewLegalOrderLogRepository(db *sql.DB) *LegalOrderLogRepository {
	return &LegalOrderLogRepository{db: db}
}

func (r *LegalOrderLogRepository) FindByOrderID(ctx context.Context, orderID int64) ([]domain.LegalOrderLog, error) {
	const q = `SELECT id, order_id, prev_status, new_status, changed_by, created_at
		FROM legal_order_logs WHERE order_id = $1 ORDER BY created_at`
	rows, err := r.db.QueryContext(ctx, q, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	logs := make([]domain.LegalOrderLog, 0)
	for rows.Next() {
		var l domain.LegalOrderLog
		if err := rows.Scan(&l.ID, &l.OrderID, &l.PrevStatus, &l.NewStatus, &l.ChangedBy, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}
