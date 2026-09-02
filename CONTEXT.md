# SentinelGo Context & Audit Summary

## Project Overview

**SentinelGo** is a lightweight, self-hosted Security Information and Event Management (SIEM) system built in Go. Its core purpose is to ingest raw security logs, normalize them into a unified schema, store events, run threat detection rules (such as SSH brute-force identification), and expose alerts and raw events via notifications and a web API.

### Tech Stack
- **Language & Runtime**: Go (`go 1.26.5` specified in `go.mod`).
- **Core Standard Libraries**: `net` (UDP listener), `sync` (thread safety), `regexp` (log parsing), `context`, `bufio`, `time`.
- **HTTP / REST Framework**: `github.com/gin-gonic/gin` (in progress).
- **Build & Deployment**: `Makefile`, Multi-stage `Dockerfile` (~20MB minimal container), GitHub Actions CI (`.github/workflows/`).

---

## Key Navigation & Entry Points

| Path / File | Purpose & Component Interaction |
| :--- | :--- |
| `cmd/sentinelgo/main.go` | **Application Entrypoint**: Initializes `InMemoryStorage`, launches `SysLogCollector` on UDP `:1514`, runs `BruteForceDetector`, spins up console alert handler, sets up Gin HTTP router, and manages graceful shutdown. |
| `internal/logs/model.go` | **Core Event Schema**: Defines `logs.Event` (`Timestamp`, `SourceIP`, `User`, `EventType`, `Outcome`, `Host`, `RawMessage`). |
| `internal/alert/alert.go` | **Alert Model**: Defines `alert.Alert` (`ID`, `Timestamp`, `Rule`, `Severity`, `SourceIP`, `User`, `Message`, `EventCount`). |
| `internal/collector/collector.go` | **Collector Interface**: Standard contract (`Start`, `Close`) for log ingestion sources. |
| `internal/collector/syslog/syslog.go` | **UDP Syslog Listener**: Implementation of `Collector` listening on UDP socket with read deadlines and non-blocking context handling. |
| `internal/parser/auth/auth_parser.go` | **Auth Log Parser**: Converts raw SSH syslog lines into normalized `logs.Event` structs using regex. Handles month parsing and status mapping (`Accepted`/`Failed`). |
| `internal/storage/storage.go` | **Storage Interface**: Defines `Storage` contract (`Save`, `Query`, `Count`, `Close`) and `QueryFilter` struct. |
| `internal/storage/memory/memory.go` | **In-Memory Store**: Thread-safe (`sync.RWMutex`) implementation of `Storage` supporting filtering (time, IP, user, outcome, event type) and offset/limit pagination. |
| `internal/detector/detector.go` | **Detector Interface**: Contract for background detection rules. |
| `internal/detector/bruteforce/bruteforce.go` | **Brute-Force Detector**: Periodic ticker scan querying failed auth events per source IP over a sliding window, identifying most targeted user, and enforcing cooldowns per IP. |
| `internal/api/server.go` | **REST API Server**: Endpoint handlers for querying events and alerts (currently in-progress). |
| `Dockerfile` | **Container Build**: Multi-stage build producing an optimized ~20MB Alpine binary container. |

---

## Completed Features

- [x] **Unified Event Schema**: Defined `logs.Event` structure in `internal/logs/model.go`.
- [x] **SSH Auth Log Parsing**: Regular expression parser capable of handling successful logins, failed attempts, and invalid users from syslog format. Evaluated with test coverage in `internal/parser/auth/auth_parser_test.go`.
- [x] **UDP Syslog Collector**: Asynchronous UDP listener running in background goroutines with context cancellation and error handling (`internal/collector/syslog/syslog.go`).
- [x] **Thread-Safe In-Memory Storage**: Concurrent read/write storage engine using `sync.RWMutex` supporting offset/limit pagination and multi-field filtering (`internal/storage/memory/memory.go`). Tested in `memory_test.go`.
- [x] **Brute-Force Rule Detector**: Configurable window and threshold scanner detecting high-frequency login failures from single IP addresses, calculating target usernames, and suppressing duplicate alerts via cooldown timers (`internal/detector/bruteforce/bruteforce.go`).
- [x] **Console Alert Channel**: Real-time alerting pipeline emitting formatted alert outputs to stdout.
- [x] **Lightweight Docker Image**: Multi-stage Dockerfile paring image size down to ~20MB.
- [x] **Automated CI Workflow**: GitHub Actions configuration running `golangci-lint` and `go test`.

---

## Current Work-In-Progress

### Active Commit State
- **Branch**: `main` (up to date with `origin/main`).
- **Latest Commit**: `bfd13df` (`wip: working on endpoints`).
- **Untracked File**: `GEMINI.md`.

### Active Development Focus
1. **REST API Implementation (`/v1/events` & `/v1/alerts`)**:
   - `cmd/sentinelgo/main.go` registers `/v1/events` and `/v1/alerts` routes on a Gin engine.
   - `internal/api/server.go` holds stub implementations for `GetEventHandler` and `GetAlertHandler`.
2. **Current Blockers / Compilation Errors**:
   - Dependency `github.com/gin-gonic/gin` is imported in `main.go` and `server.go` but missing from `go.mod`.
   - `internal/api/server.go` handlers currently return invalid function signatures (`func()` instead of `gin.HandlerFunc`).
   - `cmd/sentinelgo/main.go` calls un-prefixed `GetEventHandler` without passing the `store` argument.

---

## Next Action Items

1. **Resolve REST API Dependencies & Implementation**:
   - Add Gin dependency to `go.mod` via `go get github.com/gin-gonic/gin` & `go mod tidy`.
   - Complete handler logic in `internal/api/server.go` to query `InMemoryStorage` and return JSON responses (`c.JSON(http.StatusOK, ...)`).
   - Update `cmd/sentinelgo/main.go` to import `internal/api` and bind `api.GetEventHandler(store)` and `api.GetAlertHandler(store)`.
2. **Verify Build & Run End-to-End Tests**:
   - Verify all packages compile and pass with `go test ./...`.
   - Run `make run` and verify HTTP endpoints (`GET /v1/events` and `GET /v1/alerts`) respond correctly.
3. **Web Dashboard & Visualization**:
   - Develop frontend UI or Grafana dashboard integration for real-time security log search and alert display.
4. **Persistent Storage & Multi-Source Collectors**:
   - Introduce persistent storage backend (e.g., SQLite / PostgreSQL).
   - Expand collectors and parsers to handle firewall logs and cloud audit trails.
