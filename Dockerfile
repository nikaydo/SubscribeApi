FROM golang:1.24.6-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o subs ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/subs /app/subs

COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY --from=builder /build/migrations /migrations

EXPOSE 8080