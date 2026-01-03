# ----- Build stage -----
FROM golang:1.23.2-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# ⚠️ Make sure this path matches your main package!
RUN go build -o app ./

# ----- Final stage -----
FROM alpine:3.19

RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8088
CMD ["./app"]