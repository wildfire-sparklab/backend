package checker

import (
	"time"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/internal/wind"
)

type service struct {
	h hotspots.Service
	w wind.Service
}

func NewChecker() *service {
	return &service{}
}

func (s service) StartCheck() {
	ticker := time.NewTicker(time.Hour * 3)
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

}
