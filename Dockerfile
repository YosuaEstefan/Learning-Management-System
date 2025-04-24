# Dockerfile
# Stage 1: build
FROM golang:1.24.1-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Compile binary from root main.go
RUN go build -o lms ./main.go

# Stage 2: create minimal image
FROM alpine:3.17

# Install runtime dependencies
RUN apk update && \
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.18/main" > /etc/apk/repositories && \
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.18/community" >> /etc/apk/repositories && \
    apk update && \
    apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Create upload directories
RUN mkdir -p /home/appuser/uploads/materials /home/appuser/uploads/submissions
WORKDIR /home/appuser

# Copy the compiled binary
COPY --from=builder /app/lms .

# Expose application port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["./lms"]