# Stage 1: Development image (for Air hot-reload)
# Gunakan versi yang sama dengan go.mod untuk konsistensi
FROM golang:1.24-alpine AS development

RUN apk add --no-cache git curl tini autoconf

# Install Air (hot-reload tool) — repo resmi adalah air-verse/air
RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b /usr/local/bin

WORKDIR /app

# Volume akan override ini, copy hanya biar build berhasil
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p bin

CMD ["air", "-c", ".air.toml"]


# Stage 2: Production Builder
# Versi harus sama dengan development stage dan go.mod
FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOARCH=amd64 \
    GOOS=linux

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN mkdir -p .dist && \
    go build -a \
    -ldflags="-s -w -extldflags '-static'" \
    -trimpath \
    -o .dist/boilerplate \
    ./cmd/server/main.go


# Stage 3: Production Runtime
FROM alpine:latest AS production

RUN apk add --no-cache tini curl ca-certificates tzdata

WORKDIR /app

# Create non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

# Copy binary dari builder stage
COPY --from=builder --chown=appuser:appgroup /build/.dist/boilerplate /app/boilerplate

# Production perlu .config.yml — mount via Docker secret atau bind mount saat deploy
# Contoh: docker run -v /path/to/.config.yml:/app/.config.yml ...

# Switch to non-root user
USER appuser

EXPOSE 8000

# Add healthcheck endpoint
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8000/api/healthz || exit 1

# Use tini as init process untuk handle sinyal dengan benar
ENTRYPOINT ["/sbin/tini", "--"]

CMD ["./boilerplate"]