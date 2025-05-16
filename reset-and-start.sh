#!/bin/bash

echo "Stopping any running containers..."
docker-compose down

echo "Removing PostgreSQL volume to get a fresh start..."
docker volume rm $(docker volume ls -q | grep postgres-data) || true

echo "Starting containers..."
docker-compose up -d

echo "Containers started, you can check logs with: docker-compose logs -f" 