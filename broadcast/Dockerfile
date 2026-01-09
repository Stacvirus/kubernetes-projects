# Use official Golang image as builder
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o broadcast .

# Use minimal alpine image for final container
FROM alpine:3.18

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/broadcast .
COPY --from=builder /app/.env .

# Set default port
ENV PORT=8585

# Run the application
CMD ["./broadcast"]