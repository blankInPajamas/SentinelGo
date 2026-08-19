package bruteforce

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blankInPajamas/SentinelGo/internal/alert"
	"github.com/blankInPajamas/SentinelGo/internal/logs"
	"github.com/blankInPajamas/SentinelGo/internal/storage"
)

type BruteForceDetector struct {
	store         storage.Storage
	threshold     int
	timeWindow    time.Duration
	checkInterval time.Duration
	mu            sync.Mutex
	lastCheck     time.Time
}

func New(store storage.Storage, threshold int, timeWindow time.Duration, checkInterval time.Duration) *BruteForceDetector {
	return &BruteForceDetector{
		store:         store,
		threshold:     threshold,
		timeWindow:    timeWindow,
		checkInterval: checkInterval,
	}
}

func (d *BruteForceDetector) Start(
	ctx context.Context,
	alertChan chan<- alert.Alert,
) error {
	ticker := time.NewTicker(d.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			d.detect(ctx, alertChan)
		}
	}
}

func (d *BruteForceDetector) detect(
	ctx context.Context,
	alertChan chan<- alert.Alert,
) {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	startTime := now.Add(-d.timeWindow)

	events, err := d.store.Query(storage.QueryFilter{
		StartTime: startTime,
		EndTime:   now,
		Outcome:   "failure",
		EventType: "auth",
	})
	if err != nil {
		fmt.Printf("detector query error: %v\n", err)
		return
	}

	ipFailures := make(map[string][]logs.Event)

	for _, e := range events {
		ipFailures[e.SourceIP] = append(ipFailures[e.SourceIP], e)
	}

	for ip, failures := range ipFailures {
		if len(failures) < d.threshold {
			continue
		}

		user := d.mostTargetedUser(failures)

		alertChan <- alert.Alert{
			Timestamp: now,
			Rule:      "brute-force-ssh",
			Severity:  "high",
			SourceIP:  ip,
			User:      user,
			Message: fmt.Sprintf(
				"Detected %d failed logins from %s in %v",
				len(failures),
				ip,
				d.timeWindow,
			),
			EventCount: len(failures),
		}
	}

	d.lastCheck = now
}

func (d *BruteForceDetector) mostTargetedUser(
	failures []logs.Event,
) string {
	userCounts := make(map[string]int)

	for _, f := range failures {
		userCounts[f.User]++
	}

	maxUser := ""
	maxCount := 0

	for user, count := range userCounts {
		if count > maxCount {
			maxCount = count
			maxUser = user
		}
	}

	return maxUser
}

func (d *BruteForceDetector) Close() error {
	return nil
}
