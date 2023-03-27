// Package wind https://github.com/briandowns/openweathermap
package wind

import (
	"encoding/json"
	"io"
	"strings"
	"time"
)

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type DtTxt struct {
	time.Time
}

func (dt *DtTxt) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("2006-01-02 15:04:05", strings.Trim(string(b), "\""))
	dt.Time = t
	return err
}

func (dt *DtTxt) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt)
}

// Forecast5WeatherList holds specific query data
type Forecast5WeatherList struct {
	Dt      int       `json:"dt"`
	Main    Main      `json:"main"`
	Weather []Weather `json:"weather"`
	Clouds  Clouds    `json:"clouds"`
	Wind    Wind      `json:"wind"`
	Rain    Rain      `json:"rain"`
	Snow    Snow      `json:"snow"`
	DtTxt   DtTxt     `json:"dt_txt"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	FeelsLike float64 `json:"feels_like"`
	Pressure  float64 `json:"pressure"`
	SeaLevel  float64 `json:"sea_level"`
	GrndLevel float64 `json:"grnd_level"`
	Humidity  int     `json:"humidity"`
}

// Forecast5WeatherData will hold returned data from queries
type Forecast5WeatherData struct {
	// COD     string                `json:"cod"`
	// Message float64               `json:"message"`
	Cnt  int                    `json:"cnt"`
	List []Forecast5WeatherList `json:"list"`
}

type ForecastWeatherList struct {
	Dt   int  `json:"dt"`
	Wind Wind `json:"wind"`
}

type WeatherData struct {
	List []ForecastWeatherList `json:"list"`
}

func (f *Forecast5WeatherData) Decode(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(&f); err != nil {
		return err
	}
	return nil
}

type Clouds struct {
	All int `json:"all"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

type Rain struct {
	OneH   float64 `json:"1h,omitempty"`
	ThreeH float64 `json:"3h,omitempty"`
}

type Snow struct {
	OneH   float64 `json:"1h,omitempty"`
	ThreeH float64 `json:"3h,omitempty"`
}
