# =============================================================================
#  Local testing
# =============================================================================

# Print uncovered lines in unit tests (coverage)
cover:
	go-carpet -mincov 99.9
# Download go modules
getmod:
	@go mod download
# Run unit tests
test: getmod
	@go test -race -cover ./...
# Run golangci-lint (lint and static analysis)
lint:
	@golangci-lint run && echo "OK"
# Update go.mod and go.sum to latest version
update:
	go get -u ./...
	go mod tidy -go=1.17
# Run govulncheck (vulnerability check)
vuln:
	govulncheck ./...

# =============================================================================
#  Docker based testing
# =============================================================================

# Build container images
docker_build: docker_pull
	docker compose --file ./.github/docker-compose.yml build
# Run unit tests with golang 1.17
docker_go117:
	docker compose --file ./.github/docker-compose.yml run --rm v1_17
# Run unit tests with golang 1.18
docker_go118:
	docker compose --file ./.github/docker-compose.yml run --rm v1_18
# Run unit tests with golang 1.19
docker_go119:
	docker compose --file ./.github/docker-compose.yml run --rm v1_19
# Run unit tests with latest golang
docker_go:
	docker compose --file ./.github/docker-compose.yml run --rm latest
# Run golangci-lint (lint and static analysis)
docker_lint:
	docker compose --file ./.github/docker-compose.yml run --rm lint
# Pull image one-by-one for stability
docker_pull:
	@docker pull golang:alpine
	@docker pull golang:1.19-alpine
	@docker pull golang:1.18-alpine
	@docker pull golang:1.17-alpine
	@docker pull golangci/golangci-lint:latest
# Update go.mod and go.sum
docker_update:
	docker compose --file ./.github/docker-compose.yml run --rm update
# Run govulncheck (vulnerability check)
docker_vuln:
	docker compose --file ./.github/docker-compose.yml run --rm vuln
