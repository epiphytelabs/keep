.PHONY: all binaries compress image test test-install test-serve vendor

commands = keep keepd

binaries = $(addprefix dist/, $(commands))
sources  = $(shell find . \( -name '*.go' -o -name '*.yml' \))

all: binaries

binaries: $(binaries)

compress: $(binaries)
	upx-ucl -1 $^

image:
	docker build -t epiphytelabs/keep:dev .

test: test-install test-serve

test-install: binaries
	docker stop keep-firefly-postgres || true
	docker rm keep-firefly-postgres || true
	docker stop keep-firefly || true
	docker rm keep-firefly || true
	docker network rm keep-firefly || true
	dist/keep install firefly

test-serve: binaries image
	dist/keep server uninstall || true
	dist/keep server install
	docker logs -f keep

vendor:
	go mod tidy
	go mod vendor
	go run vendor/github.com/goware/modvendor/main.go -copy="**/*.c **/*.h"

$(binaries): dist/%: $(sources)
	go build -o $@ -mod=vendor --ldflags="-s -w" ./cmd/$*