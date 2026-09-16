package http

import (
	"errors"
	"net/http"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/repository"
	"github.com/e-notary-bprs/backend/internal/usecase"
	"github.com/e-notary-bprs/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// FinancingHandler menangani request data pembiayaan.
type FinancingHandler struct {
	financingUsecase *usecase.FinancingUcase
}

func NewFinancingHandler(financingUsecase *usecase.FinancingUcase) *FinancingHandler {
	return &FinancingHandler{financingUsecase: financingUsecase}
}

func (h *FinancingHandler) Create(c *gin.Context) {
	var f domain.FinancingApplication
	if err := c.ShouldBindJSON(&f); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.financingUsecase.Create(c.Request.Context(), &f); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, f)
}

func (h *FinancingHandler) List(c *gin.Context) {
	apps, err := h.financingUsecase.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, apps)
}

func (h *FinancingHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	app, err := h.financingUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, app)
}

func (h *FinancingHandler) SyncFromCBS(c *gin.Context) {
	var apps []domain.FinancingApplication
	if err := c.ShouldBindJSON(&apps); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.financingUsecase.SyncFromCBS(c.Request.Context(), apps); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "synced", "count": len(apps)})
}

func (h *FinancingHandler) FetchFromCBS(c *gin.Context) {
	apps, err := h.financingUsecase.FetchFromCBS(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "fetched from CBS", "count": len(apps)})
}

func (h *FinancingHandler) ValidateBPN(c *gin.Context) {
	noSertifikat := c.Query("no_sertifikat")
	if noSertifikat == "" {
		response.Error(c, http.StatusBadRequest, "no_sertifikat is required")
		return
	}
	result, err := h.financingUsecase.ValidateCollateralViaBPN(c.Request.Context(), noSertifikat)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, result)
}

func (h *FinancingHandler) Apply(c *gin.Context) {
	userIDValue, ok := c.Get("user_id")
	userID, valid := userIDValue.(int64)
	if !ok || !valid || userID == 0 {
		response.Error(c, http.StatusUnauthorized, "user context is unavailable")
		return
	}
	var f domain.FinancingApplication
	if err := c.ShouldBindJSON(&f); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.financingUsecase.ApplyForFinancing(c.Request.Context(), userID, &f); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, f)
}

func (h *FinancingHandler) ListMine(c *gin.Context) {
	userIDValue, ok := c.Get("user_id")
	userID, valid := userIDValue.(int64)
	if !ok || !valid || userID == 0 {
		response.Error(c, http.StatusUnauthorized, "user context is unavailable")
		return
	}
	apps, err := h.financingUsecase.ListMine(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, apps)
}

func (h *FinancingHandler) ValidatePegadaian(c *gin.Context) {
	noSBG := c.Query("no_sbg")
	if noSBG == "" {
		response.Error(c, http.StatusBadRequest, "no_sbg is required")
		return
	}
	result, err := h.financingUsecase.ValidateCollateralViaPegadaian(c.Request.Context(), noSBG)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, result)
}
