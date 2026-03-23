.PHONY: mocks test coverage test-complete test-db down-db
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

export DATABASE_URL = postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable

test-db:
	@if docker ps -a --format '{{.Names}}' | grep -q '^test-postgres$$'; then \
		docker start test-postgres; \
	else \
		docker run --name test-postgres \
		  -e POSTGRES_PASSWORD=postgres \
		  -e POSTGRES_USER=postgres \
		  -e POSTGRES_DB=testdb \
		  -p 5432:5432 \
		  -d postgres:16; \
	fi

down-db:
	docker stop test-postgres
	docker rm test-postgres
