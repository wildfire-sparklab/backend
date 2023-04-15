package storage

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"wildfire-backend/internal/hotspots"
	"wildfire-backend/pkg/mysql"
)

type storage struct {
	client mysql.Client
}

func NewHotspotsStorage(client mysql.Client) hotspots.Storage {
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
	loc, _ := time.LoadLocation("Asia/Yakutsk")
	dateStart := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		12,
		0,
		0,
		0,
		loc)
	dateEnd := time.Date(
		date.Year(),
		date.Month(),
		date.Day()+1,
		12,
		0,
		0,
		0,
		loc)
	err := s.client.DB.Model(&hotspots.Hotspot{}).
		Where("time BETWEEN ? AND ?", dateStart, dateEnd).
		Find(&hostspots).Error
	if err != nil {
		return nil, err
	}
	return hostspots, nil
}

func (s storage) GetHotSpotsBySite(date time.Time) ([]hotspots.Hotspot, error) {
	var hostspots []hotspots.Hotspot
	dateStart := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0,
		0,
		0,
		0,
		date.Location())
	dateEnd := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		24,
		0,
		0,
		0,
		date.Location())
	err := s.client.DB.Model(&hotspots.Hotspot{}).
		Where("time BETWEEN ? AND ?", dateStart, dateEnd).
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
