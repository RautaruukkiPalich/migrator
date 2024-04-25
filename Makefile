run:
	go run cmd/main.go

docker:
	docker-compose up -d

lint:
	golangci-lint run ./...