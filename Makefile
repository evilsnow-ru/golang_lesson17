.PHONY: all gen vet build

all: build

gen:
	protoc -I/usr/local/include/ -I./api/ --go_out=./api/ messages.proto

vet:
	go vet .
	golint .

fmtchk:
	gofmt -l .

fmt:
	gofmt -w .

build: vet fmtchk
	go build .
