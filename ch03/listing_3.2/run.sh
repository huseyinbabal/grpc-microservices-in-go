#!/bin/bash
GITHUB_USERNAME=$1
GITHUB_EMAIL=$2

SERVICE_NAME=order
RELEASE_VERSION=v1.2.3

echo "Generating Go source code"
mkdir -p golang
protoc --go_out=./golang \
  --go_opt=paths=source_relative \
  --go-grpc_out=./golang \
  --go-grpc_opt=paths=source_relative \
 ./${SERVICE_NAME}/*.proto

echo "Generated Go source code files"
ls -al ./golang/${SERVICE_NAME}

cd golang/${SERVICE_NAME}
go mod init \
  github.com/${GITHUB_USERNAME}/grpc-microservices-in-go/ch03/listing_3.2/golang/${SERVICE_NAME} || true
go mod tidy || true
cd ../../
git config --global user.email ${GITHUB_EMAIL}
git config --global user.name ${GITHUB_USERNAME}
git add . && git commit -am "proto update" || true
git push -u origin HEAD
git tag -d ch03/listing_3.2/golang/${SERVICE_NAME}/${RELEASE_VERSION}
git push --delete origin ch03/listing_3.2/golang/${SERVICE_NAME}/${RELEASE_VERSION}
git tag -fa ch03/listing_3.2/golang/${SERVICE_NAME}/${RELEASE_VERSION} \
  -m "ch03/listing_3.2/golang/${SERVICE_NAME}/${RELEASE_VERSION}"
git push origin refs/tags/ch03/listing_3.2/golang/${SERVICE_NAME}/${RELEASE_VERSION}
