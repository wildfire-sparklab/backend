package checker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/internal/wind"
	"wildfire-backend/pkg/rabbit"
)

type service struct {
	h  hotspots.Service
	w  wind.Service
	mq rabbit.Conn
}

func NewChecker(h hotspots.Service, w wind.Service, mq rabbit.Conn) *service {
	return &service{
		h:  h,
		w:  w,
		mq: mq,
	}
}

func (s service) StartCheck() {
	ticker := time.NewTicker(time.Minute * 30)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				s.Checker()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
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

type checkerSend struct {
	Hotspots []hotspots.Hotspot `json:"hotspots"`
	Winds    []wind.WeatherData `json:"winds"`
}

func (s service) Checker() {
	hotspotss := s.h.GetsHotSpots()
	s.h.AddsHotsSpots(hotspotss)
	var winds []wind.WeatherData
	for _, hotspot := range wind_cords {
		winds = append(winds, s.w.GetWind(hotspot[1], hotspot[0]))
	}
	sendMessage := checkerSend{Hotspots: hotspotss, Winds: winds}
	jsonMessage, err := json.Marshal(sendMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	sEnc := base64.StdEncoding.EncodeToString(jsonMessage)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	if err != nil {
		fmt.Println("Json not convert")
	}
	fmt.Println(sEnc)
	err = s.mq.Publish("broker", []byte(sEnc))
	if err != nil {
		fmt.Println("publish error")
	}
}
