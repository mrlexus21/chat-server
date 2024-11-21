FROM golang:1.23-alpine AS builder

COPY . /github.com/mrlexus21/chat-server/source/
WORKDIR /github.com/mrlexus21/chat-server/source/

RUN go mod download
RUN go build -o ./bin/chat_server cmd/server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/mrlexus21/chat-server/source/bin/chat_server .
COPY --from=builder /github.com/mrlexus21/chat-server/source/env .

CMD ["./chat_server", "-config-path", ".env"]