package checker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/internal/wind"
	"wildfire-backend/pkg/rabbit"
)

type service struct {
	h  hotspots.Service
	w  wind.Service
	mq *rabbit.Conn
}

func NewChecker(h hotspots.Service, w wind.Service, mq *rabbit.Conn) *service {
	return &service{
		h:  h,
		w:  w,
		mq: mq,
	}
}

func (s service) StartCheck() {
	loc, _ := time.LoadLocation("Asia/Yakutsk")
	cron := gocron.NewScheduler(loc)
	cron.Every(1).Day().At("12:00").Do(func() {
		s.StartAutomata(time.Now())
	})
	cron.Every("30m").Do(func() {
		s.Checker()
	})
	cron.StartAsync()
}

// пока что тут
// 1-lat 2-lon
var wind_cords = [][]float64{
	{58.0, 110.0},
	{59.75, 110},
	{61.5, 110},
	{63.25, 110},
	{65, 110},
	{58, 114},
	{59.75, 114},
	{61.5, 114},
	{63.25, 114},
	{65, 114},
	{58, 118},
	{59.75, 118},
	{61.5, 118},
	{63.25, 118},
	{65, 118},
	{58, 122},
	{59.75, 122},
	{61.5, 122},
	{63.25, 122},
	{65, 122},
	{58, 126},
	{59.75, 126},
	{61.5, 126},
	{63.25, 126},
	{65, 126},
	{58, 130},
	{59.75, 130},
	{61.5, 130},
	{63.25, 130},
	{65, 130},
	{58, 134},
	{59.75, 134},
	{61.5, 134},
	{63.25, 134},
	{65, 134},
	{58, 138},
	{59.75, 138},
	{61.5, 138},
	{63.25, 138},
	{65, 138},
}

//type checkerSend struct {
//	Hotspots []hotspots.HotspotJson `json:"hotspots"`
//	Winds    []wind.WeatherData     `json:"winds"`
//}

func (s service) Checker() {
	fmt.Println("Check hotspots...")
	hotspotss := s.h.GetsHotSpots()
	s.h.AddsHotsSpots(hotspotss)
}

func (s service) StartAutomata(date time.Time) {
	fmt.Println("Start automata...")
	t := date.AddDate(0, 0, -1)
	hotspotss := s.h.GetsHotSpotsByTime(t)
	if len(hotspotss) == 0 {
		fmt.Println("array hotspot zero")
		return
	}
	windss, err := s.w.GetWinds(date)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(windss) == 0 && date.Year() > 2022 {
		var winds []wind.WeatherData
		for _, w := range wind_cords {
			winds = append(winds, s.w.GetWind(w[1], w[0]))
		}
		s.w.AddWind(winds)
		fmt.Println("array winds zero")
	}
	windss, err = s.w.GetWinds(date)
	if err != nil {
		fmt.Println(err)
		return
	}
	var hotspotss1 []int64
	for _, hotspot := range hotspotss {
		hotspotss1 = append(hotspotss1, hotspot.Id)
	}

	var winds1 []int64
	for _, wind1 := range windss {
		winds1 = append(winds1, wind1.Id)
	}
	fmt.Println(windss)
	sendData := &SendDataDTO{
		Hotspots: hotspotss1,
		Winds:    winds1,
		Date:     date.Format("2006-01-02"),
	}
	jsonMessage, err := json.Marshal(sendData)
	if err != nil {
		fmt.Println(err)
		return
	}
	sEnc := base64.StdEncoding.EncodeToString(jsonMessage)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
		return
	}
	if err != nil {
		fmt.Println("Json not convert")
		return
	}
	err = s.mq.Publish("broker", []byte(sEnc))
	if err != nil {
		fmt.Println("publish error")
		return
	}
}
