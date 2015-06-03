image="sinar/sherpa"
version:=$(shell date +%Y.%m.%d)
all:
	go clean
	go build
	strip sherpa
	docker build -t ${image} .
	go clean
