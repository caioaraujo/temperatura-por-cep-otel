FROM golang:1.22.2-alpine AS builder
ENV CGO_ENABLED=1

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod verify && \
    go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o bin/cep-validator cmd/service_a/server/main.go
RUN CGO_ENABLED=0 go build -o bin/temperatura-cep cmd/service_b/server/main.go

FROM alpine:3.19
COPY --from=builder /app/bin/cep-validator /usr/bin/cep-validator
COPY --from=builder /app/bin/temperatura-cep /usr/bin/temperatura-cep