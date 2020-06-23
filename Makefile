.PHONY: build

build:
	go build -v ./cmd/monitoring

.DEFAULT_GOAL := build