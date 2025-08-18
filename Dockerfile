# Build stage
FROM golang:1.21-alpine AS builder

# Install git for downloading dependencies
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN mkdir -p build && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/l8k .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set the working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/build/l8k .

# Change ownership to appuser
RUN chown appuser:appuser /app/l8k

# Switch to non-root user
USER appuser

# Expose port (if needed for future web interface)
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./l8k"]
