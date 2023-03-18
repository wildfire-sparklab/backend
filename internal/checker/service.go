package checker

import (
	"fmt"
	"time"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/internal/wind"
)

type service struct {
	h hotspots.Service
	w wind.Service
}

func NewChecker(h hotspots.Service, w wind.Service) *service {
	return &service{
		h: h,
		w: w,
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

func (s service) Checker() {
	hotspotss := s.h.GetsHotSpots()
	s.h.AddsHotsSpots(hotspotss)
	var winds []wind.Forecast5WeatherData
	for _, hotspot := range hotspotss {
		winds = append(winds, s.w.GetWind(hotspot.Long, hotspot.Lan))
	}
	//fmt.Println(winds)
	fmt.Println(hotspotss)
}
