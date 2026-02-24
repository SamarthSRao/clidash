# Multi-stage Dockerfile for CLIDash components
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o optimizer ./cmd/optimizer/main.go
RUN go build -o agent ./cmd/agent/main.go
RUN go build -o dashboard ./main.go

# Optimizer Image
FROM alpine:latest AS optimizer
WORKDIR /root/
COPY --from=builder /app/optimizer .
EXPOSE 8080
CMD ["./optimizer"]

# Agent Image
FROM alpine:latest AS agent
WORKDIR /root/
COPY --from=builder /app/agent .
EXPOSE 5775
CMD ["./agent"]
