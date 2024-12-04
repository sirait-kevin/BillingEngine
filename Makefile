# Define variables
DOCKER_COMPOSE = docker-compose
APP_NAME = BillingEngine_app
DOCKER_COMPOSE_YML = docker-compose.yml
BUILD_DIR = .

# Build the Docker containers
build:
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

# Run Docker containers in the background
up:
	@echo "Starting the application and services..."
	$(DOCKER_COMPOSE) up -d

# Run the application locally (inside the container)
run:
	@echo "Running the app..."
	$(DOCKER_COMPOSE) exec $(APP_NAME) ./main

# Stop Docker containers
down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

# Show logs of the app
logs:
	@echo "Fetching logs from the app container..."
	$(DOCKER_COMPOSE) logs -f $(APP_NAME)

# Clean Docker volumes and images
clean:
	@echo "Cleaning up Docker volumes and images..."
	$(DOCKER_COMPOSE) down -v --rmi all --remove-orphans

# Rebuild and restart the containers
rebuild: down build up

# Run the database setup manually (if needed)
db-setup:
	@echo "Running database setup..."
	$(DOCKER_COMPOSE) run --rm setup

# Start the application in debug mode
debug:
	@echo "Starting the application in debug mode..."
	$(DOCKER_COMPOSE) exec $(APP_NAME) ./main -debug=true

# Run tests in the app (if applicable, assuming you have a test command)
test:
	@echo "Running tests..."
	$(DOCKER_COMPOSE) exec $(APP_NAME) go test ./...

# Help command to show available Makefile commands
help:
	@echo "Makefile commands:"
	@echo "  build         - Build Docker containers"
	@echo "  up            - Start all Docker containers in the background"
	@echo "  run           - Run the application inside the container"
	@echo "  down          - Stop all Docker containers"
	@echo "  logs          - Fetch and follow logs for the app"
	@echo "  clean         - Clean up Docker containers, volumes, and images"
	@echo "  rebuild       - Rebuild and restart all containers"
	@echo "  db-setup      - Run database setup script (initializing DB and inserting sample data)"
	@echo "  debug         - Start the app in debug mode"
	@echo "  test          - Run the app's tests (if defined)"
	@echo "  help          - Show this help message"
