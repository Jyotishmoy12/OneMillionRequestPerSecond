FROM golang:1.27-alpine AS builder

WORKDIR /app


COPY go.mod  go.sum ./

RUN go mod download

COPY cmd ./cmd

COPY internal ./internal

RUN go test ./...

RUN go build -o api ./cmd/api

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/api ./api

EXPOSE 8080


CMD ["./api"]