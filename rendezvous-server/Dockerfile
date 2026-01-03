# # ----- Build stage -----
# # Build stage
# FROM golang:1.21-alpine as builder
# WORKDIR /app
# COPY . .
# RUN go build -o main ./cmd/api

# # Runtime stage
# FROM alpine:latest
# WORKDIR /app
# COPY --from=builder /app/main .
# EXPOSE 8085
# CMD ["./main"]

# ----- Build stage -----
FROM golang:1.22-alpine AS builder

# Install git and dependencies
RUN apk add --no-cache git

# Set working directory inside container
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code
COPY . .

# Build the app (target main.go inside cmd/api)
RUN go build -o app ./cmd/api

# ----- Final stage -----
FROM alpine:latest

# Needed for TLS and root CA support
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy built binary from builder stage
COPY --from=builder /app/app .

# Expose app port (change if your app uses another)
EXPOSE 8088

# Start app
CMD ["./app"]