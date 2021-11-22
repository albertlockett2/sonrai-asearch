proto:
	protoc --go_out=./src --go-grpc_out=./src proto/search.proto

java:
	rm -rf ./java/src
	mkdir -p ./java/src/main/java
	protoc --plugin=protoc-gen-grpc-java=./protoc-gen-java --java_out=./java/src/main/java --grpc-java_out=./java/src/main/java ./proto/search.proto

install-java: java
	cd ./java && mvn install

golang:
	mkdir -p ./build
	go mod vendor
	go build -o ./build/asearch ./src/main.go

docker: golang
	docker build . -t asearch:latest

all: proto install-java