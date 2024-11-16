# Technical Stack

- `go-migrate` - database migrations
- `sqlc` - data model generation from schema and query files
- `gnu make` - wrapper for management commands
- `gin/gonic` - Gin web framework
- `tailwindcss` - css
- `htmx` - asynchronous refresh of partial content

# Files and Directories 

- `database/migrations` - schema in the form of migration files
- `database/query` - sql CRUD queries definiting the model behaviour
- `internal/model` - models generated from sqlc and manually created test files
- `util` - go utility glue to hold this together
- `app.env` - environment variables needed at runtime
- `sqlc.yaml` - config file for sqlc

# Useful commands

The `Makefile` has been provided to help manage commonly used actions on the project which include:

- `make createdb` - create the development postgres database
- `make dropdb` - drop the development database
- `make migrateup` - apply all up migrations
- `make migratedown` - revert all migrations

- `make sqlc` - generate the models from sqlc defintions
- `make reset` - drop and recreate the database, then apply all migrations
- `make test` - run the test suite
- `make seed` - seed the database from fixture and exported data files
- `make run` - run the dev server

- `make createtestdb` - create the test database
- `make droptestdb` - drop the test database
- `make migratetestup` - apply all up migrations to the test database
- `make migratetestdown` - revert all migrations to the test database

- `migrate create -ext sql -dir database/migrations -seq <migration-name> - to create migration up/down files`

