package delivery

import (
	"database/sql"

	"github.com/e-notary-bprs/backend/internal/config"
	deliveryhttp "github.com/e-notary-bprs/backend/internal/delivery/http"
	"github.com/e-notary-bprs/backend/internal/repository/postgres"
	"github.com/e-notary-bprs/backend/internal/usecase"
	"github.com/e-notary-bprs/backend/pkg/bpn"
	"github.com/e-notary-bprs/backend/pkg/cbs"
	"github.com/e-notary-bprs/backend/pkg/pegadaian"
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
	authUsecase := usecase.NewAuthUcase(userRepo)
	notaryUsecase := usecase.NewNotaryUcase(notaryRepo)
	cbsClient := cbs.NewClient(cfg.CBS.BaseURL, cfg.CBS.APIKey, cfg.CBS.ClientCode)
	bpnClient := bpn.NewClient(cfg.BPN.BaseURL, cfg.BPN.APIKey, cfg.BPN.SecretKey)
	pegadaianClient := pegadaian.NewClient(cfg.Pegadaian.BaseURL, cfg.Pegadaian.APIKey, cfg.Pegadaian.PartnerCode)
	financingUsecase := usecase.NewFinancingUcase(financingRepo, cbsClient, bpnClient, pegadaianClient)
	orderUsecase := usecase.NewLegalOrderUcase(orderRepo, logRepo)
	docUsecase := usecase.NewLegalDocumentUcase(docRepo)

	// Inisialisasi handlers.
	authHandler := deliveryhttp.NewAuthHandler(authUsecase, &cfg.JWT)
	notaryHandler := deliveryhttp.NewNotaryHandler(notaryUsecase)
	financingHandler := deliveryhttp.NewFinancingHandler(financingUsecase)
	orderHandler := deliveryhttp.NewLegalOrderHandler(orderUsecase)
	docHandler := deliveryhttp.NewLegalDocumentHandler(docUsecase)

	v1 := r.Group("/api/v1")
	// Route publik.
	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)

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

	// Notaries: master hanya admin yang ubah, semua role boleh lihat.
	adminOnly.POST("/notaries", notaryHandler.Create)
	adminOnly.PUT("/notaries/:id", notaryHandler.Update)
	adminOnly.DELETE("/notaries/:id", notaryHandler.Delete)
	allStaff.GET("/notaries", notaryHandler.List)
	allStaff.GET("/notaries/:id", notaryHandler.GetByID)

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
	allStaff.GET("/documents/order/:orderID", docHandler.ListByOrderID)
	allStaff.PATCH("/documents/:id/processing", docHandler.UpdateProcessingStatus)
	allStaff.GET("/documents/:id", docHandler.GetByID)

	return r
}
