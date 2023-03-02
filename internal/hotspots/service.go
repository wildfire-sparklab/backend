package hotspots

import (
	"context"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"net/http"
	"time"
	"wildfire-backend/internal/config"
)

type service struct {
	cfg config.Config
}

func NewHotSpotsService() *service {

	return &service{}
}

func (s service) AddsHotsSpots() {
	s.GetsHotSpots()
}

func (s service) GetsHotSpots() []Hotspot {
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
			Long: h.Longitude,
			Lan:  h.Latitude,
			Time: t,
		})
	}
	return hotspots1
}

func (s service) StartCheck() {
	ticker := time.NewTicker(time.Hour * 3)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				s.GetsHotSpots()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
