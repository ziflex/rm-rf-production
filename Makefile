DIR_BIN = ./bin
APP_NAME = $(shell basename $(PWD))
OAPI := $(shell go env GOPATH)/bin/oapi-codegen

export GOOS=linux
export GOARCH=amd64

export DB_PORT ?= 5432
export DB_NAME ?= mydb
export DB_USER ?= user
export DB_PASS ?= password

.PHONY: clean build install install-tools install-packages test fmt lint start up down

default: compile start

clean:
	rm -rf ${DIR_BIN}/*

build: generate lint test compile

compile:
	go build -v -o ${DIR_BIN}/${APP_NAME} \
	./main.go

generate:
	oapi-codegen -config ./spec/oapi-codegen.yaml -o ./internal/api/api.gen.go ./spec/openapi.yaml

install-tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest && \
	go install golang.org/x/tools/cmd/goimports@latest && \
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

install-packages:
	go mod tidy && \
	go mod vendor

install: install-tools install-packages

test:
	go test ./internal/... ./pkg/...

start:
	${DIR_BIN}/${APP_NAME}

fmt:
	go fmt ./... && \
	goimports -w -local ./internal ./pkg main.go

lint:
	go vet ./... && \
	staticcheck ./...

up:
	docker compose up -d --build

refresh:
	docker compose down && \
	docker compose build --no-cache && \
	docker compose up -d --build

down:
	docker compose down
