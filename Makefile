all:
	go clean
	go build
	strip service-proxy
	docker build -t sinar/service-proxy .
	go clean
