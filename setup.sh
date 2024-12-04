#!/bin/bash

# Setup the database and services using docker-compose
echo "Setting up BillingEngine services..."

# Start the containers
docker-compose up --build -d

# Wait for MySQL and NSQ services to be ready
echo "Waiting for MySQL to be ready..."
sleep 10

# Run SQL script to initialize the database
echo "Initializing database..."
docker exec -i BillingEngine_mysql_1 mysql -uroot -prootpassword BillingEngine < ./sqlfiles/init.sql

echo "Setup complete."
