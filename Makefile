API_PATH = ./cmd/api
API_BIN_NAME = expense-tracker-api

.PHONY: test
test:
	go test -v -race -buildvcs ./...

.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

.PHONY: build
build:
	@go build -o=/tmp/bin/$(API_BIN_NAME) $(API_PATH)

.PHONY: run
run: build
	@/tmp/bin/$(API_BIN_NAME)
