.PHONY: all build test

all: build test

test:
	@go test . ./pkg/alert

build:
	@./scripts/build.sh
