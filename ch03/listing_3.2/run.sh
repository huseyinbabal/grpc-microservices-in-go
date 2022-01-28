#!/bin/bash

SERVICE_NAME=order
RELEASE_VERSION=v1.2.3

protoc --go_out=./golang \
  --go_opt=paths=source_relative \
  --go-grpc_out=./golang \
  --go-grpc_opt=paths=source_relative \
 ./${SERVICE_NAME}/*.proto

cd golang/${SERVICE_NAME}
go mod init \
  github.com/huseyinbabal/microservices-proto/golang/${SERVICE_NAME} || true
go mod tidy
cd ../../
git config --global user.email "huseyinbabal88@gmail.com"
git config --global user.name "Huseyin BABAL"
git add . && git commit -am "proto update" || true
git tag -fa golang/${SERVICE_NAME}/${RELEASE_VERSION} \
  -m "golang/${SERVICE_NAME}/${RELEASE_VERSION}"
git push origin refs/tags/golang/${SERVICE_NAME}/${RELEASE_VERSION}
