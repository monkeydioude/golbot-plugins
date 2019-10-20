.PHONY: all build test

all: build test

test:
	@AGENT_FILE=/tmp/golbot_test.agent go test .

build:
	@./scripts/build.sh
