# ── Build Stage ──────────────────────────────────────────────────────────────
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o tech-news-agent ./cmd/server

# ── Final Stage ───────────────────────────────────────────────────────────────
FROM alpine:3.19

# Run as non-root
RUN addgroup -S app && adduser -S app -G app
WORKDIR /app

COPY --from=builder /app/tech-news-agent .

USER app

ENTRYPOINT ["./tech-news-agent"]
