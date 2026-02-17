# Magic Movie Stream

A movie streaming backend service built with Go and Gin framework.

## Tech Stack

- **Backend:** Go (Golang)
- **Framework:** Gin
- **Database:** MongoDB
- **Authentication:** JWT (JSON Web Tokens)
- **Containerization:** Docker & Docker Compose

## Project Structure

```
server/
├── cmd/servers/          # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database connection
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # JWT auth middleware
│   ├── models/           # Data models
│   ├── routers/          # Route definitions
│   └── servers/          # Server setup
├── Dockerfile
└── docker-compose.yaml
```

## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- MongoDB

### Environment Variables

Create a `.env` file in the `server/` directory:

```env
PORT=8080
DB_CONNECTION_URL=mongodb://mongodb:27017
DB_USERNAME=your_username
DB_PASSWORD=your_password
JWT_SECRET_KEY=your_secret_key
```

### Running with Docker

```bash
cd server
docker compose up -d --build
```

### Running Locally

```bash
cd server
go run cmd/servers/main.go
```

## API Endpoints

| Method | Endpoint         | Description          |
|--------|------------------|----------------------|
| GET    | `/hello-dear`    | Health check         |
| GET    | `/user/register` | User registration    |

## License

MIT
