run: build
	@./bin/api

build:
	@go build -o bin/api cmd/api/main.go

test:
	@go test -v ./...

clean:
	rm -rf ./bin

postgresup:
	docker-compose up -d

postgresdown:
	docker-compose down -v

migrateup:
	migrate -path db/migration -database "postgresql://postgres:example@127.0.0.1:5432/mingle_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:example@127.0.0.1:5432/mingle_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: run build test clean mysql-up mysql-down migrateup migratedown sqlc
