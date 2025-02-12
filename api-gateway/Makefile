gen-proto-shared:
	@protoc --proto_path=../libs/proto/ \
			--go_out=./internal/common/grpc \
			--go-grpc_out=./internal/common/grpc \
			--go_opt=Mshared.proto=./proto \
			--go-grpc_opt=Mshared.proto=./proto \
			"../libs/proto/shared.proto"
gen-proto-auth:
	@protoc \
    --proto_path="../libs/proto/" \
    --go_out="./internal" \
    --go-grpc_out=require_unimplemented_servers=false:./internal \
    --go_opt=Mshared.proto=github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/grpc/proto \
    --go-grpc_opt=Mshared.proto=github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/grpc/proto \
    --go_opt=Mauth.proto=./modules/auth/proto \
    --go-grpc_opt=Mauth.proto=./modules/auth/proto \
    "../libs/proto/auth.proto"
gen-proto-categories:
	@protoc \
    --proto_path="../libs/proto/" \
    --go_out="./internal" \
    --go-grpc_out=require_unimplemented_servers=false:./internal \
    --go_opt=Mshared.proto=github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/grpc/proto \
    --go-grpc_opt=Mshared.proto=github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/grpc/proto \
    --go_opt=Mcategory.proto=./modules/products/proto \
    --go-grpc_opt=Mcategory.proto=./modules/products/proto \
    "../libs/proto/category.proto"
