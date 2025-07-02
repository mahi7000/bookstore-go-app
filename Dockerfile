# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /root/

COPY --from=builder /go/bin/goose /usr/local/bin/goose
# Copy binary from build stage
COPY --from=builder /app/main .
COPY --from=builder /app/sql ./sql
COPY entrypoint.sh .

# Run the binary
ENTRYPOINT ["./entrypoint.sh"]