generate:
	protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false \
	--go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/quotes.proto

run:
	go run ./cmd/