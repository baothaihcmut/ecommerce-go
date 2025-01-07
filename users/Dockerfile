# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main ./cmd/users/main.go
RUN go build -o migrate ./cmd/migrations/main.go

# Run Stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/internal/adapter/persistence/migrations ./migrations
COPY --from=builder /app/migrate . 
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
