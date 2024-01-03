include .envrc

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

###################################
# DEVELOPMENT
###################################

## run: run the cmd/web application
.PHONY: run
run:
	go run ./cmd/web -addr="${ADDR}" -dsn="${DB_CONN_STRING}"

## migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	@echo 'Creating migration file'
	migrate create -seq -ext=.sql -dir=./migrations $(filter-out $@,$(MAKECMDGOALS))

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database="mysql://${DB_CONN_STRING}" up

###################################
# QUALITY CONTROL
###################################
.PHONY: audit
audit:
	@echo 'Verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
