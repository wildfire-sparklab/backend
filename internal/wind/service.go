package wind

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"wildfire-backend/internal/config"
)

type Service struct {
	cfg     config.Config
	storage Storage
}

func NewWindService(cfg config.Config, storage Storage) *Service {
	return &Service{
		cfg:     cfg,
		storage: storage,
	}
}

func (s Service) AddWind(winds []WeatherData, date time.Time) error {
	for _, w := range winds {
		model := Model{
			Lan:  w.Lan,
			Long: w.Long,
			Time: date,
		}
		model, err := s.storage.AddWind(model)
		if err != nil {
			return err
		}
		fmt.Println(model)
		for _, list := range w.List {
			broadcast := BroadCast{
				Dt:      list.Dt,
				Speed:   list.Wind.Speed,
				Deg:     list.Wind.Deg,
				ModelID: model.Id,
			}
			err := s.storage.AddBroadcast(broadcast)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s Service) GetWinds(time2 time.Time) (Model, error) {
	return s.storage.GetWinds(time2)
}

func (s Service) GetWind(Long float64, Lan float64) WeatherData {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&appid=%s", Lan, Long, s.cfg.WindKey), nil)
	if err != nil {
		fmt.Println(err)
	}
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var data Forecast5WeatherData
	err = json.Unmarshal(b, &data)
	if err != nil {
		panic(err)
	}
	var data1 WeatherData
	for i, list := range data.List {
		if i > 24 {
			break
		}
		data1.List = append(data1.List, ForecastWeatherList{
			Dt:   list.Dt,
			Wind: list.Wind,
		})
	}
	data1.Lan = Lan
	data1.Long = Long
	return data1
}
