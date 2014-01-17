.PHONY: bin test all fmt

all: fmt test bin

bin:
	(cd main; go build materials.go)

test:
	-./runtests.sh

fmt:
	-go fmt ./...
