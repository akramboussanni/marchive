# Multi-stage build for Go backend and Svelte frontend
FROM node:18-alpine AS frontend-builder

# Set working directory for frontend
WORKDIR /app/frontend

# Copy frontend package files
COPY frontend/package*.json ./
COPY frontend/pnpm-lock.yaml ./

# Install pnpm and frontend dependencies
RUN npm install -g pnpm
RUN pnpm install

# Copy frontend source code
COPY frontend/ ./

# Go backend build stage
FROM golang:1.24-alpine AS backend-builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory for backend
WORKDIR /app/backend

# Copy Go module files
COPY backend/go.mod backend/go.sum ./

# Download Go dependencies
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final runtime stage
FROM node:18-alpine

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Install pnpm globally
RUN npm install -g pnpm

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy Go binary from backend stage
COPY --from=backend-builder /app/backend/main ./backend/

# Copy frontend from frontend stage
COPY --from=frontend-builder /app/frontend ./frontend/

# Create necessary directories
RUN mkdir -p /app/data && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose ports for both backend and frontend
EXPOSE 9520 5173

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9520/ || exit 1

# Create a startup script to run both services
RUN echo '#!/bin/sh\n\
cd /app/frontend && pnpm start &\n\
cd /app/backend && ./main\n\
wait' > /app/start.sh && chmod +x /app/start.sh

# Run the startup script
CMD ["/app/start.sh"]
