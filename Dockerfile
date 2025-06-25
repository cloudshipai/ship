# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ship ./cmd/ship

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates docker-cli

WORKDIR /root/

COPY --from=builder /app/ship /usr/local/bin/ship

# Create a non-root user
RUN addgroup -g 1000 ship && \
    adduser -D -u 1000 -G ship ship

USER ship

ENTRYPOINT ["ship"]