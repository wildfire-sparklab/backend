package hotspots

import (
	"time"
)

type HotSpotCSV struct {
	Latitude   float64 `csv:"latitude"`
	Longitude  float64 `csv:"longitude"`
	Bright     string  `csv:"bright_ti4"`
	Scan       string  `csv:"scan"`
	Track      string  `csv:"track"`
	AcqData    string  `csv:"acq_date"`
	AcqTime    AcqTime `csv:"acq_time"`
	Satellite  string  `csv:"satellite"`
	Instrument string  `csv:"instrument"`
	Confidence string  `csv:"confidence"`
	Bright5    string  `csv:"bright_ti5"`
	frp        string  `csv:"frp"`
	DayNight   string  `csv:"daynight"`
}

type AcqTime struct {
	string
}

func (acqTime *AcqTime) UnmarshalCSV(csv string) (err error) {
	chars := []rune(csv)
	m := chars[2:]
	h := chars[:2]
	acqTime.string = string(h) + ":" + string(m)
	return nil
}

func (acqTime *AcqTime) String() string {
	return acqTime.string
}

type Hotspot struct {
	Id       int64 `gorm:"primaryKey"`
	Time     time.Time
	Long     float64
	Lan      float64
	DayNight string
}

type IgnoreHotspot struct {
	Id   int `gorm:"primaryKey"`
	Long float64
	Lan  float64
}
