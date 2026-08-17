package main

import (
	"fmt"
	"os"

	"github.com/blankInPajamas/SentinelGo/internal/parser/auth"
	"github.com/blankInPajamas/SentinelGo/internal/storage"
	"github.com/blankInPajamas/SentinelGo/internal/storage/memory"
)

func main() {
	f, err := os.Open("test/sample_logs/auth.log")

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open the sample log: %v\n", err)
	}

	defer f.Close()

	store := memory.New()
	defer store.Close()

	events, err := auth.ParseReader(f)

	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed %d events \n\n", len(events))

	
	for _, e := range events {
		if err := store.Save(e); err != nil {
			fmt.Fprintf(os.Stderr, "save error: %v\n", err)
			os.Exit(1)
		}
	}
	
	all, err := store.Query(storage.QueryFilter{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "query error: %v\n", err)
		os.Exit(1)
	}

	sc, fail := 0,0

	for _, e := range all {
		fmt.Printf("%s | %-8s | user=%-10s ip=%-15s host=%s\n",
			e.Timestamp.Format("15:04:05"), e.Outcome, e.User, e.SourceIP, e.Host)

		if e.Outcome == "success" {
			sc++
		} else {
			fail++
		}
	}

	fmt.Printf("\nsummary: %d success, %d failure\n", sc, fail)
	fmt.Printf("storage count: %d events\n", store.Count())
}