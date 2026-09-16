package repository

import (
	"context"
	"errors"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
)

var ErrNotFound = errors.New("record not found")

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type NotaryRepository interface {
	Create(ctx context.Context, notary *domain.Notary) error
	FindAll(ctx context.Context) ([]domain.Notary, error)
	FindByID(ctx context.Context, id int64) (*domain.Notary, error)
	Update(ctx context.Context, notary *domain.Notary) error
	Delete(ctx context.Context, id int64) error
}

type FinancingRepository interface {
	Create(ctx context.Context, app *domain.FinancingApplication) error
	FindAll(ctx context.Context) ([]domain.FinancingApplication, error)
	FindByID(ctx context.Context, id int64) (*domain.FinancingApplication, error)
	SyncFromCBS(ctx context.Context, apps []domain.FinancingApplication) error
}

type CollateralRepository interface {
	Create(ctx context.Context, collateral *domain.Collateral) error
	FindByFinancingID(ctx context.Context, financingID int64) ([]domain.Collateral, error)
}

type LegalOrderRepository interface {
	Create(ctx context.Context, order *domain.LegalOrder) error
	FindAll(ctx context.Context) ([]domain.OrderDetail, error)
	FindByID(ctx context.Context, id int64) (*domain.OrderDetail, error)
	FindByNotaryID(ctx context.Context, notaryID int64) ([]domain.OrderDetail, error)
	FindByAssignedTo(ctx context.Context, assignedTo int64) ([]domain.OrderDetail, error)
	UpdateStatus(ctx context.Context, orderID int64, previousStatus, newStatus string, changedBy int64) error
}

type LegalOrderLogRepository interface {
	FindByOrderID(ctx context.Context, orderID int64) ([]domain.LegalOrderLog, error)
}

type LegalDocumentRepository interface {
	Create(ctx context.Context, document *domain.LegalDocument) error
	FindByOrderID(ctx context.Context, orderID int64) ([]domain.LegalDocument, error)
	FindByID(ctx context.Context, id int64) (*domain.LegalDocument, error)
	UpdateESignStatus(ctx context.Context, docID int64, status string) error
	UpdateProcessingStatus(ctx context.Context, docID int64, status string, actNumber string, notaryFee int64, processedAt time.Time) error
}
