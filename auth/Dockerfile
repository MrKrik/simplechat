# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux \
    go build -ldflags "-linkmode external -extldflags -static" -o auth ./cmd/auth/main.go
RUN CGO_ENABLED=1 GOOS=linux \
    go build -ldflags "-linkmode external -extldflags -static" -o migrator ./cmd/migrator/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/auth .
COPY --from=builder /app/migrator .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations
RUN mkdir -p /root/storage

EXPOSE 44044

CMD ["sh", "-c", "./migrator -storage-path /root/storage/auth.db -migrations-path /root/migrations && ./auth -config /root/config/local.yaml"]
