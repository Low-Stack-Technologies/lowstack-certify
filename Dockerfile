FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/certify ./cmd/main

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/certify /app/certify

ENTRYPOINT ["/app/certify"]