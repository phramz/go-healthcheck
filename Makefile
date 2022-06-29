BUILD_COMMIT:=$(or $(BUILD_COMMIT),$(shell git rev-parse --short HEAD))

.PHONY: vendors
vendors:
	go mod download

.PHONY: format
format:
	go fmt ./...

.PHONY: lint
lint:
	go vet ./...

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: run
run:
	go run cmd/healthcheck.go

.PHONY: build
build: vendors
	mkdir -p build/bin
	CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o build/bin/healthcheck cmd/healthcheck.go

.PHONY: build-docker
build-docker:
	docker build -f Dockerfile \
		--build-arg "BUILD_COMMIT=$(BUILD_COMMIT)" \
		-t healthcheck \
		-t healthcheck:latest \
		-t healthcheck:$(BUILD_COMMIT) \
		 .
