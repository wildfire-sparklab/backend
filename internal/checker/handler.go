package checker

import (
	"github.com/gin-gonic/gin"
	"time"
	"wildfire-backend/internal/handlers"
)

type handler struct {
	s service
}

func NewCheckerHandler(s service) handlers.Handler {
	return &handler{
		s: s,
	}
}

func (h handler) Register(router *gin.Engine) {
	router.POST("/checker/:date", h.CheckerHandler)
}

func (h *handler) CheckerHandler(ctx *gin.Context) {
	t, err := time.Parse("2006-01-02", ctx.Param("date"))
	if err != nil {
		ctx.JSON(500, "Error parse time")
		return
	}
	h.s.StartAutomata(t)
	ctx.Status(200)
}
