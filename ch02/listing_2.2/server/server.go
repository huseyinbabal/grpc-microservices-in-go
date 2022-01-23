package main

import (
	"fmt"
	"google.golang.org/grpc"
	payment "listing_2.2"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	payment.RegisterPaymentServiceServer(grpcServer, payment.UnimplementedPaymentServiceServer{})
	grpcServer.Serve(listener)
}
