package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/blankInPajamas/SentinelGo/internal/collector/syslog"
	"github.com/blankInPajamas/SentinelGo/internal/parser/auth"
	"github.com/blankInPajamas/SentinelGo/internal/storage/memory"
)

func main() {
	store := memory.New()
	defer store.Close()

	collector := syslog.New(":1514")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parsing each incoming line and saving to storage
	handler := func(line string) {
		event, err := auth.ParseLine(line)
		if err != nil {
			// Skip unrecognized lines (e.g., [preauth] messages)
			return
		}

		if err := store.Save(event); err != nil {
			fmt.Fprintf(os.Stderr, "save error: %v\n", err)
			return
		}

		fmt.Printf("Saved event: %s | %s from %s\n",
			event.Timestamp.Format("15:04:05"), event.Outcome, event.SourceIP)
	}

	go func() {
		fmt.Println("SentinelGo listening on UDP :1514...")
		if err := collector.Start(ctx, handler); err != nil {
			fmt.Fprintf(os.Stderr, "collector error: %v\n", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	cancel()
	collector.Close()
}
