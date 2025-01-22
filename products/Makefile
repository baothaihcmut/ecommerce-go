gen-proto-shared:
	@protoc --proto_path=../libs/proto/ \
			--go_out=./internal/adapter/grpc \
			--go-grpc_out=./internal/adapter/grpc \
			--go_opt=Mshared.proto=./proto \
			--go-grpc_opt=Mshared.proto=./proto \
			"../libs/proto/shared.proto"
gen-proto-categories:
	@protoc \
    --proto_path="../libs/proto/" \
    --go_out="./internal/adapter/grpc" \
    --go-grpc_out=require_unimplemented_servers=false:./internal/adapter/grpc \
    --go_opt=Mshared.proto=./proto \
    --go-grpc_opt=Mshared.proto=./proto \
    --go_opt=Mcategory.proto=./proto \
    --go-grpc_opt=Mcategory.proto=./proto \
    "../libs/proto/category.proto"
gen-proto-products:
	@protoc \
    --proto_path="../libs/proto/" \
    --go_out="./internal/adapter/grpc" \
    --go-grpc_out=require_unimplemented_servers=false:./internal/adapter/grpc \
    --go_opt=Mshared.proto=./proto \
    --go-grpc_opt=Mshared.proto=./proto \
    --go_opt=Mproduct.proto=./proto \
    --go-grpc_opt=Mproduct.proto=./proto \
    "../libs/proto/product.proto"