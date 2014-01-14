.PHONY: bin test all fmt

all: fmt test bin

bin:
	(cd main; go build materials.go)
	(cd mcfs/main; go build mcfs.go)

test:
	-./runtests.sh

fmt:
	-go fmt ./...
