#!/bin/bash

#sudo apt-get install -y protobuf-compiler golang-goprotobuf-dev

echo "Installing protoc go and grpc modules..."

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo "Generating Payment Service Stubs..."

protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    shipping.proto

go mod tidy

echo "####START####"

echo "Running client..."
nohup go run client/client.go &

echo "Wait for 5 seconds and run service."
sleep 30

echo "Running server..."
go run server/server.go
echo "####END####"


