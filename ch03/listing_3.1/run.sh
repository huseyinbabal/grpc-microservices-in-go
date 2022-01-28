#!/bin/bash

echo "Generating Go source code for order.proto"
mkdir -p golang
protoc -I ./proto \
    --go_out=./golang \
    --go_opt=paths=source_relative \
    --go-grpc_out=./golang \
    --go-grpc_opt=paths=source_relative \
    ./proto/order.proto

echo "Generated files"
ls -al golang

