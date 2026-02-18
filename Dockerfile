FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/exe ./cmd/app/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

# устанавливаем migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

WORKDIR /app
COPY --from=builder /app/exe .
COPY --from=builder /app/internal/core/migrations ./internal/core/migrations

CMD ["./exe"]