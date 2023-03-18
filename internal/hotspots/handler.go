package hotspots

import (
	"fmt"
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

//type GetHotspotsDTO struct {
//	Date time.Time `query:"time" time_format:"2006-01-02" time_utc:"1"`
//}

func (h *handler) GetHotSpots(ctx *gin.Context) {
	//TODO parse time
	//var gethotspot GetHotspotsDTO
	//if ctx.ShouldBindQuery(&gethotspot) != nil {
	//	ctx.JSON(500, "Error parse time")
	//	return
	//}
	t, err := time.Parse("2006-01-02", ctx.Param("time"))
	if err != nil {
		ctx.JSON(500, "Error parse time")
		return
	}
	fmt.Println(t)
	hotspots, err := h.storage.GetHotSpots(t)
	if err != nil {
		ctx.JSON(500, "Server error")
		return
	}
	ctx.JSON(200, hotspots)
}
