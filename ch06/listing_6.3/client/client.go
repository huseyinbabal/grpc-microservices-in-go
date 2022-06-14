package main

import (
	"context"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	shipping "listing_6.3"
	"log"
	"time"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
		grpc_retry.WithMax(5),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
	)))
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect shipping service. Err: %v", err)
	}

	defer conn.Close()

	shippingServiceClient := shipping.NewShippingServiceClient(conn)
	log.Println("Creating shipping...")
	_, errCreate := shippingServiceClient.Create(context.Background(), &shipping.CreateShippingRequest{
		UserId:  23,
		OrderId: 2344,
	})
	if errCreate != nil {
		log.Printf("Failed to create shipping. Err: %v", errCreate)
	} else {
		log.Println("Shipping is created successfully.")
	}
}
