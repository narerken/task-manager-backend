FROM golang:1.25.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o attendance-session-service ./cmd/

FROM alpine:3.22

WORKDIR /app

COPY --from=build /app/attendance-session-service ./attendance-session-service
COPY config.json ./config.json

EXPOSE 8081

CMD ["./attendance-session-service"]
