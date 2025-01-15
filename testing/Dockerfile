# Use an official Golang image as the base
FROM golang:1.23.4

# Set the working directory
WORKDIR /app

# Copy the application source code
COPY . .

# Install application dependencies
RUN go mod tidy

# Build the application
RUN go build -o server .

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./server"]
