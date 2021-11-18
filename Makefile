proto:
	protoc \
		--go_out=./src \
		--go-grpc_out=./src \
		proto/search.proto