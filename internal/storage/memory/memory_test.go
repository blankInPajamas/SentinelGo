package memory_test

import (
	"testing"
	"time"

	"github.com/blankInPajamas/SentinelGo/internal/logs"
	"github.com/blankInPajamas/SentinelGo/internal/storage"
	"github.com/blankInPajamas/SentinelGo/internal/storage/memory"
)

func TestInMemoryStorage_SaveAndQueryAll(t *testing.T) {
	s := memory.New()

	event1 := logs.Event{Timestamp: time.Now(), User: "alice", Outcome: "success"}
	event2 := logs.Event{Timestamp: time.Now(), User: "bob", Outcome: "failure"}
	s.Save(event1)
	s.Save(event2)

	results, err := s.Query(storage.QueryFilter{})

	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
}

func TestInMemoryStorage_QueryByOutcome(t *testing.T) {
	s := memory.New()

	event1 := logs.Event{Timestamp: time.Now(), User: "alice", Outcome: "success"}
	event2 := logs.Event{Timestamp: time.Now(), User: "dejong", Outcome: "success"}
	event3 := logs.Event{Timestamp: time.Now(), User: "bob", Outcome: "failure"}

	s.Save(event1)
	s.Save(event2)
	s.Save(event3)

	results, err := s.Query(storage.QueryFilter{})

	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	var count int = 0

	for _, r := range results {
		if r.Outcome == "success" {
			count++
		}
	}

	if count != 2 {
		t.Fatalf("Expected 2 success results, got %d", len(results))
	}
}

func TestInMemoryStorage_QueryBySourceIP(t *testing.T) {
	s := memory.New()

	event1 := logs.Event{Timestamp: time.Now(), SourceIP: "192.168.1.1", Outcome: "success"}
	event2 := logs.Event{Timestamp: time.Now(), SourceIP: "192.168.1.1", Outcome: "success"}
	event3 := logs.Event{Timestamp: time.Now(), SourceIP: "192.168.1.2", Outcome: "failure"}

	s.Save(event1)
	s.Save(event2)
	s.Save(event3)

	results, err := s.Query(storage.QueryFilter{
		SourceIP: "192.168.1.1",
	})

	if err != nil {
		t.Fatalf("Failed to query, %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 success results, got %d", len(results))
	}
}

func TestInMemoryStorage_QueryWithLimitAndOffset(t *testing.T) {
	
	s := memory.New()

	now := time.Now()
	for i := 0; i < 10; i++ {
		s.Save(logs.Event{Timestamp: now, User: "user"})
	}

	results, err := s.Query(storage.QueryFilter{
		Limit: 5,
	})

	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 5 {
		t.Fatalf("Expected 5 success results, got %d", len(results))
	}

	
	results, err = s.Query(storage.QueryFilter{Offset: 7})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 results with Offset=7, got %d", err)
	}

	results, err = s.Query(storage.QueryFilter{Limit: 2, Offset: 3})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 events with Limit=2, Offset=3, got %d", len(results))
	}

	results, err = s.Query(storage.QueryFilter{Offset: 15})
	if err != nil {
		t.Fatalf("Query with large offset failed: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("Expected 0 events with Offset=15, got %d", len(results))
	}
}

func TestInMemoryStorage_Count(t *testing.T) {
	
	s := memory.New()

	if s.Count() != 0 {
		t.Fatalf("Expected 0 events initially, received %d", s.Count())
	}

	s.Save(logs.Event{Timestamp: time.Now(), User: "alice"})
	s.Save(logs.Event{Timestamp: time.Now(), User: "bob"})
	s.Save(logs.Event{Timestamp: time.Now(), User: "charlie"})

	if s.Count() != 3 {
		t.Fatalf("Expected 3 events, received %d", s.Count())
	}

	for i := 0; i < 5; i++ {
		s.Save(logs.Event{Timestamp: time.Now(), User: "user"})
	}

	if s.Count() != 8 {
		t.Fatalf("Expected 8 events total, got %d", s.Count())
	}
}