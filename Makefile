PROTOC      := $(HOME)/tools/protoc/bin/protoc
PROTOC_OPTS := --proto_path=proto \
               --go_out=pb --go_opt=paths=source_relative \
               --go-grpc_out=pb --go-grpc_opt=paths=source_relative

.PHONY: proto build

proto:
	PATH="$(HOME)/go/bin:$$PATH" $(PROTOC) $(PROTOC_OPTS) proto/rates.proto

build:
	go build -o gRPC_currentrate .
