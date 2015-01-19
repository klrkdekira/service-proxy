all:
	go clean
	go build
	strip go-api-mirror
	docker build -t sinar/service-proxy .
	go clean
