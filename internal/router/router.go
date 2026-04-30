package router

import (
	"net/http"

	"ecommerce-go/internal/config"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	// Logger + Recovery are the common "production learning" defaults.
	r.Use(gin.Logger(), gin.Recovery())

	// Simple liveness endpoint.
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Consistent JSON for unknown paths (Gin's default 404 is HTML).
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not found",
			"path":    c.Request.URL.Path,
			"method":  c.Request.Method,
		})
	})

	return r
}
