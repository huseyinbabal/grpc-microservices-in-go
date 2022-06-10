package main

import (
	"context"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	shipping "listing_6.1"
	"log"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor()))
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Println("It is fine, this is not a complete example.")
	}

	defer conn.Close()

	shippingClient := shipping.NewShippingServiceClient(conn)
	ctx := context.Background()
	_, err = shippingClient.Create(ctx, &shipping.CreateShippingRequest{UserId: 23, OrderId: 123})
	if err != nil {
		log.Println("Don't worry, we don't expect to see it is working.")
		log.Println(err)
	} else {
		log.Println("Shipping is created successfully.")
	}
}
