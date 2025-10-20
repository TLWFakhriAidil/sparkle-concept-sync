# Sparkle Concept Sync

An enterprise-grade WhatsApp AI chatbot platform with visual flow builder, multi-provider support, and scalability for 5000+ concurrent users.

## Features

### Core Platform
- **Visual Flow Builder**: Drag-and-drop interface with 10 different node types
- **Multi-Provider WhatsApp Integration**: Support for WAHA, Wablas, and Whacenter
- **AI Integration**: OpenRouter API with multiple AI models (GPT-5, Gemini 2.5 Pro, etc.)
- **Real-time Updates**: WebSocket-based live monitoring and updates
- **Scalable Architecture**: Designed for 5000+ concurrent users

### Node Types
- **Start Node**: Flow entry point
- **Message Node**: Send text messages
- **AI Prompt Node**: AI-powered responses with model selection
- **Condition Node**: Conditional logic branching
- **Delay Node**: Add time delays in flows
- **Image/Audio/Video Nodes**: Media message support
- **User Reply Node**: Wait for user input
- **Stage Node**: Set conversation stages

### Backend Features
- **Go Fiber Framework**: High-performance HTTP server
- **PostgreSQL Database**: Robust data persistence with optimized schemas
- **Redis Caching**: Response caching and session management
- **JWT Authentication**: Secure user authentication
- **Rate Limiting**: Per-user and global rate limiting
- **Circuit Breakers**: Fault tolerance for external services

### Frontend Features
- **React + TypeScript**: Modern web development stack
- **@xyflow/react**: Professional flow builder interface
- **Tailwind CSS**: Responsive and modern UI design
- **Real-time Dashboard**: Live conversation monitoring
- **Device Management**: Multi-device WhatsApp integration

## Quick Start

### Prerequisites
- Go 1.23+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd sparkle-concept-sync
```

2. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Install dependencies**
```bash
# Frontend dependencies
npm install

# Go dependencies
go mod download
```

4. **Set up database**
```bash
# Create PostgreSQL database
createdb sparkle_db

# Run migrations (automatic on first start)
```

5. **Start development servers**
```bash
# Backend (Go)
go run cmd/server/main.go

# Frontend (React)
npm run dev
```

### Production Deployment

#### Railway Platform
1. **Configure Railway project**
```bash
# Install Railway CLI
npm install -g @railway/cli

# Login and create project
railway login
railway init
```

2. **Set environment variables in Railway dashboard**
- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string  
- `JWT_SECRET`: Secure JWT secret
- `OPENROUTER_API_KEY`: Your OpenRouter API key

3. **Deploy**
```bash
railway up
```

#### Docker Deployment
```bash
# Build and run with Docker
docker build -t sparkle-concept-sync .
docker run -p 8080:8080 sparkle-concept-sync
```

## API Documentation

### Authentication Endpoints
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh JWT token

### Device Management
- `GET /api/devices` - List user devices
- `POST /api/devices` - Create new device
- `PUT /api/devices/:id` - Update device
- `DELETE /api/devices/:id` - Delete device

### WhatsApp Webhooks
- `POST /webhooks/waha/:device_id` - WAHA webhook endpoint
- `POST /webhooks/wablas/:device_id` - Wablas webhook endpoint  
- `POST /webhooks/whacenter/:device_id` - Whacenter webhook endpoint

### Health & Monitoring
- `GET /healthz` - Basic health check
- `GET /health` - Detailed health status
- `GET /api/webhook-info` - Webhook configuration info

### WebSocket
- `WS /ws` - Real-time updates and monitoring

## Configuration

### WhatsApp Provider Setup

#### WAHA (WhatsApp HTTP API)
1. Set up WAHA instance
2. Configure webhook URL: `https://your-domain.com/webhooks/waha/{device_id}`
3. Add device in platform with provider = "waha"

#### Wablas
1. Get Wablas API key
2. Set webhook URL: `https://your-domain.com/webhooks/wablas/{device_id}`
3. Configure device with provider = "wablas"

#### Whacenter  
1. Get Whacenter API credentials
2. Set webhook URL: `https://your-domain.com/webhooks/whacenter/{device_id}`
3. Configure device with provider = "whacenter"

### AI Model Configuration
Supported models via OpenRouter:
- `openai/gpt-5-chat` - GPT-5 Chat
- `openai/gpt-5-mini` - GPT-5 Mini  
- `openai/chatgpt-4o-latest` - GPT-4o Latest
- `google/gemini-2.5-pro` - Gemini 2.5 Pro
- `google/gemini-pro-1.5` - Gemini Pro 1.5

## Architecture

### Backend (Go)
```
cmd/server/          - Application entry point
internal/config/     - Configuration management
internal/database/   - Database connection & migrations
internal/models/     - Data models
internal/services/   - Business logic services
internal/handlers/   - HTTP handlers
internal/middleware/ - HTTP middleware
```

### Frontend (React)
```
src/components/      - React components
src/components/nodes/ - Flow builder node components
src/pages/          - Page components
src/hooks/          - Custom React hooks
src/types/          - TypeScript type definitions
```

### Database Schema
- `users` - User accounts and profiles
- `device_setting` - WhatsApp device configurations
- `conversation` - Chat conversations
- `message` - Individual messages
- `conversation_flow` - Chatbot flows
- `user_session` - Active user sessions

## Performance

### Scalability Features
- **Connection Pooling**: Optimized database connections
- **Redis Caching**: Response and session caching  
- **Rate Limiting**: Per-user and global limits
- **Circuit Breakers**: Fault tolerance
- **Horizontal Scaling**: Stateless application design

### Monitoring
- Health check endpoints
- Real-time WebSocket monitoring
- Performance metrics collection
- Error tracking and logging

## Development

### Adding New Node Types
1. Create node component in `src/components/nodes/`
2. Add node to flow builder configuration
3. Implement backend flow execution logic
4. Update database schema if needed

### Adding New WhatsApp Providers
1. Create provider-specific webhook handler
2. Implement message format conversion
3. Add provider configuration options
4. Update API documentation

## Support

For technical support and questions:
- Review the API documentation
- Check the troubleshooting guide
- Submit issues via the repository

## License

This project is proprietary software. All rights reserved.