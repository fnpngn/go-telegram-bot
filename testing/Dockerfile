# Use Go base image
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum testing/main.go testing/bot_test.go ./
RUN go mod tidy

RUN  go build -o test-server .

EXPOSE 8080

# Command to run tests
CMD ["./test-server"]
