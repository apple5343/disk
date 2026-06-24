FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o file-storage cmd/file-storage/main.go

FROM alpine:latest AS runtime

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=builder /app/file-storage .

FROM runtime AS file-storage
ENTRYPOINT [ "/app/file-storage" ]