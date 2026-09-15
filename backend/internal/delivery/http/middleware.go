package http

import (
	"net/http"
	"strings"

	"github.com/e-notary-bprs/backend/internal/config"
	"github.com/e-notary-bprs/backend/pkg/auth"
	"github.com/e-notary-bprs/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware adalah middleware untuk validasi JWT token (Gin).
type AuthMiddleware struct {
	cfg *config.JWTConfig
}

func NewAuthMiddleware(cfg *config.JWTConfig) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

// Handler memvalidasi Authorization: Bearer <token>, menyimpan user_id & user_role ke context.
func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Error(c, http.StatusUnauthorized, "invalid authorization format")
			c.Abort()
			return
		}
		claims, err := auth.ParseToken(tokenString, m.cfg.Secret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

// RequireRole memastikan user memiliki salah satu role yang diizinkan.
func (m *AuthMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		roleStr, _ := role.(string)
		for _, rr := range allowedRoles {
			if roleStr == rr {
				c.Next()
				return
			}
		}
		response.Error(c, http.StatusForbidden, "insufficient permissions")
		c.Abort()
	}
}
