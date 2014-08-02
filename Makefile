GOFMT=gofmt -s
GOFILES=$(wildcard *.go **/*.go)
PRERELEASE=alpha

default: build

format:
	${GOFMT} -w ${GOFILES}

run:
	go run cmd/goverflow/main.go

build:
	mkdir -p build
	cd cmd/goverflow && goxc -c=.goxc.json -pr="$(PRERELEASE)" -d ../../build

.PHONY: build
