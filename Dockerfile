# Multi-stage Docker build for production deployment

# Stage 1: Frontend build
FROM node:18-alpine AS frontend-builder
WORKDIR /app

# Copy package files
COPY package*.json ./
COPY bun.lockb* ./

# Install dependencies with npm (Railway compatibility)
RUN npm ci --only=production

# Copy source code
COPY . .

# Build the frontend
RUN npm run build

# Stage 2: Backend build
FROM golang:1.23-alpine AS backend-builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set environment for Go build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /src

# Copy Go modules files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN go build -a -installsuffix cgo -o /app/server ./cmd/server

# Stage 3: Production image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=UTC

# Create non-root user
RUN addgroup -g 1001 -S sparkle && \
    adduser -S sparkle -u 1001

# Create necessary directories
RUN mkdir -p /app/dist /app/uploads && \
    chown -R sparkle:sparkle /app

# Copy built files
COPY --from=backend-builder /app/server /app/server
COPY --from=frontend-builder /app/dist /app/dist

# Set ownership
RUN chown -R sparkle:sparkle /app

# Switch to non-root user
USER sparkle

WORKDIR /app

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# Run the server
CMD ["/app/server"]