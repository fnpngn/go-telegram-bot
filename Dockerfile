# Use a minimal base image for Go
FROM golang:1.20-alpine

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application code
COPY . .

# Build the application
RUN go build -o app ./cmd/main.go

# Expose the port (if needed)
EXPOSE 8080

# Run the application
CMD ["./app"]
