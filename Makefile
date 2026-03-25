## Migrations

.PHONY: migrate-create migrate-up migrate-down migrate-version migrate-force

# Create a new local migration (usage: make migrate-create name=create_users_table)
migrate-create:
	@migrate create -seq -ext=.sql -dir=./migrations $(name)

# Apply all up migrations
migrate-up:
	@migrate -path=./migrations -database=$(FLICK_DB_DSN) up

# Apply all down migrations (by default rolls back 1 version)
migrate-down:
	@migrate -path=./migrations -database=$(FLICK_DB_DSN) down 1

# Check current migration version
migrate-version:
	@migrate -path=./migrations -database=$(FLICK_DB_DSN) version

# Force migration version (usage: make migrate-force version=1)
migrate-force:
	@migrate -path=./migrations -database=$(FLICK_DB_DSN) force $(version)
