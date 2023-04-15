package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"wildfire-backend/internal/config"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/internal/wind"
	storage1 "wildfire-backend/internal/wind/storage"
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
	db.DB.AutoMigrate(&wind.Model{})
	db.DB.AutoMigrate(&wind.BroadCast{})
	//s3Client := s32.NewClient(cfg.S3)

	//hotspotStorage := storage.NewHotspotsStorage(*db)
	//hostpotService := hotspots.NewHotSpotsService(*cfg, hotspotStorage)

	windStorage := storage1.NewWindStorage(*db)
	service := wind.NewWindService(*cfg, windStorage)
	files, _ := os.ReadDir("./output")
	for _, file := range files {
		jsonFile, _ := os.Open("./output/" + file.Name())
		byteValue, _ := ioutil.ReadAll(jsonFile)
		name := strings.Split(file.Name(), "_")
		name1, _ := strconv.Atoi(strings.Split(name[1], ".json")[0])
		if name1 == 12 {
			fmt.Println(name[0])
			var result wind.TestWeatherData
			json.Unmarshal([]byte(byteValue), &result)
			date, _ := time.Parse("2006-01-02", name[0])
			service.AddWind(result.Data, date)
		}

		//fmt.Println(result)

		//service.AddWind(result.Data, time.Parse())
	}
	//
	//checkerService := checker.NewChecker(*hostpotService, *windService, nil)
	//checkerService.StartCheck()
	//
	//r := gin.Default()
	//hotspotHandler := hotspots.NewHotSpotHandler(hotspotStorage)
	//s3handler := s3.NewS3Handler(s3Client)
	//s3handler.Register(r)
	//hotspotHandler.Register(r)
	//r.Run(":8081")
}
