# log output App

A Go application that generates SHA-1 hashes periodically and exposes them via an HTTP endpoint.

## Prerequisites

- Go 1.24.5 or higher
- Docker (optional)

## Running Locally

1. Clone the repository:
```bash
git clone https://github.com/stacvirus/kubernetes-projects.git
cd kubernetes-projects
```

2. Update Go module:
```bash
go mod tidy
```

3. Run the application:
```bash
go run main.go
```

The server will start on port 8080 and will:
- Log a new hash every 5 seconds
- Respond to HTTP requests at `http://localhost:8080`

## Running with Docker

1. Build the Docker image:
```bash
docker build -t log-output-go .
```

2. Run the container:
```bash
docker run -p 8080:8080 log-output-go
```

## Testing the Application

You can test the application using curl:
```bash
curl http://localhost:8080
```

Or open in your browser:
```
http://localhost:8080
```

## Project Structure

```
.
├── Dockerfile
├── go.mod
├── main.go
└── README.md
```

## Features

- Generates SHA-1 hash from "Hello, World!" string
- Logs new hash every 5 seconds
- Exposes hash generation via HTTP endpoint
- Containerized application support

## Exercises
### Chapter 2
* [1.1](https://github.com/Stacvirus/kubernetes-projects/tree/1.1)
