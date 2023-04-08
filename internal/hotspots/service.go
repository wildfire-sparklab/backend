package hotspots

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
	"wildfire-backend/internal/config"
)

type Service struct {
	cfg     config.Config
	storage Storage
}

func NewHotSpotsService(cfg config.Config, storage Storage) *Service {
	return &Service{
		cfg:     cfg,
		storage: storage,
	}
}

func (s Service) AddsHotsSpots(hotspots []Hotspot) {
	var clearHotspots []Hotspot
	for _, hotspot := range hotspots {
		//lan := math.Floor(hotspot.Lan*10) / 10
		//long := math.Floor(hotspot.Long*10) / 10
		ok, err := s.storage.GetHotSpot(hotspot.Long, hotspot.Lan, hotspot.Time)
		fmt.Println(ok, err)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				clearHotspots = append(clearHotspots, hotspot)
			}
			fmt.Println("error", err)
			continue
		}

		//if !ok {
		//	clearHotspots = append(clearHotspots, hotspot)
		//	continue
		//} else if ok {
		//	hots, err := s.storage.GetHotSpot(lan, long, hotspot.Time)
		//	if hots != nil {
		//		continue
		//	}
		//
		//
		//	//count, err := s.storage.CountHotSpot(lan, long)
		//	//if err != nil {
		//	//	return
		//	//}
		//	//if count == 4 {
		//	//	err := s.storage.AddIgnoreHotSpot(hotspot)
		//	//	if err != nil {
		//	//		return
		//	//	}
		//	//} else {
		//	//	clearHotspots = append(clearHotspots, hotspot)
		//	//}
		//}
	}
	fmt.Println(clearHotspots)
	for _, hotspot := range clearHotspots {
		err := s.storage.AddHotSpot(hotspot)
		if err != nil {
			panic(err)
		}
	}
}

func (s Service) GetsHotSpotsByTime(t time.Time) []Hotspot {
	hotstpots, _ := s.storage.GetHotSpots(t)
	return hotstpots
}

func (s Service) GetsHotSpots() []Hotspot {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//t := time.Now()
	//req, err := http.NewRequestWithContext(ctx, "GET", "https://firms.modaps.eosdis.nasa.gov/api/area/csv/"+s.cfg.MapKey+"/VIIRS_SNPP_NRT/110,58,138,65/1/"+t.Format("2006-01-02"), nil)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://firms.modaps.eosdis.nasa.gov/api/area/csv/"+s.cfg.MapKey+"/VIIRS_SNPP_NRT/110,58,138,65/1/2023-03-25", nil)
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
		t, err := time.Parse("2006-01-02 15:04", h.AcqData+" "+h.AcqTime.String())
		if err != nil {
			fmt.Println(err)
			continue
		}
		//fmt.Println(h.AcqData+" "+h.AcqTime.String(), t)
		hotspots1 = append(hotspots1, Hotspot{
			Long:     h.Longitude,
			Lan:      h.Latitude,
			Time:     t,
			DayNight: h.DayNight,
		})
	}
	return hotspots1
}
