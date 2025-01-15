# Stage 1: Build Stage
FROM golang:1.23.4-alpine AS builder

# Set environment variables for building
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy dependency files first for efficient caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o app ./cmd/main.go

# Stage 2: Runtime Stage
FROM alpine:3.18

WORKDIR /app

# Copy only the built binary from the build stage
COPY --from=builder /app/app .

# Set a non-root user for better security
RUN adduser -D -u 1001 appuser
USER appuser

# Command to run the application
CMD ["./app"]
