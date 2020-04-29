SOURCES := $(shell find . -name '*.go')
BINARY := seed
IMAGE_TAG := dev
IMAGE := danielpacak/redis-ha-seed:$(IMAGE_TAG)

build: $(BINARY)

$(BINARY): $(SOURCES)
	GOOS=linux GO111MODULE=on CGO_ENABLED=0 go build -o $(BINARY) cmd/seed/main.go

test: $(SOURCES)
	GO111MODULE=on go test -v -short -race -timeout 30s -coverprofile=coverage.txt -covermode=atomic ./...

docker-build: build
	docker build --no-cache -t $(IMAGE) .
