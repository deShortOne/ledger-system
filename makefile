# Docker image to run shell and go utility functions in
WORKER_IMAGE = golang:1.15-alpine3.13
# Docker image to generate OAS3 specs
OAS3_GENERATOR_DOCKER_IMAGE = openapitools/openapi-generator-cli:latest-release

BINARY=bin/server
MAIN=./cmd/server
SWAG_MAIN=./cmd/server/main.go

VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SQLC_CONFIGS := $(shell find internal -name sqlc.yaml)

LDFLAGS=-ldflags "-X main.version=$(VERSION) \
                  -X main.commit=$(COMMIT) \
                  -X main.date=$(DATE)"

TOOLS = \
	github.com/swaggo/swag/cmd/swag@v1.16.6 \
	github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 \
	github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.10.1

.PHONY: build run tidy generate sqlc swagger \
        test test-nocache coverage \
        lint fmt vet verify clean tools

## Dev Setup
setup: tools tidy

tools:
	@echo "Installing dev tools..."
	@for tool in $(TOOLS); do \
		go install $$tool; \
	done

## Build
build: generate
	go build $(LDFLAGS) -o $(BINARY) $(MAIN)

run: build
	./$(BINARY)

## Dependencies
tidy:
	go mod tidy

## Code generation
generate: sqlc swagger-gen-v3

sqlc:
	@for config in $(SQLC_CONFIGS); do \
		sqlc generate -f $$config; \
	done

swagger:
	swag init -g $(SWAG_MAIN) -o ./docs

# Generate OAS3 from swaggo/swag output since that project doesn't support it
# TODO: Remove this if V3 spec is ever returned from that project
swagger-gen-v3: swagger
	@echo "[OAS3] Converting Swagger 2-to-3 (yaml)"
	@docker run --rm -v $(PWD)/docs:/work $(OAS3_GENERATOR_DOCKER_IMAGE) \
	  generate -i /work/swagger.yaml -o /work/v3 -g openapi-yaml --minimal-update
	@docker run --rm -v $(PWD)/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "rm -rf /work/.openapi-generator"
	@echo "[OAS3] Copying openapi-generator-ignore (json)"
	@docker run --rm -v $(PWD)/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "cp -f /work/.openapi-generator-ignore /work/openapi"
	@echo "[OAS3] Converting Swagger 2-to-3 (json)"
	@docker run --rm -v $(PWD)/docs:/work $(OAS3_GENERATOR_DOCKER_IMAGE) \
	  generate -s -i /work/swagger.json -o /work/v3/openapi -g openapi --minimal-update
	@echo "[OAS3] Cleaning up generated files"
	@docker run --rm -v $(PWD)/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "mv -f /work/openapi/openapi.json /work ; mv -f /work/openapi/openapi.yaml /work ; rm -rf /work/openapi"

## Formatting & Linting
fmt:
	gofmt -w .

fmt-check:
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Code not formatted:"; \
		gofmt -l .; \
		exit 1; \
	fi

vet:
	go vet ./...

lint:
	golangci-lint run

## Testing
test:
	go test ./...

test-nocache:
	go test -count=1 ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## CI Verification
verify: tidy fmt-check vet lint test generate
	git diff --exit-code

## Cleanup
clean:
	rm -rf bin coverage.out
