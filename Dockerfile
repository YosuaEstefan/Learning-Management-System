# Dockerfile
# Tahap 1: Membangun aplikasi
# Menggunakan image golang sebagai base image untuk membangun aplikasi
FROM golang:1.24.1-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Instal dependensi buildtime
RUN apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# mengunduh modul ke cache
COPY go.mod go.sum ./
RUN go mod download

# Salin kode sumber
COPY . .

# Mengkompilasi biner dari root main.go
RUN go build -o lms ./main.go

# Tahap 2: membuat image runtime
FROM alpine:3.17

# Install depenciensi runtime
# Menggunakan alpine sebagai base image untuk image runtime
RUN apk update && \
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.18/main" > /etc/apk/repositories && \
    echo "https://dl-cdn.alpinelinux.org/alpine/v3.18/community" >> /etc/apk/repositories && \
    apk update && \
    apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Membuat pengguna non-root
RUN adduser -D -g '' appuser
USER appuser

# Membuat direktori unggahan
RUN mkdir -p /home/appuser/uploads/materials /home/appuser/uploads/submissions
WORKDIR /home/appuser

# Salin biner yang dikompilasi
COPY --from=builder /app/lms .

# Mengekspos port aplikasi
EXPOSE 8080

# jalankan aplikasi

ENTRYPOINT ["./lms"]