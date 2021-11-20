gen-js:
	mkdir -p ./src/gen
	npx grpc_tools_node_protoc \
    --js_out=import_style=commonjs,binary:./src/gen \
    --grpc_out=./src/gen \
    --plugin=protoc-gen-grpc=./node_modules/.bin/grpc_tools_node_protoc_plugin \
    -I ../proto \
    ../proto/*.proto

gen-ts:
	mkdir -p ./src/gen
	npx grpc_tools_node_protoc \
    --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
    --ts_out=./src/gen \
    -I ../proto \
    ../proto/*.proto


gen: gen-js gen-ts

all: gen