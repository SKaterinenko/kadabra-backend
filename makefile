include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
PG_DERICTORY=./internal/core/migrations/postgres

run:
	go run ./cmd/app/main.go

migrate-pg-up:
	migrate -path "${PG_DERICTORY}" -database "$(DB_URL)" up

migrate-pg-down:
	migrate -path "${PG_DERICTORY}" -database "$(DB_URL)" down 1

migrate-pg-create:
	migrate create -ext sql -dir "${PG_DERICTORY}" $(name)

migrate-pg-force:
ifndef version
	$(error Please specify version: make migrate-pg-force version=123)
endif
	migrate -path "${PG_DERICTORY}" -database "$(DB_URL)" force $(version)

#TODO сделать время в названии бэкапа
pg-backup:
	docker exec $(DOCKER_CONTAINER) pg_dumpall -U $(DB_USER) > backup_all.sql

pg-restore:
	 docker exec -i $(DOCKER_CONTAINER) psql -U $(DB_USER) -d postgres < backup_all.sql

#pg-restore-windows:
#	Get-Content backup_all.sql | docker exec -i $(DOCKER_CONTAINER) psql -U $(DB_USER) $(DB_NAME)