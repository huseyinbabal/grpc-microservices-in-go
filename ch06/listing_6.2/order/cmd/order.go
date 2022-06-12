package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"listing_6.2/order"
	"listing_6.2/product"
	"log"
	"net"
	"time"
)

type server struct {
	order.UnimplementedOrderServiceServer
	productClient product.ProductServiceClient
}

func (s *server) Create(ctx context.Context, in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	time.Sleep(2 * time.Second)
	productInfo, err := s.productClient.Get(ctx, &product.GetProductRequest{ProductId: in.ProductId})
	if err != nil {
		return nil, err
	}
	return &order.CreateOrderResponse{OrderId: 123, ProductTitle: productInfo.Title}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	var productOpts []grpc.DialOption
	productOpts = append(productOpts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:8081", productOpts...)
	if err != nil {
		log.Fatalf("Failed to connect product service. Err: %v", err)
	}

	order.RegisterOrderServiceServer(grpcServer, &server{
		productClient: product.NewProductServiceClient(conn),
	})
	grpcServer.Serve(listener)
}
