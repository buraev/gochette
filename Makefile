.PHONY: build up down logs restart

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

start: build up

restart: down up

rebuild:
	docker compose down
	docker compose build --no-cache  # удаляет кэш и пересобирает
	docker compose up -d
