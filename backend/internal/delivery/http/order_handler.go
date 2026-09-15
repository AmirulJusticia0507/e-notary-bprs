package http

import (
	"net/http"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/usecase"
	"github.com/e-notary-bprs/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// LegalOrderHandler menangani request order legalitas.
type LegalOrderHandler struct {
	orderUsecase *usecase.LegalOrderUcase
}

func NewLegalOrderHandler(orderUsecase *usecase.LegalOrderUcase) *LegalOrderHandler {
	return &LegalOrderHandler{orderUsecase: orderUsecase}
}

func (h *LegalOrderHandler) Create(c *gin.Context) {
	var o domain.LegalOrder
	if err := c.ShouldBindJSON(&o); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.orderUsecase.Create(c.Request.Context(), &o); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, o)
}

func (h *LegalOrderHandler) List(c *gin.Context) {
	orders, err := h.orderUsecase.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, orders)
}

func (h *LegalOrderHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	order, err := h.orderUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, order)
}

func (h *LegalOrderHandler) ListByNotary(c *gin.Context) {
	notaryID, err := parseIDParam(c, "notaryID")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid notary id")
		return
	}
	orders, err := h.orderUsecase.ListByNotary(c.Request.Context(), notaryID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, orders)
}

func (h *LegalOrderHandler) ListByAssignedTo(c *gin.Context) {
	assignedTo, err := parseIDParam(c, "assignedTo")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid assigned id")
		return
	}
	orders, err := h.orderUsecase.ListByAssignedTo(c.Request.Context(), assignedTo)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, orders)
}

func (h *LegalOrderHandler) UpdateStatus(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Status    string `json:"status" binding:"required"`
		ChangedBy int64  `json:"changed_by" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.orderUsecase.UpdateStatus(c.Request.Context(), id, req.Status, req.ChangedBy); err != nil {
		switch err {
		case usecase.ErrOrderNotFound:
			response.Error(c, http.StatusNotFound, err.Error())
		case usecase.ErrInvalidStatus:
			response.Error(c, http.StatusBadRequest, err.Error())
		case usecase.ErrSLAOverdue:
			response.Error(c, http.StatusForbidden, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "status updated"})
}

func (h *LegalOrderHandler) GetLogs(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	logs, err := h.orderUsecase.GetLogs(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, logs)
}
