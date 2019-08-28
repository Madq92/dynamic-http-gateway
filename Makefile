.PHONY: default build dependencies image

default: build

build:
	CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o dynamic-http-gateway

dependencies:
	go mod download

image:
	docker build -t dynamic-http-gateway .
