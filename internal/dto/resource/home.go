package resource

import (
	"time"
)

type HomeResource struct {
	App  string    `json:"app"`
	Env  string    `json:"env"`
	Time time.Time `json:"time"`
}
