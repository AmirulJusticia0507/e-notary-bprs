package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Payload struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func JSON(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, Payload{Status: http.StatusText(statusCode), Data: data})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Payload{Status: http.StatusText(statusCode), Error: message})
}
