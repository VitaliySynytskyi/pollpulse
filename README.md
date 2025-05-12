# PollPulse - Survey Platform

PollPulse is a microservice-based survey platform built with Go and Vue.js. It provides a complete solution for creating, managing, and analyzing surveys.

## Architecture

The platform consists of the following microservices:

- **User Service**: Handles user authentication and management
- **Survey Service**: Manages survey creation and responses
- **Result Service**: Processes and aggregates survey results
- **API Gateway**: Provides a unified API interface
- **Frontend**: Vue.js application for the user interface

## Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for local frontend development)
- Go 1.21+ (for local backend development)

## Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/pollpulse.git
   cd pollpulse
   ```

2. Start the application using Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. Access the application:
   - Frontend: http://localhost
   - API Gateway: http://localhost:8080
   - Individual services:
     - User Service: http://localhost:8081
     - Survey Service: http://localhost:8082
     - Result Service: http://localhost:8083

## Development Setup

### Backend Development

1. Install Go dependencies:
   ```bash
   go mod download
   ```

2. Run individual services:
   ```bash
   # User Service
   go run services/user/main.go

   # Survey Service
   go run services/survey/main.go

   # Result Service
   go run services/result/main.go

   # API Gateway
   go run services/gateway/main.go
   ```

### Frontend Development

1. Install dependencies:
   ```bash
   cd frontend
   npm install
   ```

2. Start development server:
   ```bash
   npm run dev
   ```

## Testing

### Backend Tests

Run all backend tests:
```bash
go test ./...
```

### Frontend Tests

Run frontend tests:
```bash
cd frontend
npm run test
```

## Deployment

### Production Deployment

1. Set up environment variables:
   - Create a `.env` file with production values
   - Update secrets and credentials

2. Build and deploy:
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### CI/CD

The project includes GitHub Actions workflows for:
- Automated testing
- Docker image building
- Deployment to production

## API Documentation

API documentation is available at:
- Swagger UI: http://localhost:8080/swagger-ui
- OpenAPI Spec: http://localhost:8080/swagger.json

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 