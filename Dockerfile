FROM golang:1.21.1-alpine AS builder

COPY . /app/
WORKDIR /app

RUN go mod download
RUN go build -o /bin/uzum_delivery_service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /bin/uzum_delivery_service .
COPY --from=builder /app/dev/local.env dev/local.env

CMD ["./uzum_delivery_service"]
