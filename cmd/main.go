package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"wildfire-backend/internal/checker"
	"wildfire-backend/internal/config"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/internal/hotspots/storage"
	"wildfire-backend/internal/s3"
	"wildfire-backend/internal/wind"
	storage1 "wildfire-backend/internal/wind/storage"
	"wildfire-backend/pkg/mysql"
	"wildfire-backend/pkg/rabbit"
	s32 "wildfire-backend/pkg/s3"
)

func main() {
	cfg := config.GetConfig()
	conn, err := rabbit.GetConn(cfg.AMQP)
	if err != nil {
		log.Panic("Not connection rabbit")
	}
	db := mysql.NewClient(context.TODO(), 5, cfg.MySQL)
	db.DB.AutoMigrate(&hotspots.Hotspot{})
	db.DB.AutoMigrate(&hotspots.IgnoreHotspot{})
	db.DB.AutoMigrate(&wind.Model{})
	db.DB.AutoMigrate(&wind.BroadCast{})
	s3Client := s32.NewClient(cfg.S3)

	hotspotStorage := storage.NewHotspotsStorage(*db)
	hostpotService := hotspots.NewHotSpotsService(*cfg, hotspotStorage)

	windStorage := storage1.NewWindStorage(*db)
	windService := wind.NewWindService(*cfg, windStorage)
	//files, _ := os.ReadDir("./output")
	//for _, file := range files {
	//	jsonFile, _ := os.Open("./output/" + file.Name())
	//	byteValue, _ := ioutil.ReadAll(jsonFile)
	//	name := strings.Split(file.Name(), "_")
	//	name1, _ := strconv.Atoi(strings.Split(name[1], ".json")[0])
	//	if name1 == 12 {
	//		fmt.Println(name[0])
	//		var result wind.TestWeatherData
	//		json.Unmarshal([]byte(byteValue), &result)
	//		date, _ := time.Parse("2006-01-02", name[0])
	//		service.AddWind(result.Data, date)
	//	}
	//
	//	//fmt.Println(result)
	//
	//	//service.AddWind(result.Data, time.Parse())
	//}
	//
	checkerService := checker.NewChecker(*hostpotService, *windService, &conn)
	checkerService.StartCheck()

	r := gin.Default()
	r.Use(CORSMiddleware())
	gin.SetMode(gin.ReleaseMode)
	hotspotHandler := hotspots.NewHotSpotHandler(hotspotStorage)

	checkerHandler := checker.NewCheckerHandler(*checkerService)
	checkerHandler.Register(r)
	s3handler := s3.NewS3Handler(s3Client)
	s3handler.Register(r)
	hotspotHandler.Register(r)
	r.Run(":8081")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
