package http

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

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
	frontendURL string
}

func NewAuthHandler(authUsecase *usecase.AuthUcase, jwtCfg *config.JWTConfig, frontendURL string) *AuthHandler {
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	return &AuthHandler{authUsecase: authUsecase, jwtCfg: jwtCfg, frontendURL: frontendURL}
}

// Register adalah signup publik khusus nasabah (role selalu nasabah).
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
		PhotoURL string `json:"photo_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	req.Role = "nasabah"
	if err := h.authUsecase.Register(c.Request.Context(), req.FullName, req.Email, req.Password, req.Role, req.PhotoURL); err != nil {
		if errors.Is(err, usecase.ErrInvalidRole) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
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
		if errors.Is(err, usecase.ErrAccountLocked) {
			response.Error(c, http.StatusLocked, err.Error())
			return
		}
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := auth.GenerateToken(user, h.jwtCfg.Secret, h.jwtCfg.Expiry)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to generate token")
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "full_name": user.FullName, "email": user.Email, "role": user.Role, "photo_url": user.PhotoURL}})
}

// UpdatePhoto mengganti foto profil user yang sedang login.
func (h *AuthHandler) UpdatePhoto(c *gin.Context) {
	userIDValue, ok := c.Get("user_id")
	userID, valid := userIDValue.(int64)
	if !ok || !valid || userID == 0 {
		response.Error(c, http.StatusUnauthorized, "user context is unavailable")
		return
	}
	var req struct {
		PhotoURL string `json:"photo_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.authUsecase.UpdatePhoto(c.Request.Context(), userID, req.PhotoURL); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "foto profil diperbarui", "photo_url": req.PhotoURL})
}

// ForgotPassword membuat kode reset (berlaku 1 jam). Selalu 200 agar
// tidak membocorkan email mana yang terdaftar.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.authUsecase.RequestPasswordReset(c.Request.Context(), req.Email); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "Jika email terdaftar, kode reset telah dibuat. Hubungi admin/CS BPRS untuk mendapatkan kode tersebut."})
}

// ResetPassword menukar kode reset dengan password baru.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.authUsecase.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "password berhasil direset, silakan masuk"})
}

// ChangePassword untuk user yang sedang login.
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDValue, ok := c.Get("user_id")
	userID, valid := userIDValue.(int64)
	if !ok || !valid || userID == 0 {
		response.Error(c, http.StatusUnauthorized, "user context is unavailable")
		return
	}
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.authUsecase.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "password berhasil diganti"})
}

// ListPendingResets (admin): daftar permintaan reset untuk di-relay ke user.
func (h *AuthHandler) ListPendingResets(c *gin.Context) {
	resets, err := h.authUsecase.ListPendingResets(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, resets)
}

// AdminCreateResetToken (admin): terbitkan kode reset untuk user, plaintext
// dikembalikan SEKALI untuk di-relay (mis. via WA) ke user bersangkutan.
func (h *AuthHandler) AdminCreateResetToken(c *gin.Context) {
	userID, err := parseIDParam(c, "userID")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}
	token, err := h.authUsecase.AdminCreateResetToken(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"reset_token": token, "expires_in": "1h"})
}

// SSOStart redirect browser ke halaman login Keycloak.
func (h *AuthHandler) SSOStart(c *gin.Context) {
	url, err := h.authUsecase.SSOStart()
	if err != nil {
		response.Error(c, http.StatusServiceUnavailable, err.Error())
		return
	}
	c.Redirect(http.StatusFound, url)
}

// SSOCallback menerima code dari Keycloak, menukar dengan user aplikasi,
// lalu redirect ke frontend dengan JWT internal di query string.
func (h *AuthHandler) SSOCallback(c *gin.Context) {
	fail := func(msg string) {
		c.Redirect(http.StatusFound, h.frontendURL+"/login?sso_error="+url.QueryEscape(msg))
	}
	user, err := h.authUsecase.SSOCallback(c.Request.Context(), c.Query("code"), c.Query("state"))
	if err != nil {
		fail(err.Error())
		return
	}
	token, err := auth.GenerateToken(user, h.jwtCfg.Secret, h.jwtCfg.Expiry)
	if err != nil {
		fail("gagal membuat sesi")
		return
	}
	q := url.Values{}
	q.Set("token", token)
	q.Set("id", strconv.FormatInt(user.ID, 10))
	q.Set("full_name", user.FullName)
	q.Set("email", user.Email)
	q.Set("role", user.Role)
	q.Set("photo_url", user.PhotoURL)
	c.Redirect(http.StatusFound, h.frontendURL+"/sso/callback?"+q.Encode())
}
