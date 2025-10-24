

.DEFAULT_GOAL := help

up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f

restart:
	docker compose down
	docker compose up -d --build

clean:
	docker system prune -af
	docker volume prune -f

help:
	@echo "Makefile commands:"
	@echo "  up       - Start the server with Docker Compose"
	@echo "  down     - Stop the server and remove containers"
	@echo "  logs     - View the server logs"
	@echo "  restart  - Restart the server with Docker Compose"
	@echo "  help     - Show this help message"