package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blankInPajamas/SentinelGo/internal/alert"
	"github.com/blankInPajamas/SentinelGo/internal/collector/syslog"
	"github.com/blankInPajamas/SentinelGo/internal/detector/bruteforce"
	"github.com/blankInPajamas/SentinelGo/internal/parser/auth"
	"github.com/blankInPajamas/SentinelGo/internal/storage/memory"
	"github.com/gin-gonic/gin"
)

func main() {
	store := memory.New()
	defer store.Close()

	collector := syslog.New(":1514")
	detector := bruteforce.New(store, 5, 60*time.Second, 10*time.Second, 5*time.Minute)

	alertChan := make(chan alert.Alert, 100)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handler: parse each incoming line and save to storage
	handler := func(line string) {
		event, err := auth.ParseLine(line)
		if err != nil {
			// Skip unrecognized lines (e.g., [preauth] messages)
			return
		}

		event.Timestamp = time.Now()

		if err := store.Save(event); err != nil {
			fmt.Fprintf(os.Stderr, "save error: %v\n", err)
			return
		}

		fmt.Printf("Saved event: %s | %s from %s\n",
			event.Timestamp.Format("15:04:05"), event.Outcome, event.SourceIP)
	}

	// Start collector in background goroutine
	go func() {
		fmt.Println("SentinelGo listening on UDP :1514...")
		if err := collector.Start(ctx, handler); err != nil {
			fmt.Fprintf(os.Stderr, "collector error: %v\n", err)
		}
	}()

	// Start detector in background goroutine
	go func() {
		if err := detector.Start(ctx, alertChan); err != nil {
			fmt.Fprintf(os.Stderr, "detector error: %v\n", err)
		}
		close(alertChan) // ← Signal no more alerts will be sent
	}()

	// Alert handler (console notifier for now)
	go func() {
		for a := range alertChan {
			fmt.Printf("🚨 ALERT [%s]: %s - %s (user: %s, events: %d)\n",
				a.Severity, a.Rule, a.Message, a.User, a.EventCount)
		}
	}()

	// Routing
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/events", GetEventHandler)
		v1.GET("/alerts", GetAlertHandler)
	}

	// Wait for Ctrl+C (SIGINT) or SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	cancel()

	// Close collector and detector (they will exit their Start() loops)
	collector.Close()
	detector.Close()

	// Give goroutines a moment to finish cleanly
	time.Sleep(500 * time.Millisecond)
}