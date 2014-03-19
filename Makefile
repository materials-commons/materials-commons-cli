.PHONY: bin test all fmt docs

all: fmt test bin

bin:
	(cd main; go build materials.go)

docs:
	./makedocs.sh

test:
	-./runtests.sh

fmt:
	-go fmt ./...
