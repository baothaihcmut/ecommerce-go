gen-proto-shared:
	@protoc --proto_path=../libs/proto/ \
			--go_out=./internal/adapter/grpc/proto/ \
			--go-grpc_out=./internal/adapter/grpc/proto/ \
			--go_opt=Mshared/shared.proto=./shared \
			--go-grpc_opt=Mshared/shared.proto=./shared \
			"../libs/proto/shared/shared.proto"
gen-proto-users:
	@protoc --proto_path=../libs/proto/ \
			--go_out=./internal/adapter/grpc/proto \
			--go-grpc_out=./internal/adapter/grpc/proto \
			--go_opt=Musers/user.proto=./users \
			--go-grpc_opt=Musers=./users \
			"../libs/proto/users/user.proto"
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