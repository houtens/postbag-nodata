# Local dev settings
PROJECT = postbag
DEV_DB = postbag
TEST_DB = postbag_test
DB_USER = postgres
DB_PASS = 
DB_HOST = localhost
DB_PORT = 5432

#------------------------------------------------------------------------------
# Helpers
#------------------------------------------------------------------------------

## help: show make commands
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' Makefile | column -t -s ':'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]


#------------------------------------------------------------------------------
# Quality Control
#------------------------------------------------------------------------------

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

#------------------------------------------------------------------------------
# Development
#------------------------------------------------------------------------------


## createdb: create the database
.PHONY: createdb
createdb:
	-psql -Upostgres -c 'create database $(DEV_DB)'

## dropdb: drop the database
.PHONY: dropdb
dropdb:
	-psql -Upostgres -c 'drop database $(DEV_DB)'

## migrateup: apply migrations
.PHONY: migrateup
migrateup:
	migrate -path database/migrations -database "postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DEV_DB)?sslmode=disable" -verbose up

## migratedown: reverse migrations
.PHONY: migratedown
migratedown:
	migrate -path database/migrations -database "postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DEV_DB)?sslmode=disable" -verbose down


## createdb/test: create the test database
.PHONY: createdb/test
createdb/test:
	-psql -Upostgres -c 'create database $(TEST_DB)'

## dropdb/test: drop the test database
.PHONY: dropdb/test
dropdb/test:
	-psql -Upostgres -c 'drop database $(TEST_DB)'

## migrateup/test: apply migrations for the test database
.PHONY: migrateup/test
migrateup/test:
	migrate -path database/migrations -database "postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(TEST_DB)?sslmode=disable" -verbose up

## migratedown/test: revert migrations on the test database
.PHONY: migratedown/test
migratedown/test:
	migrate -path database/migrations -database "postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(TEST_DB)?sslmode=disable" -verbose down

## sqlc: compile postgres data models
.PHONY: sqlc
sqlc:
	cd database && sqlc generate

## sqlc/vet: vet the sqlc queries for errors
.PHONY: sqlc/vet
sqlc/vet:
	cd database && sqlc vet

## reset: drop database, recreate and apply all migrations
.PHONY: reset
reset: dropdb createdb migrateup


# Run a subset of tests with "make test run=User"
# Where run is set to a string within the test name eg TestCreateUser()
## test: run tests
.PHONY: test
run = .
test: dropdb/test createdb/test migrateup/test
	go clean -testcache
	grc go test -v ./... --run $(run)

## test/cover: run tests with coverage report
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	# grc go test -v -cover ./...

## seed: seed the database
.PHONY: seed seed/fixtures seed/members seed/members seed/tournaments seed/invoices seed/ratings seed/results seed/fixwins
seed:
	go run seed/*.go -all
## seed/fixtures: seed the database with static fixtures
seed/fixtures:
	go run seed/*.go -fixtures
## seed/members: seed the database with members data
seed/members:
	go run seed/*.go -members
## seed/tournaments: seed the database with tournament data
seed/tournaments:
	go run seed/*.go -tournaments
## seed/invoices: seed the database with invoices data
seed/invoices:
	go run seed/*.go -invoices
## seed/ratings: seed the database with ratings data
seed/ratings:
	go run seed/*.go -ratings
## seed/results: seed the database with results data
seed/results:
	go run seed/*.go -results
## seed/fixwins: fix data to update wins etc
seed/fixwins:
	go run seed/*.go -fixwins

## countrows: show row counts in the postbag database
countrows:
	@psql postbag -c "VACUUM;"
	@psql postbag -c "SELECT table_name, (SELECT n_live_tup FROM pg_stat_user_tables WHERE relname = table_name) AS row_count FROM information_schema. tables WHERE table_schema = 'public';"

## dumpdb: backup the database with pg_dump
.PHONY: dumpdb restoredb
dumpdb:
	@pg_dump -Fc $(PROJECT) > database/backups/$(PROJECT)-$(shell date +%Y%m%d-%H%M).pg
## restoredb: (TODO) restore the database from a dump file
restoredb:
	@echo "pg_restore -d postbag postbag-DATE.pg"


## run: run the dev server
.PHONY: run
run:
	go run ./cmd/${PROJECT}/...

## css: build the tailwind css styles
.PHONY: css
css:
	npm run build

## css/watch: listen for changes and rebuild css
.PHONY: css/watch
css/watch:
	npm run watch

