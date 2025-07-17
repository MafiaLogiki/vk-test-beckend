.PHONY: run up down help setup

all: run

run: .env
	docker-compose up --build

up: run

down:
	docker-compose down

.env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
	fi

