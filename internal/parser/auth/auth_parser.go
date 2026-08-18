package auth

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/blankInPajamas/SentinelGo/internal/logs"
)

var monthMap = map[string]time.Month{
	"Jan": time.January, "Feb": time.February, "Mar": time.March,
	"Apr": time.April, "May": time.May, "Jun": time.June,
	"Jul": time.July, "Aug": time.August, "Sep": time.September,
	"Oct": time.October, "Nov": time.November, "Dec": time.December,
}

var sshLogPattern = regexp.MustCompile(
	`^(?P<month>\w{3})\s+(?P<day>\d{1,2})\s+(?P<time>\d{2}:\d{2}:\d{2})\s+` +
		`(?P<host>\S+)\s+sshd\[\d+\]:\s+` +
		`(?P<status>Accepted|Failed)\s+(?P<method>password|publickey)\s+for\s+` +
		`(?:invalid user\s+)?(?P<user>\S+)\s+from\s+(?P<ip>\S+)\s+port\s+(?P<port>\d+)`,
)

func ParseReader(r io.Reader) ([]logs.Event, error) {
	var events []logs.Event
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		event, err := ParseLine(line)
		if err != nil {
			continue
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return events, fmt.Errorf("scanning input: %w", err)
	}

	return events, nil
}

func ParseLine(line string) (logs.Event, error) {
	match := sshLogPattern.FindStringSubmatch(line)
	if match == nil {
		return logs.Event{}, fmt.Errorf("line does not match SSH pattern")
	}

	groups := make(map[string]string)
	for i, name := range sshLogPattern.SubexpNames() {
		if i != 0 && name != "" {
			groups[name] = match[i]
		}
	}

	ts, err := buildTimestamp(groups["month"], groups["day"], groups["time"])
	if err != nil {
		return logs.Event{}, fmt.Errorf("parsing timestamp: %w", err)
	}

	outcome := "failure"
	if groups["status"] == "Accepted" {
		outcome = "success"
	}

	return logs.Event{
		Timestamp:  ts,
		SourceIP:   groups["ip"],
		User:       groups["user"],
		EventType:  "auth",
		Host:       groups["host"],
		Outcome:    outcome,
		RawMessage: line,
	}, nil
}

func buildTimestamp(month, day, clk string) (time.Time, error) {
	m, ok := monthMap[month]
	if !ok {
		return time.Time{}, fmt.Errorf("invalid month: %s", month)
	}

	d, err := strconv.Atoi(day)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day: %w", err)
	}

	parsedClock, err := time.Parse("15:04:05", clk)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time: %w", err)
	}

	year := time.Now().Year()
	return time.Date(
		year, m, d,
		parsedClock.Hour(), parsedClock.Minute(), parsedClock.Second(),
		0, time.Local,
	), nil
}
