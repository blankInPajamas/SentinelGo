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

**Next steps in progress:**

5. Write unit tests for the parser to ensure it handles all auth log variants and gracefully skips unrecognized lines.
6. Build a basic in-memory storage layer to persist parsed events.
7. Implement the first detection rule: brute-force detection (N failed logins from the same source IP within T seconds).
8. Add a console notifier to print alerts when the brute-force rule fires.
9. Expose a minimal REST API (`GET /events`, `GET /alerts`) to verify data end-to-end.

Once this loop works — raw auth log in, brute-force alert out — the plan is to expand to additional log sources (firewall, cloud audit logs, endpoint/EDR) using the same collector/parser/storage/detector interfaces, without needing to rework the pipeline itself.

## Project Structure

```
sentinelgo/
├── cmd/sentinelgo/       # entry point (main.go)
├── internal/
│   ├── logs/             # shared Event schema (model.go)
│   └── parser/
│      └── auth/         # SSH auth log parser
├── test/sample_logs/     # sample logs for local testing
├── go.mod
└── README.md
```