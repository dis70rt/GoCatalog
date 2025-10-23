up:
	@echo "Starting services..."
	docker compose up --build -d

down:
	@echo "Stopping services..."
	docker compose down

stop:	down

build:
	@echo "Building services..."
	docker compose build

clean:
	@echo "Cleaning up services, volumes, and networks..."
	docker compose down -v --remove-orphans