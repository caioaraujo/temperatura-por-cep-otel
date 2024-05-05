FROM golang:1.22.2-alpine AS builder
ENV CGO_ENABLED=1

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod verify && \
    go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o bin/service_A cmd/service_a/server/main.go
RUN CGO_ENABLED=0 go build -o bin/service_B cmd/service_b/server/main.go

FROM alpine:3.19
COPY --from=builder /app/bin/service_A /usr/bin/service_A
COPY --from=builder /app/bin/service_B /usr/bin/service_B