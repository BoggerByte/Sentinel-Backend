# Build stage
FROM golang:1.18-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

# Run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY /pub/html /app/pub/html
COPY app.env .
COPY .env .

EXPOSE 8080
CMD [ "/app/main" ]
