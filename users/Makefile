gen-proto:
	protoc \
		--proto_path=../libs/proto/ "../libs/proto/shared.proto" \
		--go_out=./internal/adapter/grpc/proto/ --go_opt=paths=source_relative \
  	--go-grpc_out=./internal/adapter/grpc/proto/ --go-grpc_opt=paths=source_relative
	@protoc \
		--proto_path=../libs/proto/ "../libs/proto/user.proto" \
		--go_out=./internal/adapter/grpc/proto/ --go_opt=paths=source_relative \
  	--go-grpc_out=./internal/adapter/grpc/proto/ --go-grpc_opt=paths=source_relative
gen-sqlc:
	@cd internal/adapter/persistence/sqlc && sqlc generate
	@cd ../../../../

run:
	@go run cmd/users/main.go