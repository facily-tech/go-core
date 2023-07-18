.PHONY: lint
lint:
	@golangci-lint run --timeout 5m

.PHONY: test
test:
	@go test ./... -race -covermode=atomic -coverprofile=coverage.out
