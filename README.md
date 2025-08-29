# marchive

Book library with over 100M+ books using Anna's Archive.

## Features

- **Book Discovery**: Automated book scraping and metadata extraction
- **Download Management**: Queue-based download system with progress tracking
- **Admin Panel**: Comprehensive administration tools for system management
- **Real-time Updates**: Live progress tracking and status updates

## Screenshots

![Marchived Dashboard](static/dashboard.png)

*The main dashboard*

## 🐳 Docker Deployment

### Quick Start with Docker Compose
Use the `docker-compose.yml` file.

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `JWT_SECRET` | Secret key for JWT token signing | ✅ | - |
| `ANNAS_API_KEY` | API key for Anna book service | ✅ | - |
| `APP_PORT` | Backend server port | ❌ | `9520` |
| `DB_CONNECTION_STRING` | PostgreSQL connection string | ❌ | `postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable` |
| `LOGGER_TYPE` | Logging system type | ❌ | `zap` |
| `FRONTEND_CORS` | Frontend CORS origin | ❌ | `*` |
| `COOKIE_DOMAIN` | Cookie domain for authentication | ❌ | `localhost` |
| `TRUST_PROXY_IP_HEADERS` | Trust proxy IP headers | ❌ | `false` |

## 🛠️ Development

### Prerequisites

- Go 1.21+
- Node.js 18+
- pnpm
- you don't need postgres, if ran with -tags=debug it will run as sqlite db.

### Backend Setup

```bash
cd backend
go mod download
go run -tags=debug cmd/server/main.go
```

### Frontend Setup

```bash
cd frontend
pnpm install
pnpm dev
```

## 📁 Project Structure

```
marchived/
├── backend/                 # Go backend API
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   │   ├── anna/          # Book scraping service
│   │   ├── api/           # HTTP handlers and routes
│   │   └── applog/        # Logging configuration
│   ├── config/            # Configuration management
│   └── docs/              # Documentation
├── frontend/               # Svelte frontend application
│   ├── src/               # Source code
│   │   ├── components/    # Reusable UI components
│   │   ├── routes/        # Page components
│   │   └── stores/        # State management
│   └── public/            # Static assets
├── Dockerfile             # Multi-stage Docker build
├── docker-compose.yml     # Development environment
└── README.md             # This file
```

## 🔧 Configuration

### Backend Configuration

The backend uses a configuration system that supports:
- Environment variables
- Configuration files
- Default values

Key configuration areas:
- **Database**: Connection settings and pooling
- **Authentication**: JWT configuration and security
- **Logging**: Log levels and output formats
- **API**: Rate limiting and CORS settings

### Frontend Configuration

The frontend is configured through:
- Environment variables for API endpoints
- Tailwind CSS for styling
- Svelte stores for state management
