package controller

import (
	"net/http"

	"github.com/deshortone/ledger-system/internal/platform/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	platformService domain.PlatformHealth
}

func NewHandler(platformService domain.PlatformHealth) Handler {
	return Handler{
		platformService: platformService,
	}
}

func (h *Handler) RegisterRoutes(c *gin.RouterGroup) {
	c.GET("/liveness", h.getLiveness)
	c.GET("/readiness", h.getReadiness)
}

// getLiveness godoc
//
//	@Summary		Check if app is live
//	@Description	Check if app is live
//	@Tags			health
//	@Success		200	{object}	string
//	@Failure		503	{object}	string
//	@Router			/health/liveness		[get]
func (h *Handler) getLiveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

// getReadiness godoc
//
//	@Summary		Check if app is ready to handle connections
//	@Description	Check if app is ready to handle connections
//	@Tags			health
//	@Success		200	{object}	string
//	@Failure		503	{object}	string
//	@Router			/health/readiness		[get]
func (h *Handler) getReadiness(c *gin.Context) {
	if err := h.platformService.IsUp(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
