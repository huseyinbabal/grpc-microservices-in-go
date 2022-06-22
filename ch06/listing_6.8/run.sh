#!/bin/bash
sudo apt-get install -y protobuf-compiler golang-goprotobuf-dev

echo "Installing protoc go and grpc modules..."

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo "Generating Order Service Stubs..."


protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    order.proto

go mod tidy

echo "####START####"


mkdir -p cert && cd cert
echo "Generating private key and self-signed certificate for CA..."
openssl req -x509 \
    -sha256 \
    -newkey rsa:4096 \
    -days 365 \
    -keyout ca.key \
    -out ca.crt \
    -subj "/C=TR/ST=EURASIA/L=ISTANBUL/O=Software/OU=Microservices/CN=*.microservices.dev/emailAddress=huseyin@microservices.dev" \
    -nodes

echo "Generate private key and certificate signing request for server"
openssl req \
    -sha256 \
    -newkey rsa:4096 \
    -keyout server.key \
    -out server-req.pem \
    -subj "/C=TR/ST=EURASIA/L=ISTANBUL/O=Microservices/OU=PaymentService/CN=*.microservices.dev/emailAddress=huseyin@microservices.dev" \
    -nodes

echo "Sign certificate signing request for server by using private key of CA"
rm server-ext.cnf || true && echo "subjectAltName=DNS:*.microservices.dev,DNS:*.microservices.dev,IP:0.0.0.0" >> server-ext.cnf
openssl x509 \
    -req -in server-req.pem \
    -sha256 \
    -days 60 \
    -CA ca.crt \
    -CAkey ca.key \
    -CAcreateserial \
    -out server.crt \
    -extfile server-ext.cnf

echo "Generate private key and certificate signing request for client"
openssl req \
    -sha256 \
    -newkey rsa:4096 \
    -keyout client.key \
    -out client-req.pem \
    -subj "/C=TR/ST=EURASIA/L=Istanbul/O=Microservices/OU=OrderService/CN=*.microservices.dev/emailAddress=huseyin@microservices.dev" \
    -nodes

echo "Sign certificate signing request for client by using private key of CA"
rm client-ext.cnf || true && echo "subjectAltName=DNS:*.microservices.dev,DNS:*.microservices.dev,IP:0.0.0.0" >> client-ext.cnf
openssl x509 \
    -req -in client-req.pem \
    -sha256 \
    -days 60 \
    -CA ca.crt \
    -CAkey ca.key \
    -CAcreateserial \
    -out client.crt \
    -extfile client-ext.cnf

cd ../
echo "Running server..."
nohup go run server/server.go &

echo "Waiting for order service to be up..."
sleep 5

echo "Running client..."
go run client/client.go
killall server
echo "####END####"


