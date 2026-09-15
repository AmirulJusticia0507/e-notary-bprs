package http

import (
	"net/http"

	"github.com/e-notary-bprs/backend/internal/domain"
	"github.com/e-notary-bprs/backend/internal/usecase"
	"github.com/e-notary-bprs/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// NotaryHandler menangani request data notaris.
type NotaryHandler struct {
	notaryUsecase *usecase.NotaryUcase
}

func NewNotaryHandler(notaryUsecase *usecase.NotaryUcase) *NotaryHandler {
	return &NotaryHandler{notaryUsecase: notaryUsecase}
}

func (h *NotaryHandler) Create(c *gin.Context) {
	var n domain.Notary
	if err := c.ShouldBindJSON(&n); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.notaryUsecase.Create(c.Request.Context(), &n); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, n)
}

func (h *NotaryHandler) List(c *gin.Context) {
	notaries, err := h.notaryUsecase.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, notaries)
}

func (h *NotaryHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	notary, err := h.notaryUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, notary)
}

func (h *NotaryHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	var n domain.Notary
	if err := c.ShouldBindJSON(&n); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	n.ID = id
	if err := h.notaryUsecase.Update(c.Request.Context(), &n); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, n)
}

func (h *NotaryHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.notaryUsecase.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "notary deleted"})
}
