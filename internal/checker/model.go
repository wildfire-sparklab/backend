package checker

type SendDataDTO struct {
	Winds    []int64 `json:"winds"`
	Hotspots []int64 `json:"hotspots"`
	Date     string  `json:"date"`
}
