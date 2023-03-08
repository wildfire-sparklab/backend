package hotspots

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"io"
	"math"
	"net/http"
	"time"
	"wildfire-backend/internal/config"
)

type Service struct {
	cfg     config.Config
	storage Storage
}

func NewHotSpotsService() *Service {
	return &Service{}
}

func (s Service) AddsHotsSpots(hotspots []Hotspot) {
	var clearHotspots []Hotspot
	for _, hotspot := range hotspots {
		lan := math.Floor(hotspot.Lan*10) / 10
		long := math.Floor(hotspot.Long*10) / 10
		ok, err := s.storage.CheckHotSpot(lan, long)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				clearHotspots = append(clearHotspots, hotspot)
				continue
			}
			fmt.Println("error", err)
			continue
		}
		if ok {
			hots, _ := s.storage.GetHotSpot(lan, long, hotspot.Time)
			if hots != nil {
				continue
			}
			count, err := s.storage.CountHotSpot(lan, long)
			if err != nil {
				return
			}
			if count == 4 {
				err := s.storage.AddIgnoreHotSpot(hotspot)
				if err != nil {
					return
				}
			} else {
				clearHotspots = append(clearHotspots, hotspot)
			}
		}
	}
	for _, hotspot := range clearHotspots {
		err := s.storage.AddHotSpot(hotspot)
		if err != nil {
			return
		}
	}
}

func (s Service) GetsHotSpots() []Hotspot {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	t := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://firms.modaps.eosdis.nasa.gov/api/area/csv/"+s.cfg.MapKey+"/VIIRS_SNPP_NRT/90,40,175,80/1/"+t.Format("2006-01-02"), nil)
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
	gocsv.TagSeparator = ","
	var hotspots []HotSpotCSV
	err = gocsv.UnmarshalBytes(b, &hotspots)
	if err != nil {
		fmt.Println(err)
	}
	var hotspots1 []Hotspot
	for _, h := range hotspots {
		t, _ := time.Parse("2006-02-01 15:04", h.AcqData+" "+h.AcqTime.String())
		hotspots1 = append(hotspots1, Hotspot{
			Long:     h.Longitude,
			Lan:      h.Latitude,
			Time:     t,
			DayNight: h.DayNight,
		})
	}
	return hotspots1
}
