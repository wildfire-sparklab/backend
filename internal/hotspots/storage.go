package hotspots

import "time"

type Storage interface {
	AddHotSpot(hotspot Hotspot) error
	GetHotSpots(date time.Time) ([]Hotspot, error)
	CheckHotSpot(Long float64, Lan float64) (bool, error)
}
