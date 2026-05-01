FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service ./cmd/
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth-service .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8081
CMD ["./auth-service"]