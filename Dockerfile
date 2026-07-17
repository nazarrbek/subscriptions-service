FROM golang:1.26.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/subscriptions-service ./cmd/server

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/subscriptions-service /app/subscriptions-service

EXPOSE 8080

ENTRYPOINT ["/app/subscriptions-service"]
