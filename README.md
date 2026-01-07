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
Use the `docker-compose.yml` file, set env vars and deploy. Should work instantly

!!! Currently not working because of frontend rewrite, will be fixed soon !!!

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
- Node.js 22+
- you don't need postgres, if your go is ran with `-tags=debug` it will run as sqlite db. (ex: `go run -tags=debug cmd/server/main.go`)

### Local Development Setup

#### 1. Run the Frontend

```bash
cd frontend
npm install
npm run dev
```

#### 2. Run the Backend

```bash
# From the root directory, so maybe do cd ../
go run -tags=debug cmd/server/main.go
```

The `-tags=debug` flag enables SQLite mode for local development, so you don't need PostgreSQL running.

#### 3. Access the Application

Open your browser and go to `http://localhost:5173`.
