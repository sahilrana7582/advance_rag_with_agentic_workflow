package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sahilrana27582/sentinel-ai/sre-api/internal/handler"
)

func Setup(r *gin.Engine, h *handler.Handler) {
	r.GET("/health", h.Health)
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": "0.1.0"})
	})
}
