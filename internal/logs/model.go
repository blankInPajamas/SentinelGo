package logs

import "time"

type Event struct {
	Timestamp 	time.Time
	SourceIP	string
	User		string
	EventType 	string
	Outcome 	string
	Host 		string
	RawMessage	string
}