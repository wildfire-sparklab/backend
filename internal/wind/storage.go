package wind

import "time"

type Storage interface {
	AddWind(wind Model) (Model, error)
	GetWinds(date time.Time) ([]Model, error)
	AddBroadcast(broadcast BroadCast) error
}
