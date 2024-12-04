# Use an official Golang runtime as a parent image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o bot .

# Expose the default port (optional, depends on how you intend to access logs or metrics)
EXPOSE 8080

# Command to run the bot
CMD ["./bot"]
