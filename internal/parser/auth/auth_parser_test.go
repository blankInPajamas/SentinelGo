package auth_test

import (
	"strings"
	"testing"

	"github.com/blankInPajamas/SentinelGo/internal/parser/auth"
)

func TestParseLine_AcceptedPassword(t *testing.T) {
	line := "Aug 17 20:58:11 web01 sshd[10234]: Accepted password for deploy from 198.51.100.7 port 51520 ssh2"

	event, err := auth.ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if event.Outcome != "success" {
		t.Errorf("Expected outcome = success, got %s", event.Outcome)
	}

	if event.Host != "web01" {
		t.Errorf("Expected Host = web01, got %s", event.Host)
	}

	if event.SourceIP != "198.51.100.7" {
		t.Errorf("Expected SourceIP = 198.51.100.7, got %s", event.SourceIP)
	}

	if event.User != "deploy" {
		t.Errorf("Expected user deploy, got %s", event.User)
	}
}

func TestParseLine_FailedPasswordInvalidUser(t *testing.T) {
	line := "Aug 17 20:59:02 web01 sshd[10245]: Failed password for invalid user admin from 203.0.113.5 port 41332 ssh2"

	event, err := auth.ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if event.Outcome != "failure" {
		t.Errorf("expected outcome=failure, got %s", event.Outcome)
	}
	if event.User != "admin" {
		t.Errorf("expected user=admin, got %s", event.User)
	}
	if event.SourceIP != "203.0.113.5" {
		t.Errorf("expected ip=203.0.113.5, got %s", event.SourceIP)
	}
}

func TestParseLine_FailedPasswordValidUsername(t *testing.T) {
	line := "Aug 17 21:00:45 web01 sshd[10260]: Failed password for root from 203.0.113.5 port 41360 ssh2"

	event, err := auth.ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if event.User != "root" {
		t.Errorf("expected user=root, got %s", event.User)
	}
	if event.Outcome != "failure" {
		t.Errorf("expected outcome=failure, got %s", event.Outcome)
	}
}

func TestParseLine_AcceptedPublickey(t *testing.T) {
	line := "Aug 17 21:02:10 web01 sshd[10275]: Accepted publickey for deploy from 198.51.100.7 port 51600 ssh2"

	event, err := auth.ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Outcome != "success" {
		t.Errorf("expected outcome=success, got %s", event.Outcome)
	}
}

func TestParseLine_EmptyLine(t *testing.T) {
	_, err := auth.ParseLine("")

	if err == nil {
		t.Errorf("expected ErrNoMatch for empty line, got %v", err)
	}
}

func TestParseLine_UnrecognizedFormat(t *testing.T) {
	line := "Aug 17 21:25:47 web01 sshd[10370]: Connection closed by authenticating user root 203.0.113.5 port 41700 [preauth]"

	_, err := auth.ParseLine(line)
	if err == nil {
		t.Errorf("expected an error for unrecognized format, got nil")
	}
}

func TestParseReader_SkipsUnrecognizedLines(t *testing.T) {
	input := strings.Join([]string{
		"Aug 17 20:58:11 web01 sshd[10234]: Accepted password for deploy from 198.51.100.7 port 51520 ssh2",
		"Aug 17 21:25:47 web01 sshd[10370]: Connection closed by authenticating user root 203.0.113.5 port 41700 [preauth]",
		"Aug 17 20:59:02 web01 sshd[10245]: Failed password for invalid user admin from 203.0.113.5 port 41332 ssh2",
	}, "\n")

	events, err := auth.ParseReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 parsed events (1 skipped), got %d", len(events))
	}
}

func TestParseReader_BruteForceBurstFromSampleData(t *testing.T) {
	input := strings.Join([]string{
		"Aug 17 20:59:02 web01 sshd[10245]: Failed password for invalid user admin from 203.0.113.5 port 41332 ssh2",
		"Aug 17 20:59:05 web01 sshd[10246]: Failed password for invalid user admin from 203.0.113.5 port 41335 ssh2",
		"Aug 17 20:59:09 web01 sshd[10247]: Failed password for invalid user admin from 203.0.113.5 port 41338 ssh2",
	}, "\n")

	events, err := auth.ParseReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}

	for _, e := range events {
		if e.SourceIP != "203.0.113.5" {
			t.Errorf("expected all events from 203.0.113.5, got %s", e.SourceIP)
		}
		if e.Outcome != "failure" {
			t.Errorf("expected failure outcome, got %s", e.Outcome)
		}
	}
}