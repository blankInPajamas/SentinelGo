# Stage 01
FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/sentinelgo

# Stage 02
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 1514

CMD ["./main"]

