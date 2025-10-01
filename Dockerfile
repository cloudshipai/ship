# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ship ./cmd/ship

# Final stage
FROM alpine:latest

# Install Docker CLI and other necessary tools
RUN apk --no-cache add ca-certificates docker-cli

WORKDIR /root/

COPY --from=builder /app/ship /usr/local/bin/ship

# Create a non-root user and add to docker group for socket access
RUN addgroup -g 1000 ship && \
    adduser -D -u 1000 -G ship ship && \
    mkdir -p /home/ship/.ship && \
    chown -R ship:ship /home/ship

# Note: Run container with --group-add=$(stat -c '%g' /var/run/docker.sock)
# or use root user for Docker socket access
USER ship
WORKDIR /home/ship

# Set environment variables
ENV SHIP_DISABLE_TELEMETRY=true

ENTRYPOINT ["ship"]