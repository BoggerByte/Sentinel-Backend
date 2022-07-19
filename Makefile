postgres:
	docker run -d \
		-e POSTGRES_USER='root' \
		-e POSTGRES_PASSWORD='qwerty' \
		-p 5432:5432 \
		--rm \
		--name=sentinel-db \
		postgres:14.3

redis:
	docker run -d \
		-p 6379:6379 \
		--rm \
		--name=sentinel-redis \
		redis /bin/sh -c 'redis-server --appendonly yes --requirepass qwerty'

createdb:
	docker exec -it sentinel-db createdb --username=root --owner=root sentinel_db

dropdb:
	docker exec -it sentinel-db dropdb --username=root sentinel_db

migrateup:
	migrate -path ./pkg/db/migration -database "postgresql://root:qwerty@localhost:5432/sentinel_db?sslmode=disable" up

migratedown:
	migrate -path ./pkg/db/migration -database "postgresql://root:qwerty@localhost:5432/sentinel_db?sslmode=disable" down

sqlc:
	sqlc generate -f ./cfg/sqlc.yaml

mock:
	mockgen -package mockdb -destination pkg/db/mock/store.go github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc Store
	mockgen -package mockmemdb -destination pkg/db/memory_mock/store.go github.com/BoggerByte/Sentinel-backend.git/pkg/db/memory Store

server:
	go run cmd/main.go

.PHONY: postgres redis createdb dropdb migrateup migratedown sqlc server
