.PHONY: all

gen:
	go generate

build:
	go build .

all: gen build