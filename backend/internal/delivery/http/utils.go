package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseID mengambil ID dari URL parameter ":id".
func parseID(c *gin.Context) (int64, error) {
	return parseIDParam(c, "id")
}

// parseIDParam mengambil ID dari URL parameter dengan nama tertentu.
func parseIDParam(c *gin.Context, name string) (int64, error) {
	return strconv.ParseInt(c.Param(name), 10, 64)
}
