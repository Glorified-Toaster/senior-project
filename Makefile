CONFIG_FILE = config/config.yaml
MIGRATIONS_DIR = db/migrations

DB_HOST := $(shell yq '.database.host' $(CONFIG_FILE))
DB_PORT := $(shell yq '.database.port' $(CONFIG_FILE))
DB_USER := $(shell yq '.database.username' $(CONFIG_FILE))
DB_PASS := $(shell yq '.database.password' $(CONFIG_FILE))
DB_NAME := $(shell yq '.database.database_name' $(CONFIG_FILE))
DB_SSL  := $(shell yq '.database.ssl_mode' $(CONFIG_FILE))

DB_URL := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL)

GOOSE = ~/go/bin/goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)"
SQLC = ~/go/bin/sqlc
TEMPL = ~/go/bin/templ
.PHONY: migrate-up migrate-down migrate-status

# Goose 
migrate-up:
	$(GOOSE) up

migrate-down:
	$(GOOSE) down

migrate-status:
	$(GOOSE) status

# Sqlc
generate-query:
	$(SQLC) generate -f ./db/sqlc.yaml

# Tailwindcss 
tailwind-dev:
	tailwindcss -w -i ./web/static/src/css/input.css -o ./web/static/css/output.css --minify

tailwind-build:
	tailwindcss -m -i ./web/static/src/css/input.css -o ./web/static/css/output.css --minify

# Templ
generate-templ:
	$(TEMPL) generate

