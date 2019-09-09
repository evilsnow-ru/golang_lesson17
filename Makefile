.PHONY: vet test gen build

vet:
	go vet .

test:
	go test .

gen:
	protoc --go_out=. messages.proto

build: vet test
	go build -o distr/otuslesson17 .