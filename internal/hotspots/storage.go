package hotspots

import "time"

type Storage interface {
	AddHotSpot(hotspot Hotspot) error
	AddIgnoreHotSpot(hotspot IgnoreHotspot) error
	GetHotSpots(date time.Time) ([]Hotspot, error)
	GetHotSpot(Long float64, Lan float64, date time.Time) (*Hotspot, error)
	CheckHotSpot(Long float64, Lan float64) (bool, error)
	CountHotSpot(Long float64, Lan float64) (int, error)
}
