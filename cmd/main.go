package main

import (
	"context"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"net/http"
	"time"
	"wildfire-backend/internal/config"
	hotspots2 "wildfire-backend/internal/hotspots"
)

func main() {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	t := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://firms.modaps.eosdis.nasa.gov/api/area/csv/"+cfg.MapKey+"/VIIRS_SNPP_NRT/90,40,175,80/1/"+t.Format("2006-01-02"), nil)
	if err != nil {
		panic(err)
	}
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	gocsv.TagSeparator = ","
	var hotspots []hotspots2.HotSpotCSV
	err = gocsv.UnmarshalBytes(b, &hotspots)
	if err != nil {
		panic(err)
	}
	var hotspots1 []hotspots2.Hotspot
	fmt.Println("hotspots:")
	for _, h := range hotspots {
		t, _ := time.Parse("2006-02-01 15:04", h.AcqData+" "+h.AcqTime.String())
		hotspots1 = append(hotspots1, hotspots2.Hotspot{
			Long: h.Longitude,
			Lan:  h.Latitude,
			Time: t,
		})
	}
	for _, h := range hotspots1 {
		fmt.Println(h.Lan, h.Long, h.Time)
	}

}
