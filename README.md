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