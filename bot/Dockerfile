# Stage 1: Build Stage
FROM golang:1.23.4-alpine AS builder

# Set environment variables for building
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files first from the bot directory
COPY . .

# Download Go dependencies
RUN go mod download


# Build the Go application
RUN go build -o app ./main.go

# Stage 2: Runtime Stage
FROM alpine:3.18

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app .

# Set a non-root user for better security
RUN adduser -D -u 1001 appuser
USER appuser

# Command to run the application
CMD ["./app"]
