package storage

import (
	"time"

	"github.com/blankInPajamas/SentinelGo/internal/logs"
)

type QueryFilter struct {
	StartTime time.Time
	EndTime   time.Time
	SourceIP  string
	User      string
	Outcome   string
	EventType string
	Limit     int
	Offset    int
}

type Storage interface {
	Save(event logs.Event) error
	Query(filter QueryFilter) ([]logs.Event, error)
	Count() int
	Close() error
}