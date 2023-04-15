package hotspots

import (
	"github.com/gin-gonic/gin"
	"time"
	"wildfire-backend/internal/handlers"
)

type handler struct {
	storage Storage
}

func NewHotSpotHandler(storage Storage) handlers.Handler {
	return &handler{
		storage: storage,
	}
}

func (h handler) Register(router *gin.Engine) {
	router.GET("/hotspot/:time", h.GetHotSpots)
}

func (h *handler) GetHotSpots(ctx *gin.Context) {
	t, err := time.Parse("2006-01-02", ctx.Param("time"))
	if err != nil {
		ctx.JSON(500, "Error parse time")
		return
	}
	hotspots, err := h.storage.GetHotSpotsBySite(t)
	if err != nil {
		ctx.JSON(500, "Server error")
		return
	}
	ctx.JSON(200, hotspots)
}
