#!/bin/sh

protoc --proto_path=proto/ --go_out=models/ --go_opt=paths=source_relative \
    --go-grpc_out=models/ --go-grpc_opt=paths=source_relative \
    proto/models.proto

protoc --proto_path=proto/ --go_out=twitter/ --go_opt=paths=source_relative \
    --go-grpc_out=twitter/ --go-grpc_opt=paths=source_relative \
    proto/twitter.proto
