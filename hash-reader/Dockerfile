# hash-reader-app/Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o hash-reader .

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/hash-reader .

ENV PORT=5000

CMD ["./hash-reader"]
