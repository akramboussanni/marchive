# marchive

Book library with over 100M+ books using Anna's Archive.

## Features

- **Book Discovery**: Automated book scraping and metadata extraction
- **Download Management**: Queue-based download system with progress tracking
- **Admin Panel**: Comprehensive administration tools for system management
- **Real-time Updates**: Live progress tracking and status updates

## Screenshots

![marchived Dashboard](static/dashboard.png)

*The main dashboard*

## 🐳 Docker Deployment

### Quick Start with Docker Compose
Use the `docker-compose.yml` file.

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `JWT_SECRET` | Secret key for JWT token signing | ✅ | - |
| `ANNAS_API_KEY` | API key for Anna book service | ✅ | - |
| `DOMAIN` | Domain for cookies and CORS (e.g., localhost, yourdomain.com) | ❌ | `localhost` |
| `APP_PORT` | Backend server port | ❌ | `9520` |
| `DB_CONNECTION_STRING` | PostgreSQL connection string | ❌ | `postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable` |
| `LOGGER_TYPE` | Logging system type | ❌ | `std` |
| `TRUST_PROXY_IP_HEADERS` | Trust proxy IP headers | ❌ | `false` |

**Note**: The `DOMAIN` variable is used for both cookie domain and CORS origin configuration. Set this to your production domain when deploying (ex: example.com)

## 🛠️ Development

### Prerequisites

- Go 1.21+
- Node.js 18+
- pnpm
- you don't need postgres, if your go is ran with `-tags=debug` it will run as sqlite db. (ex: `go run -tags=debug cmd/server/main.go`)

### Local Development Setup

#### 1. Build the Frontend

```bash
cd frontend
pnpm install
pnpm build
```

This will create a `build/` directory with the compiled frontend assets.

#### 2. Run the Backend

```bash
# From the root directory
go run -tags=debug cmd/server/main.go
```

The `-tags=debug` flag enables SQLite mode for local development, so you don't need PostgreSQL running.

#### 3. Access the Application

Open your browser and go to `http://localhost:9520` - the backend will serve both the API and the frontend.

#### Hot Reload Development

For development with hot reload (frontend on localhost:5173, backend on localhost:9520):

```bash
python dev.py
```

This will:
- Start the frontend dev server on http://localhost:5173 with hot reload
- Start the backend server on http://localhost:9520
- Automatically install dependencies
- Provide real-time logs for both services

## 📁 Project Structure

```
marchive/
├── cmd/                     # Application entry points
│   └── server/             # Main server binary
├── internal/                # Private application code
│   ├── anna/               # Book scraping service
│   ├── api/                # HTTP handlers and routes
│   ├── applog/             # Logging configuration
│   ├── config/             # Configuration management
│   ├── db/                 # Database layer
│   ├── jwt/                # JWT authentication
│   ├── middleware/         # HTTP middleware
│   ├── model/              # Data models
│   ├── repo/               # Repository layer
│   ├── services/           # Business logic services
│   └── utils/              # Utility functions
├── frontend/                # Svelte frontend application
│   ├── src/                # Source code
│   │   ├── components/     # Reusable UI components
│   │   ├── routes/         # Page components
│   │   ├── stores/         # State management
│   │   └── utils/          # Frontend utilities
│   ├── build/              # Built frontend assets (generated)
│   ├── package.json        # Frontend dependencies
│   └── pnpm-lock.yaml      # Locked dependency versions
├── config/                  # Configuration files
├── downloads/               # Downloaded book storage
├── Dockerfile              # Multi-stage Docker build
├── docker-compose.yml      # Development environment
├── docker-compose-https.yml # HTTPS development environment
├── go.mod                  # Go module definition
├── go.sum                  # Go dependency checksums
└── README.md               # This file
```