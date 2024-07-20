.PHONY: lint
lint:
	golangci-lint run --config .golangci-lint.yaml

.PHONY: build
build:
	go build -o ./.bin/app cmd/main/main.go

.PHONY: run
run: build
	./.bin/app

.PHONY: test
test:
	go test -v ./...

.PHONY: test-cover
test-cover:
	go test -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: clean
clean:
	rm -rf ./.bin