DIR_BIN = ./bin

APP_NAME = $(shell basename $(PWD))

.PHONY: build install install-tools install-packages test fmt lint start up down

default: compile start

build: generate lint test compile

compile:
	go build -v -o ${DIR_BIN}/${APP_NAME} \
	./main.go

generate:
	oapi-codegen -config ./api/oapi-codegen.yaml -o ./internal/api/api.gen.go ./api/openapi.yaml

install-tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest && \
	go install golang.org/x/tools/cmd/goimports@latest && \
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

install-packages:
	go mod tidy && \
	go mod vendor

install: install-tools install-packages

test:
	go test ./

start:
	${DIR_BIN}/${APP_NAME}

fmt:
	go fmt ./... && \
	goimports -w -local ./internal ./pkg main.go

lint:
	go vet ./... && \
	staticcheck ./...

up:
	docker-compose up -d

down:
	docker-compose down
