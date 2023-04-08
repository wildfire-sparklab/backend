package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"wildfire-backend/internal/checker"
	"wildfire-backend/internal/config"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/internal/hotspots/storage"
	"wildfire-backend/internal/s3"
	"wildfire-backend/internal/wind"
	"wildfire-backend/pkg/mysql"
	s32 "wildfire-backend/pkg/s3"
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
	s3Client := s32.NewClient(cfg.S3)

	hotspotStorage := storage.NewHotspotsStorage(*db)
	hostpotService := hotspots.NewHotSpotsService(*cfg, hotspotStorage)
	windService := wind.NewWindService(*cfg)
	checkerService := checker.NewChecker(*hostpotService, *windService, nil)
	checkerService.StartCheck()

	r := gin.Default()
	hotspotHandler := hotspots.NewHotSpotHandler(hotspotStorage)
	s3handler := s3.NewS3Handler(s3Client)
	s3handler.Register(r)
	hotspotHandler.Register(r)
	r.Run(":8081")
}
