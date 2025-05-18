FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

#  Build a fully static binary for Alpine compatibility
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o log-sender .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/log-sender .

RUN chmod +x log-sender

ENTRYPOINT ["./log-sender"]
