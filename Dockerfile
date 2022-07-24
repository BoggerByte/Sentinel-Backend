# Build stage
FROM golang:1.18-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

# Run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
COPY --from=builder /app/cfg ./cfg
COPY --from=builder /app/pub/html ./pub/html

EXPOSE 8080
CMD [ "/app/main" ]
