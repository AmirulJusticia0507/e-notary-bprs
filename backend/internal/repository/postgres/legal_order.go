package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/repository"
)

type LegalOrderRepository struct{ db *sql.DB }

func NewLegalOrderRepository(db *sql.DB) *LegalOrderRepository { return &LegalOrderRepository{db: db} }

func (r *LegalOrderRepository) Create(ctx context.Context, order *domain.LegalOrder) error {
	const q = `INSERT INTO legal_orders (order_number, financing_id, notary_id, assigned_to, status, sla_deadline, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return r.db.QueryRowContext(ctx, q, order.OrderNumber, order.FinancingID, order.NotaryID, order.AssignedTo, order.Status, order.SLADeadline, time.Now(), time.Now()).Scan(&order.ID)
}

func (r *LegalOrderRepository) FindAll(ctx context.Context) ([]domain.OrderDetail, error) {
	return r.queryOrderDetails(ctx, `SELECT o.id, o.order_number, o.financing_id, o.notary_id, o.assigned_to, o.status, o.sla_deadline,
		o.completed_at, o.created_at, o.updated_at, f.customer_name, n.full_name AS notary_name, u.full_name AS assigned_name
		FROM legal_orders o
		JOIN financing_applications f ON f.id = o.financing_id
		JOIN notaries n ON n.id = o.notary_id
		JOIN users u ON u.id = o.assigned_to
		WHERE o.deleted_at IS NULL ORDER BY o.id`)
}

func (r *LegalOrderRepository) FindByID(ctx context.Context, id int64) (*domain.OrderDetail, error) {
	const q = `SELECT o.id, o.order_number, o.financing_id, o.notary_id, o.assigned_to, o.status, o.sla_deadline,
		o.completed_at, o.created_at, o.updated_at, f.customer_name, n.full_name AS notary_name, u.full_name AS assigned_name
		FROM legal_orders o
		JOIN financing_applications f ON f.id = o.financing_id
		JOIN notaries n ON n.id = o.notary_id
		JOIN users u ON u.id = o.assigned_to
		WHERE o.id = $1 AND o.deleted_at IS NULL`
	order := &domain.OrderDetail{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(&order.ID, &order.OrderNumber, &order.FinancingID, &order.NotaryID, &order.AssignedTo, &order.Status, &order.SLADeadline, &order.CompletedAt, &order.CreatedAt, &order.UpdatedAt, &order.CustomerName, &order.NotaryName, &order.AssignedName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("legal order not found: %w", repository.ErrNotFound)
		}
		return nil, err
	}
	return order, nil
}

func (r *LegalOrderRepository) FindByNotaryID(ctx context.Context, notaryID int64) ([]domain.OrderDetail, error) {
	return r.queryOrderDetails(ctx, `SELECT o.id, o.order_number, o.financing_id, o.notary_id, o.assigned_to, o.status, o.sla_deadline,
		o.completed_at, o.created_at, o.updated_at, f.customer_name, n.full_name AS notary_name, u.full_name AS assigned_name
		FROM legal_orders o
		JOIN financing_applications f ON f.id = o.financing_id
		JOIN notaries n ON n.id = o.notary_id
		JOIN users u ON u.id = o.assigned_to
		WHERE o.notary_id = $1 AND o.deleted_at IS NULL ORDER BY o.id`, notaryID)
}

func (r *LegalOrderRepository) FindByAssignedTo(ctx context.Context, assignedTo int64) ([]domain.OrderDetail, error) {
	return r.queryOrderDetails(ctx, `SELECT o.id, o.order_number, o.financing_id, o.notary_id, o.assigned_to, o.status, o.sla_deadline,
		o.completed_at, o.created_at, o.updated_at, f.customer_name, n.full_name AS notary_name, u.full_name AS assigned_name
		FROM legal_orders o
		JOIN financing_applications f ON f.id = o.financing_id
		JOIN notaries n ON n.id = o.notary_id
		JOIN users u ON u.id = o.assigned_to
		WHERE o.assigned_to = $1 AND o.deleted_at IS NULL ORDER BY o.id`, assignedTo)
}

func (r *LegalOrderRepository) queryOrderDetails(ctx context.Context, q string, args ...any) ([]domain.OrderDetail, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := make([]domain.OrderDetail, 0)
	for rows.Next() {
		var order domain.OrderDetail
		if err := rows.Scan(&order.ID, &order.OrderNumber, &order.FinancingID, &order.NotaryID, &order.AssignedTo, &order.Status, &order.SLADeadline, &order.CompletedAt, &order.CreatedAt, &order.UpdatedAt, &order.CustomerName, &order.NotaryName, &order.AssignedName); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *LegalOrderRepository) UpdateStatus(ctx context.Context, orderID int64, newStatus string, changedBy int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := time.Now()
	var prevStatus string
	if err := tx.QueryRowContext(ctx, `SELECT status FROM legal_orders WHERE id = $1`, orderID).Scan(&prevStatus); err != nil {
		return err
	}
	result, err := tx.ExecContext(ctx, `UPDATE legal_orders SET status = $1, updated_at = $2,
		completed_at = CASE WHEN $1 = 'completed' THEN $2 ELSE completed_at END WHERE id = $3`, newStatus, now, orderID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("legal order not found: %w", repository.ErrNotFound)
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO legal_order_logs (order_id, prev_status, new_status, changed_by, created_at)
		VALUES ($1, $2, $3, $4, $5)`, orderID, prevStatus, newStatus, changedBy, now)
	if err != nil {
		return err
	}
	return tx.Commit()
}
