package delivery

import (
	"database/sql"

	"github.com/e-notary-bprs/backend/internal/config"
	deliveryhttp "github.com/e-notary-bprs/backend/internal/delivery/http"
	"github.com/e-notary-bprs/backend/internal/repository/postgres"
	"github.com/e-notary-bprs/backend/internal/usecase"
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
		AllowCredentials: true,
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
	financingUsecase := usecase.NewFinancingUcase(financingRepo)
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

	// Route terlindungi.
	authMiddleware := deliveryhttp.NewAuthMiddleware(&cfg.JWT)
	protected := v1.Group("")
	protected.Use(authMiddleware.Handler())
	{
		// Notaries
		protected.POST("/notaries", notaryHandler.Create)
		protected.GET("/notaries", notaryHandler.List)
		protected.GET("/notaries/:id", notaryHandler.GetByID)
		protected.PUT("/notaries/:id", notaryHandler.Update)
		protected.DELETE("/notaries/:id", notaryHandler.Delete)

		// Financing Applications
		protected.POST("/financings", financingHandler.Create)
		protected.GET("/financings", financingHandler.List)
		protected.GET("/financings/:id", financingHandler.GetByID)
		protected.POST("/financings/sync", financingHandler.SyncFromCBS)

		// Legal Orders
		protected.POST("/orders", orderHandler.Create)
		protected.GET("/orders", orderHandler.List)
		protected.GET("/orders/:id", orderHandler.GetByID)
		protected.GET("/orders/notary/:notaryID", orderHandler.ListByNotary)
		protected.GET("/orders/assigned/:assignedTo", orderHandler.ListByAssignedTo)
		protected.PATCH("/orders/:id/status", orderHandler.UpdateStatus)
		protected.GET("/orders/:id/logs", orderHandler.GetLogs)

		// Legal Documents
		protected.POST("/documents", docHandler.Create)
		protected.GET("/documents/order/:orderID", docHandler.ListByOrderID)
		protected.GET("/documents/:id", docHandler.GetByID)
		protected.PATCH("/documents/:id/sign", docHandler.UpdateESignStatus)
	}

	return r
}
