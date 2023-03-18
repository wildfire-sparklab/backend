package storage

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/pkg/postgres"
)

type storage struct {
	client postgres.Client
}

func NewHotspotsStorage(client postgres.Client) hotspots.Storage {
	return storage{
		client: client,
	}
}

func (s storage) AddHotSpot(hotspot hotspots.Hotspot) error {
	err := s.client.DB.Create(&hotspot).Error
	return err
}

func (s storage) AddIgnoreHotSpot(hotspot hotspots.IgnoreHotspot) error {
	err := s.client.DB.Create(&hotspot).Error
	return err
}

func (s storage) GetHotSpots(date time.Time) ([]hotspots.Hotspot, error) {
	var hostspots []hotspots.Hotspot
	dateStart := time.Date(
		date.YearDay(),
		date.Month(),
		date.Day(),
		0,
		0,
		0,
		0,
		date.Location())
	err := s.client.DB.Model(&hotspots.Hotspot{}).
		Where("time BETWEEN ? AND ?", dateStart, date).
		Find(&hostspots).Error
	if err != nil {
		return nil, err
	}
	return hostspots, nil
}

func (s storage) GetHotSpot(Long float64, Lan float64, date time.Time) (*hotspots.Hotspot, error) {
	var hotspot hotspots.Hotspot
	err := s.client.DB.
		Model(&hotspots.Hotspot{}).
		Where(
			&hotspots.Hotspot{
				Long: Long,
				Lan:  Lan,
				Time: date}).
		Last(&hotspot).
		Error
	return &hotspot, err
}

func (s storage) CheckHotSpot(Long float64, Lan float64) (bool, error) {
	var hotspot hotspots.Hotspot
	err := s.client.DB.Model(&hotspots.Hotspot{}).
		Where(&hotspots.Hotspot{Lan: Lan, Long: Long}).Find(&hotspot).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s storage) CountHotSpot(Long float64, Lan float64) (int, error) {
	var count *int64
	err := s.client.DB.Model(&hotspots.Hotspot{}).
		Where(&hotspots.Hotspot{Lan: Lan, Long: Long}).Count(count).Error
	return int(*count), err
}
