package alert

import "time"

type Alert struct {
	ID         string
	Timestamp  time.Time
	Rule       string
	Severity   string // low, medium, high, critical
	SourceIP   string
	User       string
	Message    string
	EventCount int
}
