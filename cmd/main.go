package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"wildfire-backend/internal/checker"
	"wildfire-backend/internal/config"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/internal/hotspots/storage"
	"wildfire-backend/internal/wind"
	"wildfire-backend/pkg/postgres"
)

func main() {
	cfg := config.GetConfig()
	db := postgres.NewClient(context.TODO(), 5, cfg.Postgres)
	db.DB.AutoMigrate(&hotspots.Hotspot{})
	db.DB.AutoMigrate(&hotspots.IgnoreHotspot{})
	hotspotStorage := storage.NewHotspotsStorage(*db)
	hostpotService := hotspots.NewHotSpotsService(*cfg, hotspotStorage)
	windService := wind.NewWindService(*cfg)
	checkerService := checker.NewChecker(*hostpotService, *windService)
	checkerService.Checker()

	r := gin.Default()
	hotspotHandler := hotspots.NewHotSpotHandler(hotspotStorage)
	hotspotHandler.Register(r)
	r.Run(":8081")
}
