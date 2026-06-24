FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o user-service cmd/user-service/main.go
RUN go build -o migrate cmd/migrator/main.go

FROM alpine:latest AS runtime

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=builder /app/user-service .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations

FROM runtime AS user-service
EXPOSE 8002
ENTRYPOINT [ "/app/user-service" ]

FROM runtime AS user-migrate
ENTRYPOINT [ "/app/migrate" ]
