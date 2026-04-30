package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, err error) {
	if err == nil {
		c.JSON(status, gin.H{"error": http.StatusText(status)})
		return
	}

	// Keep error responses predictable for learning/debugging.
	// In production, you may want to hide internal error details.
	c.JSON(status, gin.H{
		"error": err.Error(),
	})
}

func ErrorMessage(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}
