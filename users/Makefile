gen-proto:
	@go run cmd/proto-gen/main.go
gen-sqlc:
	@cd internal/adapter/persistence/sqlc && sqlc generate
	@cd ../../../../
migrate-up-dev:
	@go run cmd/migrations/main.go up dev

migrate-up-prod:
	@go run cmd/migrations/main.go up prod

run-dev:
	@go run cmd/users/main.go dev

run-prod:
	@go run cmd/users/main.go prod