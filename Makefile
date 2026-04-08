.PHONY: help confirm run/api db/migrations/create db/migrations/up db/migrations/down db/migrations/version force

# Help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

## run/api: run the cmd/api application
run/api:
	@echo "Starting the application..."
	@go run ./cmd/api

## db/migrations/create name=$1: create a new database migration
db/migrations/create:
	@echo "Creating new migration: $(name)"
	@migrate create -seq -ext=.sql -dir=./migrations $(name)

## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo "Applying up migrations..."
	@migrate -path=./migrations -database=$(FLICK_DB_DSN) up

## db/migrations/down: apply one down database migration
db/migrations/down:
	@echo "Applying down migrations..."
	@migrate -path=./migrations -database=$(FLICK_DB_DSN) down 1

## db/migrations/version: print the current migration version
db/migrations/version:
	@echo "Current migration version: "
	@migrate -path=./migrations -database=$(FLICK_DB_DSN) version

## force version=$1: force the migration version to a specific value
force:
	@echo "Forcing migration version to: $(version)"
	@migrate -path=./migrations -database=$(FLICK_DB_DSN) force $(version)
