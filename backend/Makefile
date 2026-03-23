.PHONY: mocks test coverage test-complete
mocks:
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make mocks FILE=path/to/file.go"; \
		exit 1; \
	fi
	@./scripts/makemocks.sh "$(FILE)"

test:
	go test ./...

coverage:
	go test ./... -cover

test-complete:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
