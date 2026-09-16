package delivery

import (
	"context"
	"database/sql"
	"log"

	"github.com/e-notary-bprs/backend/internal/config"
	deliveryhttp "github.com/e-notary-bprs/backend/internal/delivery/http"
	"github.com/e-notary-bprs/backend/internal/repository/postgres"
	"github.com/e-notary-bprs/backend/internal/usecase"
	"github.com/e-notary-bprs/backend/pkg/bpn"
	"github.com/e-notary-bprs/backend/pkg/cbs"
	"github.com/e-notary-bprs/backend/pkg/emeterai"
	"github.com/e-notary-bprs/backend/pkg/esign"
	"github.com/e-notary-bprs/backend/pkg/pegadaian"
	"github.com/e-notary-bprs/backend/pkg/sso"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter menginisialisasi semua handler dan mengembalikan gin Engine.
func NewRouter(cfg *config.Config, db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	// Inisialisasi repositories.
	userRepo := postgres.NewUserRepository(db)
	notaryRepo := postgres.NewNotaryRepository(db)
	financingRepo := postgres.NewFinancingRepository(db)
	orderRepo := postgres.NewLegalOrderRepository(db)
	logRepo := postgres.NewLegalOrderLogRepository(db)
	docRepo := postgres.NewLegalDocumentRepository(db)

	// Inisialisasi usecases.
	resetRepo := postgres.NewPasswordResetRepository(db)
	var ssoClient *sso.Client
	if cfg.Keycloak.Enabled() {
		var err error
		ssoClient, err = sso.NewClient(context.Background(), cfg.Keycloak.Issuer, cfg.Keycloak.InternalURL, cfg.Keycloak.ClientID, cfg.Keycloak.ClientSecret, cfg.Keycloak.RedirectURL)
		if err != nil {
			log.Printf("keycloak dinonaktifkan: %v", err)
			ssoClient = nil
		}
	}
	authUsecase := usecase.NewAuthUcase(userRepo, resetRepo, ssoClient)
	notaryUsecase := usecase.NewNotaryUcase(notaryRepo)
	cbsClient := cbs.NewClient(cfg.CBS.BaseURL, cfg.CBS.APIKey, cfg.CBS.ClientCode)
	bpnClient := bpn.NewClient(cfg.BPN.BaseURL, cfg.BPN.APIKey, cfg.BPN.SecretKey)
	pegadaianClient := pegadaian.NewClient(cfg.Pegadaian.BaseURL, cfg.Pegadaian.APIKey, cfg.Pegadaian.PartnerCode)
	financingUsecase := usecase.NewFinancingUcase(financingRepo, cbsClient, bpnClient, pegadaianClient)
	orderUsecase := usecase.NewLegalOrderUcase(orderRepo, logRepo)
	emeteraiClient := emeterai.NewClient(cfg.EMeterai.BaseURL, cfg.EMeterai.APIKey)
	esignClient := esign.NewClient(cfg.ESign.BaseURL, cfg.ESign.APIKey)
	docUsecase := usecase.NewLegalDocumentUcase(docRepo, emeteraiClient, esignClient)

	// Inisialisasi handlers.
	authHandler := deliveryhttp.NewAuthHandler(authUsecase, &cfg.JWT, cfg.Keycloak.FrontendURL)
	notaryHandler := deliveryhttp.NewNotaryHandler(notaryUsecase)
	financingHandler := deliveryhttp.NewFinancingHandler(financingUsecase)
	orderHandler := deliveryhttp.NewLegalOrderHandler(orderUsecase)
	docHandler := deliveryhttp.NewLegalDocumentHandler(docUsecase)

	// Health check publik (tanpa auth) agar root URL tidak 404.
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "e-notary-bprs-api"})
	})

	v1 := r.Group("/api/v1")
	// Route publik.
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "e-notary-bprs-api"})
	})
	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)
	v1.POST("/auth/forgot-password", authHandler.ForgotPassword)
	v1.POST("/auth/reset-password", authHandler.ResetPassword)
	v1.GET("/auth/sso/login", authHandler.SSOStart)
	v1.GET("/auth/sso/callback", authHandler.SSOCallback)

	// Route terlindungi (wajib JWT).
	authMiddleware := deliveryhttp.NewAuthMiddleware(&cfg.JWT)
	protected := v1.Group("")
	protected.Use(authMiddleware.Handler())

	// Grup role:
	// - adminOnly: kelola master (notaris create/update/delete)
	// - legalTeam: admin + legal_officer (buat order, sync CBS, kelola dokumen)
	// - allStaff: admin + legal_officer + notary (lihat data + update processing notaris)
	adminOnly := protected.Group("")
	adminOnly.Use(authMiddleware.RequireRole("admin"))

	legalTeam := protected.Group("")
	legalTeam.Use(authMiddleware.RequireRole("admin", "legal_officer"))

	allStaff := protected.Group("")
	allStaff.Use(authMiddleware.RequireRole("admin", "legal_officer", "notary"))

	// Auth mandiri: semua role yang login boleh ganti password sendiri.
	allAuthed := protected.Group("")
	allAuthed.Use(authMiddleware.RequireRole("admin", "legal_officer", "notary", "nasabah"))
	allAuthed.PATCH("/auth/change-password", authHandler.ChangePassword)
	allAuthed.PATCH("/auth/photo", authHandler.UpdatePhoto)

	// Admin: relay kode reset password ke user (via WA/telepon).
	adminOnly.GET("/users/password-resets", authHandler.ListPendingResets)
	adminOnly.POST("/users/:userID/reset-token", authHandler.AdminCreateResetToken)

	// Notaries: master hanya admin yang ubah, semua role boleh lihat.
	adminOnly.POST("/notaries", notaryHandler.Create)
	adminOnly.PUT("/notaries/:id", notaryHandler.Update)
	adminOnly.DELETE("/notaries/:id", notaryHandler.Delete)
	allStaff.GET("/notaries", notaryHandler.List)
	allStaff.GET("/notaries/:id", notaryHandler.GetByID)

	// Nasabah: ajukan + pantau pengajuan sendiri.
	nasabahOnly := protected.Group("")
	nasabahOnly.Use(authMiddleware.RequireRole("nasabah"))
	nasabahOnly.POST("/financings/apply", financingHandler.Apply)
	nasabahOnly.GET("/financings/mine", financingHandler.ListMine)

	// Financing Applications: legal yang input/sync, semua role boleh lihat + validasi.
	legalTeam.POST("/financings", financingHandler.Create)
	legalTeam.POST("/financings/sync", financingHandler.SyncFromCBS)
	allStaff.GET("/financings", financingHandler.List)
	allStaff.GET("/financings/fetch-cbs", financingHandler.FetchFromCBS)
	allStaff.GET("/financings/validate-bpn", financingHandler.ValidateBPN)
	allStaff.GET("/financings/validate-pegadaian", financingHandler.ValidatePegadaian)
	allStaff.GET("/financings/:id", financingHandler.GetByID)

	// Legal Orders: legal yang buat/ubah status, semua role boleh lihat.
	legalTeam.POST("/orders", orderHandler.Create)
	legalTeam.PATCH("/orders/:id/status", orderHandler.UpdateStatus)
	allStaff.GET("/orders", orderHandler.List)
	allStaff.GET("/orders/notary/:notaryID", orderHandler.ListByNotary)
	allStaff.GET("/orders/assigned/:assignedTo", orderHandler.ListByAssignedTo)
	allStaff.GET("/orders/:id/logs", orderHandler.GetLogs)
	allStaff.GET("/orders/:id", orderHandler.GetByID)

	// Legal Documents: legal yang upload/sign, notaris update processing + semua boleh baca.
	legalTeam.POST("/documents", docHandler.Create)
	legalTeam.PATCH("/documents/:id/sign", docHandler.UpdateESignStatus)
	legalTeam.POST("/documents/:id/stamp", docHandler.StampMeterai)
	legalTeam.POST("/documents/:id/sign-request", docHandler.RequestSign)
	allStaff.GET("/documents/order/:orderID", docHandler.ListByOrderID)
	allStaff.PATCH("/documents/:id/processing", docHandler.UpdateProcessingStatus)
	allStaff.GET("/documents/:id", docHandler.GetByID)

	return r
}
