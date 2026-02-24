# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 👇 ADD THIS
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o app ./cmd/API

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]