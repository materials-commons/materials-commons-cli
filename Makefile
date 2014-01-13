.PHONY: bin test all fmt

all: fmt test bin

bin:
	(cd materials; go build materials.go)
	(cd mcfs; go build mcfs.go)

test:
	-./runtests.sh

fmt:
	-go fmt ./...
