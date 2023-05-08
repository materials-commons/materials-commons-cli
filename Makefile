.PHONY: bin all fmt deploy cli

all: fmt bin

fmt:
	-go fmt ./...

bin: cli

cli:
	(cd ./cmd/mcc; go build)

run: cli
	./cmd/mcc/mcc

deploy:  cli
	sudo cp cmd/mcc/mcc /usr/local/bin
	sudo chmod a+rx /usr/local/bin/mcc
