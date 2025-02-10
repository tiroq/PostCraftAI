.PHONY: all build-backend build-frontend up down clean

all: build-backend build-frontend

build-backend:
	docker build -t postcraft-ai-backend ./backend

build-frontend:
	docker build -t postcraft-ai-frontend ./frontend

up:
	docker compose up --build

down:
	docker compose down

clean:
	docker compose down -v --rmi all
