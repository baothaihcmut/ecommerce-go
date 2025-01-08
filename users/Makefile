gen-proto-shared:
	@protoc --proto_path=../libs/proto/ \
			--go_out=./internal/adapter/grpc \
			--go-grpc_out=./internal/adapter/grpc \
			--go_opt=Mshared.proto=./proto \
			--go-grpc_opt=Mshared.proto=./proto \
			"../libs/proto/shared.proto"
gen-proto-users:
	@protoc --proto_path "../libs/proto/" \
			--go_out "./internal/adapter/grpc" \
			--go-grpc_out "./internal/adapter/grpc" \
			--go_opt "Mshared.proto=./proto" \
			--go-grpc_opt=Mshared.proto=./proto \
			--go_opt "Muser.proto=./proto" \
			--go-grpc_opt=Muser.proto=./proto \
			"../libs/proto/user.proto"
gen-proto-auth:
	@protoc --proto_path "../libs/proto/" \
			--go_out "./internal/adapter/grpc" \
			--go-grpc_out "./internal/adapter/grpc" \
			--go_opt "Mshared.proto=./proto" \
			--go-grpc_opt=Mshared.proto=./proto \
			--go_opt "Mauth.proto=./proto" \
			--go-grpc_opt=Mauth.proto=./proto \
			"../libs/proto/auth.proto"
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