module github.com/baothaihcmut/Ecommerce-go/users

go 1.23.7

require (
	github.com/google/uuid v1.6.0
	github.com/samber/lo v1.49.1
	golang.org/x/crypto v0.36.0
	golang.org/x/text v0.23.0 // indirect
)

require (
	github.com/baothaihcmut/Ecommerce-go/libs v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v5 v5.7.2
	google.golang.org/grpc v1.71.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250313205543-e70fdf4c4cb4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

replace github.com/baothaihcmut/Ecommerce-go/libs => ../libs
