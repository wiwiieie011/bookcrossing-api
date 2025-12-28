# syntax=docker/dockerfile:1.6
# =========================================
# Builder stage
# =========================================
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

# Cache deps
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy sources
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/app ./cmd/bookcrossing

RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/seed ./cmd/seed

# =========================================
# Runtime stage
# =========================================
FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates tzdata && adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/app /app/app
COPY --from=builder /app/seed /app/seed

USER appuser

EXPOSE 1010
ENV GIN_MODE=release

CMD ["/app/app"]
