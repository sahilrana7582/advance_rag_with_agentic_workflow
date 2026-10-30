package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sahilrana27582/sentinel-ai/sre-api/internal/config"
)

type Handler struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"env":    h.cfg.APIEnv,
		"port":   h.cfg.APIPort,
	})
}
