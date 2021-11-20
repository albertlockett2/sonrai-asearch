proto:
	protoc --go_out=./src --go-grpc_out=./src proto/search.proto

java:
	rm -rf ./java/src
	mkdir -p ./java/src/main/java
	protoc --plugin=protoc-gen-grpc-java=./protoc-gen-java --java_out=./java/src/main/java --grpc-java_out=./java/src/main/java ./proto/search.proto

install-java:
	cd ./java
	mvn package
	mvn install