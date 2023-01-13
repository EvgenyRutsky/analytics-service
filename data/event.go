package data

import "time"

type Event struct {
	Timestamp time.Time
	Value     int `json:"value"`
}
