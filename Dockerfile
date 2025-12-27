FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies (including CGO requirements for SQLite)
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Enable CGO for SQLite support
ENV CGO_ENABLED=1

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -a -installsuffix cgo -o whatsapp-worker ./cmd/whatsapp-worker

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata sqlite-libs

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/whatsapp-worker .

# Expose health check port
EXPOSE 8080

# Run the binary
CMD ["./whatsapp-worker"]

