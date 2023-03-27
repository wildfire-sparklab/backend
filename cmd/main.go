package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"wildfire-backend/internal/checker"
	"wildfire-backend/internal/config"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/internal/hotspots/storage"
	"wildfire-backend/internal/wind"
	"wildfire-backend/pkg/postgres"
	"wildfire-backend/pkg/rabbit"
)

func main() {
	cfg := config.GetConfig()
	conn, err := rabbit.GetConn(cfg.AMQP)
	if err != nil {
		log.Panic("Not connection rabbit")
	}
	db := postgres.NewClient(context.TODO(), 5, cfg.Postgres)
	db.DB.AutoMigrate(&hotspots.Hotspot{})
	db.DB.AutoMigrate(&hotspots.IgnoreHotspot{})
	hotspotStorage := storage.NewHotspotsStorage(*db)
	hostpotService := hotspots.NewHotSpotsService(*cfg, hotspotStorage)
	windService := wind.NewWindService(*cfg)
	checkerService := checker.NewChecker(*hostpotService, *windService, conn)
	checkerService.Checker()
	checkerService.StartCheck()

	r := gin.Default()
	hotspotHandler := hotspots.NewHotSpotHandler(hotspotStorage)
	hotspotHandler.Register(r)
	r.Run(":8081")
}
