FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o manager

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/manager .
ENTRYPOINT ["/manager"]