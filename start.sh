#!/bin/bash

echo "Starting Resume.in application..."

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Make the script executable
chmod +x ./start.sh

# Build and start the containers
echo "Building and starting Docker containers..."
docker-compose up --build

# Wait for user input to exit
echo "Press Ctrl+C to stop the application." 