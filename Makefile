run:
	go run cmd/app/main.go --config=./config/config.yaml

docker:
	docker-compose up -d

makemigrations:
	migrate create -ext sql -dir migrations $(name)

migrate:
	go run ./cmd/migrator --config=./config/config.yaml

test: migrate run
