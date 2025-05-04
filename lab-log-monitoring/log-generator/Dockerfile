# Use official Go image to build the binary
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

# Build the Go app as a statically linked binary
RUN go build -o log-sender .

# Final minimal image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/log-sender .

# Set entrypoint
ENTRYPOINT ["./log-sender"]
