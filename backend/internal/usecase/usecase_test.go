package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/repository"
)

func futureTime() time.Time { return time.Now().Add(24 * time.Hour) }

type stubOrderRepo struct {
	detail *domain.OrderDetail
	orders []domain.OrderDetail
	status map[int64]string
}

func (s *stubOrderRepo) Create(_ context.Context, order *domain.LegalOrder) error {
	order.ID = 1
	return nil
}

func (s *stubOrderRepo) FindAll(_ context.Context) ([]domain.OrderDetail, error) {
	return s.orders, nil
}

func (s *stubOrderRepo) FindByID(_ context.Context, _ int64) (*domain.OrderDetail, error) {
	if s.detail == nil {
		return nil, repository.ErrNotFound
	}
	return s.detail, nil
}

func (s *stubOrderRepo) FindByNotaryID(_ context.Context, _ int64) ([]domain.OrderDetail, error) {
	return s.orders, nil
}

func (s *stubOrderRepo) FindByAssignedTo(_ context.Context, _ int64) ([]domain.OrderDetail, error) {
	return s.orders, nil
}

func (s *stubOrderRepo) UpdateStatus(_ context.Context, orderID int64, previousStatus, newStatus string, _ int64) error {
	if s.status == nil {
		s.status = map[int64]string{}
	}
	if s.detail != nil && previousStatus != s.detail.Status {
		return repository.ErrNotFound
	}
	s.status[orderID] = newStatus
	return nil
}

type stubLogRepo struct{}

func (s *stubLogRepo) FindByOrderID(_ context.Context, _ int64) ([]domain.LegalOrderLog, error) {
	return []domain.LegalOrderLog{}, nil
}

func newOrderUcase(detail *domain.OrderDetail) (*LegalOrderUcase, *stubOrderRepo) {
	repo := &stubOrderRepo{detail: detail}
	return NewLegalOrderUcase(repo, &stubLogRepo{}), repo
}

func TestIsValidTransition(t *testing.T) {
	cases := []struct {
		from, to string
		want     bool
	}{
		{"pending", "in_progress", true},
		{"pending", "completed", false},
		{"in_progress", "completed", true},
		{"in_progress", "rejected", true},
		{"in_progress", "pending", false},
		{"completed", "in_progress", false},
		{"rejected", "completed", false},
	}
	for _, tc := range cases {
		if got := isValidTransition(tc.from, tc.to); got != tc.want {
			t.Errorf("isValidTransition(%q,%q) = %v, want %v", tc.from, tc.to, got, tc.want)
		}
	}
}

func TestUpdateStatusInvalidTransition(t *testing.T) {
	uc, _ := newOrderUcase(&domain.OrderDetail{
		ID: 1, OrderNumber: "ORD-1", Status: "pending",
		SLADeadline: futureTime(),
	})
	if err := uc.UpdateStatus(context.Background(), 1, "completed", 7); err != ErrInvalidStatus {
		t.Errorf("UpdateStatus() = %v, want ErrInvalidStatus", err)
	}
}

func TestUpdateStatusValid(t *testing.T) {
	uc, repo := newOrderUcase(&domain.OrderDetail{
		ID: 1, OrderNumber: "ORD-1", Status: "pending",
		SLADeadline: futureTime(),
	})
	if err := uc.UpdateStatus(context.Background(), 1, "in_progress", 7); err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}
	if repo.status[1] != "in_progress" {
		t.Errorf("status = %q, want in_progress", repo.status[1])
	}
}

func TestUpdateStatusNotFound(t *testing.T) {
	uc, _ := newOrderUcase(nil)
	if err := uc.UpdateStatus(context.Background(), 99, "in_progress", 7); err != ErrOrderNotFound {
		t.Errorf("UpdateStatus() = %v, want ErrOrderNotFound", err)
	}
}
