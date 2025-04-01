# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o url-shortener .

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/url-shortener .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./url-shortener"]