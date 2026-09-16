package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/repository"
	"github.com/e-notary-bprs/backend/pkg/auth"
	cbs "github.com/e-notary-bprs/backend/pkg/cbs"
)

var (
	ErrOrderNotFound = errors.New("legal order not found")
	ErrInvalidStatus = errors.New("invalid status transition")
	ErrSLAOverdue    = errors.New("SLA deadline exceeded")
	ErrInvalidRole   = errors.New("invalid user role")
)

type AuthUcase struct {
	userRepo repository.UserRepository
}

func NewAuthUcase(userRepo repository.UserRepository) *AuthUcase {
	return &AuthUcase{userRepo: userRepo}
}

func (u *AuthUcase) Register(ctx context.Context, fullName, email, password, role string) error {
	if role == "" {
		role = "legal_officer"
	}
	if role != "legal_officer" && role != "admin" && role != "notary" {
		return ErrInvalidRole
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	return u.userRepo.Create(ctx, &domain.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
		IsActive:     true,
	})
}

func (u *AuthUcase) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if !user.IsActive || !auth.CheckPassword(user.PasswordHash, password) {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}

type NotaryUcase struct {
	notaryRepo repository.NotaryRepository
}

func NewNotaryUcase(notaryRepo repository.NotaryRepository) *NotaryUcase {
	return &NotaryUcase{notaryRepo: notaryRepo}
}

func (u *NotaryUcase) Create(ctx context.Context, notary *domain.Notary) error {
	return u.notaryRepo.Create(ctx, notary)
}

func (u *NotaryUcase) List(ctx context.Context) ([]domain.Notary, error) {
	return u.notaryRepo.FindAll(ctx)
}

func (u *NotaryUcase) GetByID(ctx context.Context, id int64) (*domain.Notary, error) {
	return u.notaryRepo.FindByID(ctx, id)
}

func (u *NotaryUcase) Update(ctx context.Context, notary *domain.Notary) error {
	return u.notaryRepo.Update(ctx, notary)
}

func (u *NotaryUcase) Delete(ctx context.Context, id int64) error {
	return u.notaryRepo.Delete(ctx, id)
}

type FinancingUcase struct {
	financingRepo repository.FinancingRepository
	cbsClient     *cbs.Client
}

func NewFinancingUcase(financingRepo repository.FinancingRepository, cbsClient *cbs.Client) *FinancingUcase {
	return &FinancingUcase{financingRepo: financingRepo, cbsClient: cbsClient}
}

func (u *FinancingUcase) Create(ctx context.Context, app *domain.FinancingApplication) error {
	return u.financingRepo.Create(ctx, app)
}

func (u *FinancingUcase) List(ctx context.Context) ([]domain.FinancingApplication, error) {
	return u.financingRepo.FindAll(ctx)
}

func (u *FinancingUcase) GetByID(ctx context.Context, id int64) (*domain.FinancingApplication, error) {
	return u.financingRepo.FindByID(ctx, id)
}

func (u *FinancingUcase) SyncFromCBS(ctx context.Context, apps []domain.FinancingApplication) error {
	return u.financingRepo.SyncFromCBS(ctx, apps)
}

func (u *FinancingUcase) FetchFromCBS(ctx context.Context) ([]domain.FinancingApplication, error) {
	if u.cbsClient == nil {
		return nil, errors.New("CBS client not configured")
	}

	financings, err := u.cbsClient.GetFinancings(ctx)
	if err != nil {
		return nil, err
	}

	// Convert CBS data to domain models
	apps := make([]domain.FinancingApplication, 0)
	for _, f := range financings {
		app := domain.FinancingApplication{
			CustomerName:      f.CustomerName,
			CustomerNIK:       f.CustomerNIK,
			FinancingAmount:   int64(f.FinancingAmount),
			CollateralType:    f.CollateralType,
			CollateralDetails: f.CollateralDetails,
			Status:            f.Status,
			SyncedAt:          time.Now(),
		}
		apps = append(apps, app)
	}

	return apps, nil
}

type LegalOrderUcase struct {
	orderRepo repository.LegalOrderRepository
	logRepo   repository.LegalOrderLogRepository
}

func NewLegalOrderUcase(orderRepo repository.LegalOrderRepository, logRepo repository.LegalOrderLogRepository) *LegalOrderUcase {
	return &LegalOrderUcase{orderRepo: orderRepo, logRepo: logRepo}
}

func (u *LegalOrderUcase) Create(ctx context.Context, order *domain.LegalOrder) error {
	order.Status = "pending"
	order.SLADeadline = time.Now().Add(14 * 24 * time.Hour)
	return u.orderRepo.Create(ctx, order)
}

func (u *LegalOrderUcase) List(ctx context.Context) ([]domain.OrderDetail, error) {
	return u.orderRepo.FindAll(ctx)
}

func (u *LegalOrderUcase) GetByID(ctx context.Context, id int64) (*domain.OrderDetail, error) {
	return u.orderRepo.FindByID(ctx, id)
}

func (u *LegalOrderUcase) ListByNotary(ctx context.Context, notaryID int64) ([]domain.OrderDetail, error) {
	return u.orderRepo.FindByNotaryID(ctx, notaryID)
}

func (u *LegalOrderUcase) ListByAssignedTo(ctx context.Context, assignedTo int64) ([]domain.OrderDetail, error) {
	return u.orderRepo.FindByAssignedTo(ctx, assignedTo)
}

func (u *LegalOrderUcase) UpdateStatus(ctx context.Context, orderID int64, newStatus string, changedBy int64) error {
	detail, err := u.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrOrderNotFound
		}
		return err
	}
	if !isValidTransition(detail.Status, newStatus) {
		return ErrInvalidStatus
	}
	if time.Now().After(detail.SLADeadline) && newStatus != "completed" {
		return ErrSLAOverdue
	}
	return u.orderRepo.UpdateStatus(ctx, orderID, detail.Status, newStatus, changedBy)
}

func (u *LegalOrderUcase) GetLogs(ctx context.Context, orderID int64) ([]domain.LegalOrderLog, error) {
	return u.logRepo.FindByOrderID(ctx, orderID)
}

func isValidTransition(from, to string) bool {
	transitions := map[string][]string{
		"pending":     {"in_progress"},
		"in_progress": {"completed", "rejected"},
	}
	for _, status := range transitions[from] {
		if status == to {
			return true
		}
	}
	return false
}

type LegalDocumentUcase struct {
	docRepo repository.LegalDocumentRepository
}

func NewLegalDocumentUcase(docRepo repository.LegalDocumentRepository) *LegalDocumentUcase {
	return &LegalDocumentUcase{docRepo: docRepo}
}

func (u *LegalDocumentUcase) Create(ctx context.Context, document *domain.LegalDocument) error {
	return u.docRepo.Create(ctx, document)
}

func (u *LegalDocumentUcase) ListByOrderID(ctx context.Context, orderID int64) ([]domain.LegalDocument, error) {
	return u.docRepo.FindByOrderID(ctx, orderID)
}

func (u *LegalDocumentUcase) GetByID(ctx context.Context, id int64) (*domain.LegalDocument, error) {
	return u.docRepo.FindByID(ctx, id)
}

func (u *LegalDocumentUcase) UpdateESignStatus(ctx context.Context, docID int64, status string) error {
	return u.docRepo.UpdateESignStatus(ctx, docID, status)
}

func (u *LegalDocumentUcase) UpdateProcessingStatus(ctx context.Context, docID int64, status string, actNumber string, notaryFee int64, processedAt time.Time) error {
	return u.docRepo.UpdateProcessingStatus(ctx, docID, status, actNumber, notaryFee, processedAt)
}
