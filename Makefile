image="sinar/sherpa"
version:=$(shell date +%Y.%m.%d)
all: build
build:
	go clean
	go build
	strip sherpa
	docker build -t ${image} .
	docker tag -f ${image}:latest ${image}:$(version)
	go clean
push:
	docker push ${image}
