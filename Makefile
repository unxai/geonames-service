.PHONY: build run test clean

BINARY_NAME=geonames-service

build:
	go build -o bin/$(BINARY_NAME) main.go

run:
	go run main.go

test:
	go test ./...

clean:
	rm -rf bin/
	docker-compose down

dev: clean
	docker-compose up -d
	go run main.go