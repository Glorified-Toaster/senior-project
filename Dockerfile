# ── Stage 1: Build ──────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates curl libc6-compat

# Install templ CLI
RUN go install github.com/a-h/templ/cmd/templ@latest

# Install standalone tailwindcss CLI
RUN arch=$(uname -m) && \
    if [ "$arch" = "x86_64" ]; then arch="x64"; elif [ "$arch" = "aarch64" ]; then arch="arm64"; fi && \
    curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-$arch && \
    chmod +x tailwindcss-linux-$arch && \
    mv tailwindcss-linux-$arch /usr/local/bin/tailwindcss

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy entire source
COPY . .

# Generate templ files
RUN templ generate

# Build tailwind CSS
RUN tailwindcss -i ./web/static/src/css/input.css -o ./web/static/css/output.css --minify

# Build the Go binary (static, no CGO)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/api/main.go

# ── Stage 2: Runtime ───────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata postgresql-client

WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/server .

# Copy runtime assets
COPY --from=builder /app/web/static ./web/static
COPY --from=builder /app/web/templates ./web/templates
COPY --from=builder /app/config ./config
COPY --from=builder /app/db/migrations ./db/migrations

# Create directories the app expects
RUN mkdir -p logs certs web/static/uploads

EXPOSE 8443

CMD ["./server"]
