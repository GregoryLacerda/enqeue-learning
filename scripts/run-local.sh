#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

# Start RabbitMQ using Docker
docker-compose up -d rabbitmq

# Start the producer API
go run cmd/main.go &

# Wait for all background processes to finish
wait