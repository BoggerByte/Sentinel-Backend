postgres:
	docker run --name=sentinel-db -e POSTGRES_USER='root' -e POSTGRES_PASSWORD='qwerty' -p 0.0.0.0:5432:5432 -d --rm postgres:14.3

createdb:
	docker exec -it sentinel-db createdb --username=root --owner=root sentinel_db

dropdb:
	docker exec -it sentinel-db dropdb --username=root sentinel_db

migrateup:
	migrate -path ./migrations -database "postgresql://root:qwerty@localhost:5432/sentinel_db?sslmode=disable" -verbose up

migratedown:
	migrate -path ./migrations -database "postgresql://root:qwerty@localhost:5432/sentinel_db?sslmode=disable" -verbose down

migratedrop:
	migrate -path ./migrations -database "postgresql://root:qwerty@localhost:5432/sentinel_db?sslmode=disable" -verbose drop

.PHONY: postgres createdb dropdb migrateup migratedown
