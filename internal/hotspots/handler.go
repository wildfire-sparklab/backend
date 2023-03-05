package hotspots

import (
	"github.com/gin-gonic/gin"
	"wildfire-backend/internal/handlers"
)

type handler struct {
}

func NewHotSpotHandler() handlers.Handler {
	return &handler{}
}

func (h handler) Register(router *gin.Engine) {
	router.GET("/hotspot/:time", h.GetHotSpots)
}

func (h *handler) GetHotSpots(ctx *gin.Context) {

}
