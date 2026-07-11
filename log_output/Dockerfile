# Use official Golang image as builder
FROM golang:latest AS builder

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o hash-generator .

# Use minimal alpine image for final container
FROM alpine:3.18

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/hash-generator .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./hash-generator"]