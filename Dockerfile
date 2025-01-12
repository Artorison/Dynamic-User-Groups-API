FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/user_groups_api

RUN go build -o /app/bin/main .


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/bin/main .
COPY config/config.yml /app/config/config.yml

CMD ["/app/main"]