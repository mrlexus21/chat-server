include .env

LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-chat-api

generate-chat-api:
	mkdir -p pkg/chat/v1
	protoc --proto_path api/chat/v1 \
	--go_out=pkg/chat/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go.exe \
	--go-grpc_out=pkg/chat/v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc.exe \
	api/chat/v1/chat.proto

#build:
#	GOOS=linux GOARCH=amd64 go build -o chat_server cmd/server/main.go

#copy-to-server:
#	scp service_linux root@188.130.207.122:

#docker-build-and-push:
#	docker buildx build --no-cache --platform linux/amd64 -t cr.yandex/crpa9fe1mvd86ibqm7qn/chat-server:v0.0.1 .
#	#docker login -u mrlexus21 -p CRgAAAAA6ELvdoGnA7EL4qSbkICKeDkVldDC2OOU cr.selcloud.ru/olezhek28
#	echo y0_AgAAAAADQna9AATuwQAAAAEYZvYPAADxQEN56WxMCpDw7t3NqCQr8ikZNw|docker login --username oauth --password-stdin cr.yandex
#	docker push cr.yandex/crpa9fe1mvd86ibqm7qn/chat-server:v0.0.1

local-migration-create:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} create $(TABLE_NAME)_table sql

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v