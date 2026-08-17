package main

import (
	"fmt"
	"os"

	"github.com/blankInPajamas/SentinelGo/internal/parser/auth"
)

func main() {
	f, err := os.Open("test/sample_logs/auth.log")

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open the sample log: %v\n", err)
	}

	defer f.Close()

	events, err := auth.ParseReader(f)

	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed %d events \n\n", len(events))

	sc, fail := 0,0

	for _, e := range events {
		fmt.Printf("%s | %-8s | user=%-10s ip=%-15s host=%s\n",
			e.Timestamp.Format("15:04:05"), e.Outcome, e.User, e.SourceIP, e.Host)

		if e.Outcome == "success" {
			sc++
		} else {
			fail++
		}
	}

	fmt.Printf("\nsummary: %d success, %d failure\n", sc, fail)
}