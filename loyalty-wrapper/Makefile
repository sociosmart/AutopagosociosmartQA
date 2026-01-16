include .env.backend

migrate:
	docker compose run --rm backend beanie migrate -uri ${CONNECTION_URI} -db loyalty -p migrations/
