createdb = createdb
dropdb = dropdb
migrate = migrate

config_path ?= ./etc/bank.ini.example
migrations_path ?= ./bank/postgres/migrations

dsn = $(shell sed -nr 's/dsn="(.*)"/\1/p' $(config_path) | tr -d [:space:])
db_host = $(shell echo $(dsn) | sed -nr 's|.*//[^@]*@([^/]*).*|\1|p')
db_user = $(shell echo $(dsn) | sed -nr 's|.*//([^:@]*)[@:].*|\1|p')
db_name = $(shell echo $(dsn) | sed -nr 's|.*/([^?]*).*|\1|p')

default: help

## db:create: Create database.
.PHONY: db\:create
db\:create:
	$(createdb) -h $(db_host) -U $(db_user) $(db_name)

## db:drop: Drop database.
.PHONY: db\:drop
db\:drop:
	$(dropdb) -h $(db_host) -U $(db_user) $(db_name)

## db:migrate:up: Run database migrations (UP). Optional: config=<path-to-config>
.PHONY: db\:migrate\:up
db\:migrate\:up:
	$(migrate) -path=$(migrations_path) -database=$(dsn) up

## db:migrate:down: Run database migrations (DOWN). Optional: config=<path-to-config>
.PHONY: db\:migrate\:down
db\:migrate\:down:
	$(migrate) -path=$(migrations_path) -database=$(dsn) down

## db:migrate:force: Force dirty database migrations for given version=<version>. Optional: config=<path-to-config>
.PHONY: db\:migrate\:force
db\:migrate\:force:
	$(migrate) -path=$(migrations_path) -database=$(dsn) force $(version)

## migration:create: Create a migration with specified name=<name>. Optional: migrations_path=<path-to-migrations>
.PHONY: migration\:create
migration\:create:
	$(migrate) create -ext=.sql -dir=$(migrations_path) $(name)

## help: Show this help.
.PHONY: help
help: Makefile
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'