#!/bin/bash
echo "Generating Payment Service Stubs..."

protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    payment.proto

echo "Running server..."
nohup go run server/server.go &

echo "Running client..."
go run client/client.go
