package storage

import (
	"time"
)

// Storage interface is an object which can perform DB functions
type Storage interface {
	SendAsync(time time.Time, value int) error
	GetValuesForLastMin(t time.Time) ([]int, error)
}
