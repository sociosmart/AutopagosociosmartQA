connect-db:
	docker compose exec db mariadb -u payments1 --password=randompassword smartgas-payments

docs:
	docker compose exec backend swag init -g internal/server/server.go

gen-wire:
	docker compose exec backend wire internal/injectors/wire.go

schedule-sync:
	docker compose exec backend /go/bin/app scheduleSync

up:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up
qa-up:
	docker compose -f docker-compose.yml -f docker-compose.qa.yml up

up-d:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d
qa-up-d:
	docker compose -f docker-compose.yml -f docker-compose.qa.yml up -d

stop:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml stop

qa-stop:
	docker compose -f docker-compose.yml -f docker-compose.qa.yml stop

build:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml build

qa-build:
	docker compose -f docker-compose.yml -f docker-compose.qa.yml build

migrate:
	docker compose exec backend go run . migrate

generate:
	docker compose exec backend go generate ./...

createuser:
	docker compose exec backend go run . createuser

test:
	docker compose exec -e GIN_MODE=test backend go test -v -count=1 \
																					./api/v1/controllers \
																					./internal/middlewares \
																					./internal/server/app
test-coverage:
	docker compose exec -e GIN_MODE=test backend go test -v -count=1 \
																					-coverprofile=./tmp/coverage.out \
																					./api/v1/controllers \
																					./internal/middlewares \
																					./internal/server/app \
																					&& clear \
																					&& docker compose exec backend go tool cover -func=./tmp/coverage.out

sync:
	docker compose exec backend go run . sync $(filter-out $@,$(MAKECMDGOALS))

sh:
	docker compose exec backend sh

create-authorized-app:
	docker compose exec backend go run . authorizedApp $(filter-out $@,$(MAKECMDGOALS))

attach-app:
	docker compose exec backend sh
