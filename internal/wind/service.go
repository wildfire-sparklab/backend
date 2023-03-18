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
	cfg config.Config
}

func NewWindService(cfg config.Config) *Service {
	return &Service{
		cfg: cfg,
	}
}

func (s Service) GetWind(Long float64, Lan float64) Forecast5WeatherData {
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
	return data
}
