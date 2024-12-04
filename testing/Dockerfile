# Use Go base image
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum main.go bot_test.go ./
RUN go mod tidy

RUN  go build -o test-server .

# Copy the token file into the container
COPY ../token.env /app/token.env

EXPOSE 8080

# Command to run tests
CMD ["./test-server"]
