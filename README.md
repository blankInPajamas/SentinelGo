# SentinelGo

A lightweight, self-hosted SIEM (Security Information and Event Management) system written in Go. SentinelGo ingests logs from multiple sources, normalizes them into a common schema, correlates events to detect threats, and generates real-time alerts through a searchable dashboard.

## What It Does

SentinelGo centralizes security-relevant data from across an environment so that threats hidden in isolated logs become visible when correlated together. At its core, the project handles:

- **Log collection** — pulling in raw events from sources like syslog, cloud audit trails, and endpoint logs.
- **Normalization** — mapping varied log formats into a consistent internal event schema.
- **Storage & search** — persisting normalized events so they can be queried quickly.
- **Correlation & detection** — applying rules to spot patterns indicative of attacks (e.g., brute-force login attempts).
- **Alerting** — surfacing detections as actionable alerts through notifications.
- **Visualization** — exposing an API/dashboard to search events and review alerts.

## Current Task

We're starting with **authentication logs** (SSH/syslog-based auth events) as the first log source, since they're easy to generate and test against, and they let us validate the full pipeline end-to-end quickly.

**Completed milestones:**

1. Defined the normalized `Event` schema (`internal/logs/model.go`).
2. Created sample auth log data (`test/sample_logs/auth.log`) with realistic SSH login patterns.
3. Implemented the SSH auth log parser (`internal/parser/auth/auth_parser.go`) that converts raw syslog lines into structured `Event` values.
4. Verified end-to-end parsing: running `go run cmd/sentinelgo/main.go` successfully parses 24 events from the sample log, correctly distinguishing successful logins from failed attempts.
5. Built in-memory storage layer (`internal/storage/memory/memory.go`) with thread-safe Save/Query operations using `sync.RWMutex`.
6. Implemented comprehensive storage tests covering filtering by time/IP/user/outcome, pagination (limit/offset), and concurrent access.
7. Added integration tests (`cmd/sentinelgo/main_test.go`) verifying parser → storage end-to-end flow.
8. Implement the first detection rule: brute-force detection (N failed logins from the same source IP within T seconds).
9. Add a console notifier to print alerts when the brute-force rule fires.

**Next steps in progress:**

10. Add alert deduplication to prevent repeated alerts for the same IP within a cooldown period.
11. Expose a minimal REST API (`GET /events`, `GET /alerts`) to verify data end-to-end.
12. Build a simple web dashboard or integrate Grafana for visualization.

Once this loop works — raw auth log in, brute-force alert out — the plan is to expand to additional log sources (firewall, cloud audit logs, endpoint/EDR) using the same collector/parser/storage/detector interfaces, without needing to rework the pipeline itself.

## Project Structure

```
sentinelgo/
├── cmd/sentinelgo/       # entry point (main.go)
├── internal/
│   ├── logs/             # shared Event schema (model.go)
│   ├── parser/
│   │   └── auth/         # SSH auth log parser
│   ├── collector/
│   │   └── syslog/       # UDP syslog listener
│   ├── storage/
│   │   ├── storage.go    # Storage interface
│   │   └── memory/       # In-memory storage implementation
│   ├── detector/
│   │   └── bruteforce/   # Brute-force detection rule
│   └── alert/            # Alert struct definition
├── test/sample_logs/     # sample logs for local testing
├── go.mod
└── README.md
```

## Quick Start

```bash
# Run the parser + storage integration against the sample auth log
go run cmd/sentinelgo/main.go
```

```bash
# Run all tests
go test -v ./...
```

Expected output: 24 parsed events with a mix of `success` and `failure` outcomes, and all tests passing.

## Live Ingestion (Syslog Collector)

```bash
# Start the collector (listens on UDP port 1514)
go run cmd/sentinelgo/main.go
```

```bash
# In another terminal, send test syslog messages:
echo "Aug 18 13:00:00 web01 sshd[1234]: Failed password for admin from 203.0.113.5 port 12345" | nc -u -w1 127.0.0.1 1514

# Expected output: "Saved event: 13:00:00 | failure from 203.0.113.5"
```

Press `Ctrl+C` to gracefully shut down the collector.

## Testing Brute-Force Detection

To see the detector in action, send 5+ failed logins from the same IP within 60 seconds:

```bash
# Start SentinelGo
go run cmd/sentinelgo/main.go
```

```bash
# In another terminal, send a burst of failed logins:
for i in {1..6}; do
  echo "Aug 18 14:00:0$i web01 sshd[200$i]: Failed password for admin from 10.0.0.99 port 4000$i" | nc -u -w1 127.0.0.1 1514
  sleep 0.5
done

# Wait up to 10 seconds for the detector to scan
# Expected output: "🚨 ALERT [high]: brute-force-ssh - Detected 6 failed logins from 10.0.0.99 in 1m0s"
```

The detector scans every 10 seconds, so alerts may take up to 10 seconds to appear after the threshold is breached.