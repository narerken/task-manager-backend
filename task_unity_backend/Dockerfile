FROM golang:1.25.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o task-unity-backend ./cmd/

FROM alpine:3.22

WORKDIR /app

COPY --from=build /app/task-unity-backend ./task-unity-backend
COPY config.json ./config.json

EXPOSE 8080

CMD ["./task-unity-backend"]
