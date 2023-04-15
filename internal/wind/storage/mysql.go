package storage

import (
	"fmt"
	"time"
	"wildfire-backend/internal/wind"
	"wildfire-backend/pkg/mysql"
)

type storage struct {
	client mysql.Client
}

func NewWindStorage(client mysql.Client) wind.Storage {
	return storage{
		client: client,
	}
}

func (s storage) AddWind(wind wind.Model) (wind.Model, error) {
	err := s.client.DB.Create(&wind).Error
	return wind, err
}

func (s storage) AddBroadcast(broadcast wind.BroadCast) error {
	err := s.client.DB.Create(&broadcast).Error
	return err
}

func (s storage) GetWinds(date time.Time) ([]wind.Model, error) {
	var winds []wind.Model
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
	fmt.Println(dateStart, dateEnd)
	err := s.client.DB.Model(&wind.Model{}).
		Where("time BETWEEN ? AND ?", dateStart, dateEnd).
		Find(&winds).Error
	if err != nil {
		return nil, err
	}
	return winds, nil
}
