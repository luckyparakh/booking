EVENT_SERVICE_BINARY=eventservice
## up: starts all containers in the background
up: 
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"
## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	sudo rm -rf db-data/*
	docker-compose down --remove-orphans
	@echo "Done!"
## build_broker: builds the event binary as a linux executable
build_event:
	@echo "Building Event Service binary..."
	cd ./src/eventservice && env GOOS=linux CGO_ENABLED=0 go build -o ${EVENT_SERVICE_BINARY} .
	@echo "Done!"
up_build: build_event down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"