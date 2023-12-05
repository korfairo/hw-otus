BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
POSTGRES_DSN="postgresql://postgres:password@127.0.0.1:5432/calendar?sslmode=disable"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

lint-install:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: lint-install
	golangci-lint run ./...

migrate-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

migrate: migrate-install
	goose -dir migrations postgres $(POSTGRES_DSN) up

.PHONY: build run build-img run-img version test lint-install lint migrate-install migrate
