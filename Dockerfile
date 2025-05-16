FROM golang:alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOSUMDB=off go install github.com/pressly/goose/cmd/goose@v2.7.0

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /go/bin/load_balancer \
    cmd/load_balancer/main.go

FROM alpine:3.17

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /go/bin/goose   /usr/local/bin/goose
COPY --from=builder /go/bin/load_balancer /usr/local/bin/load_balancer

COPY --from=builder /app/migrations     /app/migrations
COPY --from=builder /app/config/config.yml /app/config/config.yml

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "-c", "\
    goose -dir /app/migrations \
      postgres \"$DB_DSN\" up && \
    /usr/local/bin/load_balancer \
"]
