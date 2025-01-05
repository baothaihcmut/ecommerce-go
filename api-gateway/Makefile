gen-proto-shared:
	@protoc --proto_path=../libs/proto/ \
			--go_out=./internal/common/proto/ \
			--go-grpc_out=./internal/common/proto/ \
			--go_opt=paths=source_relative \
			--go-grpc_opt=paths=source_relative \
			"../libs/proto/shared/shared.proto"
gen-proto-users:
	@protoc --proto_path=../libs/proto/ \
			--go_out=./internal/modules/users/proto/ \
			--go-grpc_out=./internal/modules/users/proto/ \
			--go_opt=paths=source_relative \
			--go-grpc_opt=paths=source_relative \
			"../libs/proto/users/user.proto"