postgres:
	docker run --name=sentinel-db -e POSTGRES_USER='root' -e POSTGRES_PASSWORD='qwerty' -p 0.0.0.0:5432:5432 -d --rm postgres:14.3

createdb:
	docker exec -it sentinel-db createdb --username=root --owner=root sentinel_db

dropdb:
	docker exec -it sentinel-db dropdb --username=root sentinel_db

migrateup:
	migrate -path ./db/migration -database "postgresql://root:qwerty@localhost:5432/sentinel_db?sslmode=disable" up

migratedown:
	migrate -path ./db/migration -database "postgresql://root:qwerty@localhost:5432/sentinel_db?sslmode=disable" down

sqlc:
	sqlc generate

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server
