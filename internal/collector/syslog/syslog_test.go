package syslog_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/blankInPajamas/SentinelGo/internal/collector/syslog"
)

func TestSysCollector_StartAndClose(t *testing.T) {
	collector := syslog.New(":0") // Port 0 = let OS pick random available port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var received []string
	var mu sync.Mutex

	handler := func(line string) {
		mu.Lock()
		received = append(received, line)
		mu.Unlock()
	}

	// Start in background
	done := make(chan error, 1)
	go func() {
		done <- collector.Start(ctx, handler)
	}()

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	// Cancel context to stop
	cancel()
	err := <-done

	// Should exit with context cancelled error, not a panic
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
