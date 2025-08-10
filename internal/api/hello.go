package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloHandler responds with a simple greeting.
// @Summary      Returns a greeting
// @Description  Simple hello world endpoint
// @Tags         example
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/hello [get]
func HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello world"})
}
