package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"wildfire-backend/internal/checker"
	"wildfire-backend/internal/config"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/internal/hotspots/storage"
	"wildfire-backend/internal/wind"
	"wildfire-backend/pkg/mysql"
)

func main() {
	cfg := config.GetConfig()
	//conn, err := rabbit.GetConn(cfg.AMQP)
	//if err != nil {
	//	log.Panic("Not connection rabbit")
	//}
	db := mysql.NewClient(context.TODO(), 5, cfg.MySQL)
	db.DB.AutoMigrate(&hotspots.Hotspot{})
	db.DB.AutoMigrate(&hotspots.IgnoreHotspot{})
	hotspotStorage := storage.NewHotspotsStorage(*db)
	hostpotService := hotspots.NewHotSpotsService(*cfg, hotspotStorage)
	windService := wind.NewWindService(*cfg)
	checkerService := checker.NewChecker(*hostpotService, *windService, nil)
	checkerService.StartCheck()

	r := gin.Default()
	hotspotHandler := hotspots.NewHotSpotHandler(hotspotStorage)
	hotspotHandler.Register(r)
	r.Run(":8081")
}
