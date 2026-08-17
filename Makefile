.DEFAULT_GOAL := help

.PHONY: run test

run:
	go run ./cmd/sentinelgo/

test:
	go test -v ./...

help:
	@echo "Usage: make [run|test]"
