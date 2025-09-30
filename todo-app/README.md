# Todo App

A simple Todo application written in Go with configurable port settings through environment variables.

## Prerequisites

- Go 1.24.5 or higher
- Docker (optional)

## Project Structure

```
.
├── Dockerfile
├── go.mod
├── main.go
├── .env
└── README.md
```

## Environment Variables

Create a `.env` file in the project root:

```plaintext
PORT=8000
```

## Running Locally

1. Clone the repository:
```bash
git clone https://github.com:Stacvirus/hash-generator-app/tree/1.2
cd todo-app
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run main.go
```

The server will start on the port specified in your `.env` file (default: 8000)

## Running with Docker

1. Build the Docker image:
```bash
docker build -t todo-app .
```

2. Run the container:
```bash
# Using default port (8000)
docker run -p 8000:8000 todo-app

# Override port using environment variable
docker run -p 3000:3000 -e PORT=3000 todo-app

# Using custom .env file
docker run -p 8000:8000 -v $(pwd)/.env:/app/.env todo-app
```

## API Endpoints

- `GET /`: Returns a simple message indicating the app is running

## Development

### Modifying the Port

You can change the port in three ways:

1. Update the `.env` file
2. Set the PORT environment variable:
```bash
PORT=3000 go run main.go
```
3. Use Docker environment variables as shown above

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port number | 8000 |