start_docker:
	@echo "Starting Docker containers..."
	docker compose -f security/deployment/docker-compose.yml up --build -d

stop_docker:
	@echo "Stopping Docker containers..."
	docker compose -f security/deployment/docker-compose.yml down