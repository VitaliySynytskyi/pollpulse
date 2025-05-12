FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service ./services/user-service

# Use a small alpine image
FROM alpine:3.17

# Add ca-certificates for HTTPS and tzdata for timezones
RUN apk --no-cache add ca-certificates tzdata

# Set up non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the binary from builder
COPY --from=builder /app/user-service /usr/local/bin/

# Copy migrations
COPY --from=builder /app/services/user-service/migrations /migrations

# Set the user to non-root
USER appuser

# Expose the port
EXPOSE 8081

# Set the entrypoint
ENTRYPOINT ["user-service"] 