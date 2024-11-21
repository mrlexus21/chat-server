FROM golang:1.25.5-alpine AS builder

COPY . /app/auth
WORKDIR /app/auth

RUN go mod download && \
    go build -o /app/auth/bin/chat_server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/auth/bin/chat_server .
COPY .env .

CMD ["./chat_server", "-config-path", ".env"]