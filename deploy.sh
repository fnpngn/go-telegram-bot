#!/bin/bash

# Function to build Docker images
build_images() {
    echo "Building Docker images..."
    docker build -t telegram-bot .
    docker build -t test-results-server .
}

# Function to run Docker containers
run_containers() {
    echo "Running Docker containers..."

    # Read the Telegram bot token from the file
    TELEGRAM_BOT_TOKEN=$(cat token.txt)

    # Run the Telegram bot container
    docker run --rm -it -e TOKEN=$TELEGRAM_BOT_TOKEN telegram-bot &

    # Run the test results server container
    docker run -p 8080:8080 test-results-server &
}

# Main script execution
main() {
    echo "Starting deployment process..."

    # Build and run the Docker containers
    build_images
    run_containers

    echo "Deployment complete."
}

# Execute the main function
main
