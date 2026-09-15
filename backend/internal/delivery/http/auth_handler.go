package http

import (
	"net/http"

	"github.com/e-notary-bprs/backend/internal/config"
	"github.com/e-notary-bprs/backend/internal/usecase"
	"github.com/e-notary-bprs/backend/pkg/auth"
	"github.com/e-notary-bprs/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthHandler menangani request autentikasi.
type AuthHandler struct {
	authUsecase *usecase.AuthUcase
	jwtCfg      *config.JWTConfig
}

func NewAuthHandler(authUsecase *usecase.AuthUcase, jwtCfg *config.JWTConfig) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase, jwtCfg: jwtCfg}
}

// Register handles user registration.
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.authUsecase.Register(c.Request.Context(), req.FullName, req.Email, req.Password, req.Role); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, gin.H{"email": req.Email})
}

// Login handles user authentication.
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.authUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := auth.GenerateToken(user, h.jwtCfg.Secret, h.jwtCfg.Expiry)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to generate token")
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "email": user.Email, "role": user.Role}})
}
