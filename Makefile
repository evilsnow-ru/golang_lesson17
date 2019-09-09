.PHONY: all

gen:
	go generate

all: gen
	go build .