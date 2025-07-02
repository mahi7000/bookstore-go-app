#!/bin/sh
set -e

echo "Running database migrations..."
goose -dir sql/schema postgres "$DB_URL" up

echo "Starting the application..."
./main
