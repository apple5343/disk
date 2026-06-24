FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o file-data cmd/file-data/main.go
RUN go build -o migrate cmd/migrator/main.go

FROM alpine:latest AS runtime

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=builder /app/file-data .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations

FROM runtime AS file-data
ENTRYPOINT [ "/app/file-data" ]

FROM runtime AS migrate
ENTRYPOINT [ "/app/migrate" ]