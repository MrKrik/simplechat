FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o apigateway ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/apigateway .
COPY --from=builder /app/config ./config

EXPOSE 8082

CMD ["./apigateway", "-config", "/root/config/local.yaml"]
