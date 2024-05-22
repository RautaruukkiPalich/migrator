run:
	go run cmd/app/main.go --config=./config/config.yaml

docker:
	docker-compose up -d

makemigrations:
	migrate create -ext sql -dir migrations $(name)

migrate:
	go run ./cmd/migrate --config=./config/config.yaml

lint:
	golangci-lint run ./...

filldb:
	go run ./cmd/filldb --config=./config/config.yaml
