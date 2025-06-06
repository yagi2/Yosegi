# Build stage
FROM golang:1.21-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o yosegi \
    ./main.go

# Final stage
FROM scratch

# Copy ca-certificates from builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/yosegi /usr/local/bin/yosegi

# Copy documentation
COPY --from=builder /app/README.md /README.md
COPY --from=builder /app/LICENSE /LICENSE
COPY --from=builder /app/SECURITY.md /SECURITY.md

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/yosegi"]

# Default command
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.title="Yosegi"
LABEL org.opencontainers.image.description="Interactive git worktree management tool with beautiful TUI"
LABEL org.opencontainers.image.url="https://github.com/yagi2/Yosegi"
LABEL org.opencontainers.image.source="https://github.com/yagi2/Yosegi"
LABEL org.opencontainers.image.licenses="MIT"