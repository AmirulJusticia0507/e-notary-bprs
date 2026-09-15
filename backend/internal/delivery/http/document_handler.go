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

// LegalDocumentHandler menangani request dokumen legal.
type LegalDocumentHandler struct {
	docUsecase *usecase.LegalDocumentUcase
}

func NewLegalDocumentHandler(docUsecase *usecase.LegalDocumentUcase) *LegalDocumentHandler {
	return &LegalDocumentHandler{docUsecase: docUsecase}
}

func (h *LegalDocumentHandler) Create(c *gin.Context) {
	var doc domain.LegalDocument
	if err := c.ShouldBindJSON(&doc); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.docUsecase.Create(c.Request.Context(), &doc); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, doc)
}

func (h *LegalDocumentHandler) ListByOrderID(c *gin.Context) {
	orderID, err := parseIDParam(c, "orderID")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order id")
		return
	}
	docs, err := h.docUsecase.ListByOrderID(c.Request.Context(), orderID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, docs)
}

func (h *LegalDocumentHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	doc, err := h.docUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, doc)
}

func (h *LegalDocumentHandler) UpdateESignStatus(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.docUsecase.UpdateESignStatus(c.Request.Context(), id, req.Status); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "e-sign status updated"})
}
