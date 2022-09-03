.PHONY: all binaries compress image test test-install test-serve vendor

commands = keep keepd

binaries = $(addprefix dist/, $(commands))
sources  = $(shell find . \( -name '*.go' -o -name '*.yml' \))

all: binaries

binaries: $(binaries)

compress: $(binaries)
	upx-ucl -1 $^

image:
	docker build -t epiphytelabs/keep:latest .

test: test-install test-serve

test-install: binaries
	docker stop keep-firefly-postgres || true
	docker rm keep-firefly-postgres || true
	docker stop keep-firefly || true
	docker rm keep-firefly || true
	docker network rm keep-firefly || true
	dist/keep install firefly

test-serve: binaries image
	@docker run \
		-it -p 40001:80 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		--net keep \
		epiphytelabs/keep:latest

vendor:
	go mod tidy
	go mod vendor
	go run vendor/github.com/goware/modvendor/main.go -copy="**/*.c **/*.h"

$(binaries): dist/%: $(sources)
	go build -o $@ -mod=vendor --ldflags="-s -w" ./cmd/$*