.PHONY: all gen vet build

all: build

gen:
	go generate .

vet:
	go vet .
	golint .

build: vet
	go build .
