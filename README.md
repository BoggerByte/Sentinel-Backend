# Sentinel-backend

API server of **[DMS] Sentinel** project written on Golang.  
Uses gRPC to communicate with [Sentinel-discord-bot](https://github.com/BoggerByte/Sentinel-Discord-Bot) app.

## Build

Create and run the entire backend server using docker-compose.yaml

```
docker compose up --remove-orphans
```

### Create and run by parts

Building server image

```
docker build -t sentinel-backend:latest .
```

Don't forget to set up docker network for all containers

```
// start postgres container
docker run -d \
    -e POSTGRES_USER='root' \
    -e POSTGRES_PASSWORD='qwerty' \
    -e POSTGRES_DB='sentinel_db' \
    -p 5432:5432 \
    --rm \
    --name=sentinel-db \
    --network=sentinel-network \
    postgres:14.3
    
// migrate db
make migrateup
    
// start redis container
docker run -d \
    -p 6379:6379 \
    --rm \
    --name=sentinel-redis \
    --network=sentinel-network \
    redis /bin/sh -c 'redis-server --appendonly yes --requirepass qwerty'

// start server container
docker run \
    -e DB_HOST=$db_host \
    -e REDIS_HOST=$redis_host \
    -e GIN_MODE='release' \
    -p 8080:8080 \
    --name=sentinel-backend \
    --network=sentinel-network \
    sentinel-backend:latest
```

## Dev

Running dev server.  
It is necessary to install tools like [golang-migrate](https://github.com/golang-migrate/migrate),
[sqlc](https://github.com/kyleconroy/sqlc), [gomock](https://github.com/golang/mock).

```
// generate db queries
make sqlc

// start postgres container
make postgres
make migrateup

// start redis container
make redis

// start server
make server
```

Testing

```
// generate db and memory db mocks
make mock

// run tests
make test
```