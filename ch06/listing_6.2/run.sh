#!/bin/bash

sudo apt-get install -y protobuf-compiler golang-goprotobuf-dev

echo "Installing protoc go and grpc modules..."

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo "Generating Order Service Stubs..."
cd order
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    order.proto

echo "Generating Product Service Stubs..."
cd ../product
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    product.proto
cd ..
go mod tidy

echo "####START####"

echo "Running order..."
nohup go run order/cmd/order.go &

echo "Running product..."
nohup go run product/cmd/product.go &

echo "Wait for 5 seconds and run service."
sleep 5

echo "Running client..."
go run client/client.go
echo "####END####"


