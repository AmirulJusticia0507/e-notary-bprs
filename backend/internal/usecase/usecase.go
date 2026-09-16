package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/repository"
	"github.com/e-notary-bprs/backend/pkg/auth"
	cbs "github.com/e-notary-bprs/backend/pkg/cbs"
	bpn "github.com/e-notary-bprs/backend/pkg/bpn"
	pegadaian "github.com/e-notary-bprs/backend/pkg/pegadaian"
	emeterai "github.com/e-notary-bprs/backend/pkg/emeterai"
	esign "github.com/e-notary-bprs/backend/pkg/esign"
)

const (
	// MaxFailedLogins: akun dikunci setelah 5x salah password beruntun.
	MaxFailedLogins = 5
	// LockoutDuration: lama kunci akun.
	LockoutDuration = 15 * time.Minute
)

var (
	ErrOrderNotFound = errors.New("legal order not found")
	ErrInvalidStatus = errors.New("invalid status transition")
	ErrSLAOverdue    = errors.New("SLA deadline exceeded")
	ErrInvalidRole   = errors.New("invalid user role")
	ErrAccountLocked = errors.New("account locked")
)

type AuthUcase struct {
	userRepo  repository.UserRepository
	resetRepo repository.PasswordResetRepository
}

func NewAuthUcase(userRepo repository.UserRepository, resetRepo repository.PasswordResetRepository) *AuthUcase {
	return &AuthUcase{userRepo: userRepo, resetRepo: resetRepo}
}

// ResetTokenTTL: masa berlaku kode reset password.
const ResetTokenTTL = time.Hour

// MaxPhotoChars: batas data-URL foto profil (~500KB, hasil downscale frontend jauh di bawah ini).
const MaxPhotoChars = 700000

func (u *AuthUcase) Register(ctx context.Context, fullName, email, password, role, photoURL string) error {
	if role == "" {
		role = "legal_officer"
	}
	if role != "legal_officer" && role != "admin" && role != "notary" && role != "nasabah" {
		return ErrInvalidRole
	}
	if len(photoURL) > MaxPhotoChars {
		return errors.New("foto terlalu besar, maksimal ~500KB")
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
		PhotoURL:     photoURL,
	})
}

// UpdatePhoto mengganti foto profil user yang sedang login.
func (u *AuthUcase) UpdatePhoto(ctx context.Context, userID int64, photoURL string) error {
	if len(photoURL) > MaxPhotoChars {
		return errors.New("foto terlalu besar, maksimal ~500KB")
	}
	return u.userRepo.UpdatePhoto(ctx, userID, photoURL)
}

func (u *AuthUcase) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if !user.IsActive {
		return nil, errors.New("invalid email or password")
	}
	if user.LockedUntil != nil {
		if remain := time.Until(*user.LockedUntil); remain > 0 {
			return nil, fmt.Errorf("%w: coba lagi dalam %d menit", ErrAccountLocked, int(remain.Minutes())+1)
		}
	}
	if !auth.CheckPassword(user.PasswordHash, password) {
		attempts := user.FailedLoginAttempts + 1
		var lockedUntil *time.Time
		if attempts >= MaxFailedLogins {
			t := time.Now().Add(LockoutDuration)
			lockedUntil = &t
			attempts = 0
		}
		_ = u.userRepo.RecordFailedLogin(ctx, user.ID, attempts, lockedUntil)
		if lockedUntil != nil {
			return nil, fmt.Errorf("%w: 5x salah password, coba lagi dalam 15 menit", ErrAccountLocked)
		}
		return nil, errors.New("invalid email or password")
	}
	_ = u.userRepo.ResetLoginAttempts(ctx, user.ID)
	return user, nil
}

// newResetToken membuat token acak + hash SHA256-nya untuk disimpan.
func newResetToken() (token, tokenHash string, err error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", err
	}
	token = hex.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	return token, hex.EncodeToString(sum[:]), nil
}

func hashResetToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// RequestPasswordReset membuat kode reset (berlaku 1 jam). Selalu sukses
// tanpa membocorkan apakah email terdaftar (anti enumeration).
func (u *AuthUcase) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil || !user.IsActive {
		return nil
	}
	_ = u.resetRepo.DeleteByUserID(ctx, user.ID)
	_, tokenHash, err := newResetToken()
	if err != nil {
		return err
	}
	_ = u.resetRepo.DeleteExpired(ctx)
	return u.resetRepo.Create(ctx, &domain.PasswordReset{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(ResetTokenTTL),
	})
}

// AdminCreateResetToken menerbitkan kode reset baru untuk user dan
// mengembalikan plaintext-nya SEKALI saja untuk di-relay admin ke user.
func (u *AuthUcase) AdminCreateResetToken(ctx context.Context, userID int64) (string, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return "", err
	}
	_ = u.resetRepo.DeleteByUserID(ctx, userID)
	token, tokenHash, err := newResetToken()
	if err != nil {
		return "", err
	}
	if err := u.resetRepo.Create(ctx, &domain.PasswordReset{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(ResetTokenTTL),
	}); err != nil {
		return "", err
	}
	return token, nil
}

// ResetPassword menukar kode reset yang valid dengan password baru.
func (u *AuthUcase) ResetPassword(ctx context.Context, token, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("password minimal 6 karakter")
	}
	reset, err := u.resetRepo.FindValidByTokenHash(ctx, hashResetToken(token))
	if err != nil {
		return errors.New("kode reset tidak valid atau kedaluwarsa")
	}
	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := u.userRepo.UpdatePassword(ctx, reset.UserID, hash); err != nil {
		return err
	}
	_ = u.userRepo.ResetLoginAttempts(ctx, reset.UserID)
	_ = u.resetRepo.DeleteByID(ctx, reset.ID)
	return nil
}

// ChangePassword untuk user yang masih bisa login.
func (u *AuthUcase) ChangePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("password minimal 6 karakter")
	}
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if !auth.CheckPassword(user.PasswordHash, currentPassword) {
		return errors.New("password lama salah")
	}
	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return u.userRepo.UpdatePassword(ctx, userID, hash)
}

// ListPendingResets untuk relay admin (tanpa token, hanya siapa + expiry).
func (u *AuthUcase) ListPendingResets(ctx context.Context) ([]domain.PasswordReset, error) {
	return u.resetRepo.ListPending(ctx)
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
	financingRepo   repository.FinancingRepository
	cbsClient       *cbs.Client
	bpnClient       *bpn.Client
	pegadaianClient *pegadaian.Client
}

func NewFinancingUcase(financingRepo repository.FinancingRepository, cbsClient *cbs.Client, bpnClient *bpn.Client, pegadaianClient *pegadaian.Client) *FinancingUcase {
	return &FinancingUcase{financingRepo: financingRepo, cbsClient: cbsClient, bpnClient: bpnClient, pegadaianClient: pegadaianClient}
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

// ApplyForFinancing mencatat pengajuan pembiayaan dari akun nasabah sendiri.
func (u *FinancingUcase) ApplyForFinancing(ctx context.Context, userID int64, app *domain.FinancingApplication) error {
	app.UserID = &userID
	app.Status = "pending"
	return u.financingRepo.Create(ctx, app)
}

// ListMine mengembalikan pengajuan milik satu nasabah.
func (u *FinancingUcase) ListMine(ctx context.Context, userID int64) ([]domain.FinancingApplication, error) {
	return u.financingRepo.FindByUserID(ctx, userID)
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

func (u *FinancingUcase) ValidateCollateralViaBPN(ctx context.Context, noSertifikat string) (*bpn.ValidasiResult, error) {
	if u.bpnClient == nil {
		return nil, errors.New("BPN client not configured")
	}
	return u.bpnClient.ValidateSertifikat(ctx, noSertifikat)
}

func (u *FinancingUcase) ValidateCollateralViaPegadaian(ctx context.Context, noSBG string) (*pegadaian.ValidasiResult, error) {
	if u.pegadaianClient == nil {
		return nil, errors.New("pegadaian client not configured")
	}
	return u.pegadaianClient.ValidateGadai(ctx, noSBG)
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
	docRepo        repository.LegalDocumentRepository
	emeteraiClient *emeterai.Client
	esignClient    *esign.Client
}

func NewLegalDocumentUcase(docRepo repository.LegalDocumentRepository, emeteraiClient *emeterai.Client, esignClient *esign.Client) *LegalDocumentUcase {
	return &LegalDocumentUcase{docRepo: docRepo, emeteraiClient: emeteraiClient, esignClient: esignClient}
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

// StampMeterai membeli e-Meterai ke distributor resmi lalu menyimpan SN-nya ke dokumen.
func (u *LegalDocumentUcase) StampMeterai(ctx context.Context, docID int64, purchaser string) (*emeterai.PurchaseResponse, error) {
	if u.emeteraiClient == nil || !u.emeteraiClient.IsConfigured() {
		return nil, errors.New("e-Meterai provider is not configured (isi EMETERAI_BASE_URL dan EMETERAI_API_KEY)")
	}
	doc, err := u.docRepo.FindByID(ctx, docID)
	if err != nil {
		return nil, err
	}
	resp, err := u.emeteraiClient.Purchase(ctx, emeterai.PurchaseRequest{
		DocumentID: doc.FileURL,
		Purchaser:  purchaser,
		Amount:     1,
	})
	if err != nil {
		return nil, err
	}
	if err := u.docRepo.UpdateMeteraiSN(ctx, docID, resp.MeteraiSN); err != nil {
		return nil, err
	}
	return resp, nil
}

// RequestSign mengajukan TTE tersertifikasi ke PSrE lalu menandai status e-sign dokumen.
func (u *LegalDocumentUcase) RequestSign(ctx context.Context, docID int64, signerName, signerEmail string) (*esign.SignResponse, error) {
	if u.esignClient == nil || !u.esignClient.IsConfigured() {
		return nil, errors.New("e-Sign provider is not configured (isi ESIGN_BASE_URL dan ESIGN_API_KEY)")
	}
	if _, err := u.docRepo.FindByID(ctx, docID); err != nil {
		return nil, err
	}
	resp, err := u.esignClient.Sign(ctx, esign.SignRequest{
		DocumentID:  strconv.FormatInt(docID, 10),
		SignerName:  signerName,
		SignerEmail: signerEmail,
	})
	if err != nil {
		return nil, err
	}
	status := "signed"
	if resp.Status != "" {
		status = resp.Status
	}
	if err := u.docRepo.UpdateESignStatus(ctx, docID, status); err != nil {
		return nil, err
	}
	return resp, nil
}
